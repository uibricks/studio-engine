package mapping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/streadway/amqp"
	"github.com/uibricks/studio-engine/internal/app/mapping/constants"
	"github.com/uibricks/studio-engine/internal/app/mapping/utils"
	"github.com/uibricks/studio-engine/internal/pkg/db"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	expressionpb "github.com/uibricks/studio-engine/internal/pkg/proto/expression"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	"github.com/uibricks/studio-engine/internal/pkg/rabbitmq"
	"github.com/uibricks/studio-engine/internal/pkg/redis"
	"github.com/uibricks/studio-engine/internal/pkg/request"
	pkgUtils "github.com/uibricks/studio-engine/internal/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type MappingServer struct {
	DB         *pg.DB
	Rabbit     *rabbitmq.Rabbit
	Redis      *redis.Redis
	ReplyQueue amqp.Queue
	Channel    *rabbitmq.Channel
}

func ProvideMappingServer(dbClient db.DbClient, redis *redis.Redis, rabbit *rabbitmq.Rabbit, replyQ amqp.Queue, ch *rabbitmq.Channel) *MappingServer {
	return &MappingServer{
		DB:         dbClient.Connection,
		Redis:      redis,
		Rabbit:     rabbit,
		ReplyQueue: replyQ,
		Channel:    ch,
	}
}

func (m *MappingServer) UpdateRepositoryDetails(ctx context.Context, req *mappingpb.UpdateRepositoryRequest) (*mappingpb.UpdateRepositoryResponse, error) {

	// get mapping from cache
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if not in cache, fetch from DB : both projectId and projectVersion are required
	if len(mapping.GetProjectId()) == 0 {
		if err = m.LoadMappingWithVersion(mapping, req.GetProjectId(), req.GetProjectVersion()); err != nil && !strings.Contains(err.Error(), constants.NotFound) {
			return nil, err
		}
	}

	// update mapping and response object with data from request
	res := &mappingpb.UpdateRepositoryResponse{ProjectId: req.GetProjectId()}

	if req.RepositoryMenu != nil || req.EmptyRepositoryMenu {
		mapping.GetConfig().RepositoryMenu = req.RepositoryMenu
		res.RepositoryMenu = req.RepositoryMenu
	}

	if req.Repositories != nil {

		if mapping.GetConfig().Repositories == nil {
			mapping.GetConfig().Repositories = make(map[string]*mappingpb.Repository)
		}

		for key, repo := range req.Repositories {
			if repo.GetAuthentication().GetType() == "" {
				repo.Authentication = &mappingpb.Authentication{Type: constants.DefaultAuthenticationType}
			}
			mapping.GetConfig().Repositories[key] = repo
		}
		res.Repositories = req.Repositories
	}

	// set updated mapping back to cache. updated mapping in cache can be used to get, save and delete mapping
	err = m.SetMappingInCache(ctx, mapping, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetMapping get the mapping info depending on the mapping request
// First priority is the cache and if the cache does not exist then it will retrieve it from database
func (m *MappingServer) GetMapping(ctx context.Context, req *mappingpb.MappingRequest) (*mappingpb.MappingResponse, error) {

	// get mapping from cache
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if not in cache, fetch from DB only if the request contains a project version.
	if len(mapping.GetProjectId()) == 0 {
		if err = m.LoadMappingWithVersion(mapping, req.GetProjectId(), req.GetProjectVersion()); err != nil {
			return nil, err
		}
	}

	// set response object with data from cache/db
	res := utils.GetEmptyConfig()

	res.DefaultEnvironment = mapping.GetConfig().GetDefaultEnvironment()

	if req.GetIncludeEnvs() {
		res.Environments = mapping.GetConfig().GetEnvironments()
		res.EnvironmentVariables = mapping.GetConfig().GetEnvironmentVariables()
	}

	if req.GetIncludeMenu() {
		res.RepositoryMenu = mapping.GetConfig().GetRepositoryMenu()
	}

	for _, repo := range req.GetRepositoryIds() {
		if _, found := mapping.GetConfig().GetRepositories()[repo]; found {
			res.Repositories[repo] = mapping.GetConfig().GetRepositories()[repo]
		}
	}

	return &mappingpb.MappingResponse{Config: res}, nil
}

// SaveMapping save the repository details from cache to db
func (m *MappingServer) SaveMapping(ctx context.Context, req *mappingpb.SaveMappingRequest) (*mappingpb.SaveMappingResponse, error) {
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if cache does not exist then duplicate the latest copy in the database with project version
	if len(mapping.GetProjectId()) == 0 {
		res, err := m.DB.Query(nil, fmt.Sprintf(constants.QueryDuplicateActiveMapping, req.GetNewProjectVersion(), req.GetProjectId()))
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to duplicate mapping %v", err))
		}
		// if there are no existing records in the database then we need to insert empty mapping record
		if res.RowsAffected() == 0 {
			_, err := m.DB.Query(nil, fmt.Sprintf(constants.QueryInsertEmptyMapping, req.GetProjectId(), req.GetNewProjectVersion()))
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("failed to insert new mapping record %v", err))
			}
		}
	} else {
		mapping.ProjectVersion = req.GetNewProjectVersion()
		mapping.ProjectId = req.GetProjectId()

		// to overwrite createdAt in cache
		mapping.Id = 0
		mapping.CreatedAt = pkgUtils.DbFormatDate(time.Now())

		_, err = m.DB.Model(mapping).Insert()
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to save the mapping with project(%s), version(%d) - %v", req.GetProjectId(), req.GetNewProjectVersion(), err))
		}
	}

	err = m.Redis.HDel(ctx, "test", redis.Misc, "admin", req.GetProjectId())
	if err != nil {
		logger.WithContext(ctx).Errorf("Failed to delete cache for repo with project id : %s - %v.", req.GetProjectId(), err)
	}

	m.MoveReposToHistory(ctx, req.GetProjectId(), req.CurrProjectVersion)

	return &mappingpb.SaveMappingResponse{ProjectId: req.GetProjectId(), CreatedAt: pkgUtils.DbFormatDate(time.Now())}, nil
}

func (m *MappingServer) DeleteRepository(ctx context.Context, req *mappingpb.DeleteRepositoryRequest) (*mappingpb.DeleteRepositoryResponse, error) {
	// get mapping from cache
	mapping, err := m.GetMappingFromCache(ctx, req.GetProjectId())

	if err != nil {
		return nil, err
	}

	// if not in cache, fetch from DB : both projectId and projectVersion are required
	if len(mapping.GetProjectId()) == 0 {
		if err = m.LoadMappingWithVersion(mapping, req.GetProjectId(), req.GetProjectVersion()); err != nil {
			return nil, err
		}
	}

	// if the project not found in db or cache, return Not Found
	if len(mapping.GetProjectId()) == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Mapping not found in cache or db with project id(%s)", req.GetProjectId()))
	}

	// if the repo not found in project key in cache, return Not Found
	if mapping.GetConfig().GetRepositories()[req.GetRepositoryId()] == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Repository not found with id(%s) in project with id in(%s)", req.GetRepositoryId(), req.GetProjectId()))
	}

	dependencies := make(map[string]string)
	// Check the repository dependencies
	for repoId, repo := range mapping.GetConfig().GetRepositories() {
		if repoId != req.GetRepositoryId() {
			repoStr, err := utils.MarshalToString(repo)
			if err != nil {
				return nil, status.Error(codes.Unknown, fmt.Sprintf("Failed in string conversion while checking for the dependencies for the repo id(%s)", repoId))
			}
			if strings.Contains(repoStr, req.GetRepositoryId()) {
				dependencies[repoId] = ""
			}
		}
	}

	// update the dependency names and send the dependencies
	if len(dependencies) > 0 {
		utils.UpdateRepoNames(dependencies, mapping.GetConfig().GetRepositoryMenu())
		return &mappingpb.DeleteRepositoryResponse{
			ProjectId:    req.GetProjectId(),
			Dependencies: utils.MapToArrayDependencies(dependencies),
		}, nil
	}

	//update the cache with changes
	deleted := false
	mapping.Config.RepositoryMenu = utils.DeleteRepoFromMenu(req.GetRepositoryId(), mapping.GetConfig().GetRepositoryMenu(), &deleted)
	delete(mapping.Config.Repositories, req.GetRepositoryId())

	err = m.SetMappingInCache(ctx, mapping, req.GetProjectId())

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Failed to delete repository(%s) in project(%s) from cache. Error - %v", req.GetRepositoryId(), req.GetProjectId(), err))
	}

	return &mappingpb.DeleteRepositoryResponse{
		ProjectId:      req.GetProjectId(),
		UpdatedAt:      pkgUtils.DbFormatDate(time.Now()),
		RepositoryMenu: mapping.GetConfig().GetRepositoryMenu(),
	}, nil
}

func (m *MappingServer) DeleteMapping(ctx context.Context, req *mappingpb.DeleteMappingRequest) (*mappingpb.DeleteMappingResponse, error) {

	softDeleted, err := m.SoftDeleteMapping(ctx, req)

	if err != nil {
		return nil, err
	}

	permDeleted := false

	if !softDeleted {
		permDeleted, err = m.PermanentlyDeleteMapping(ctx, req.GetProjectId())

		if err != nil {
			return nil, err
		}
	}

	if !softDeleted && !permDeleted {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("project not found with id %s", req.GetProjectId()))
	}

	return &mappingpb.DeleteMappingResponse{ProjectId: req.GetProjectId(), DeletedAt: pkgUtils.DbFormatDate(time.Now())}, nil
}

func (m *MappingServer) ExecuteAPI(ctx context.Context, req *mappingpb.ExecuteApiRequest) (*mappingpb.ExecuteApiResponse, error) {

	reqHeaders := make(map[string]string)
	headers := req.GetHeaders()

	for _, hdr := range headers {
		reqHeaders[hdr.GetKey()] = hdr.GetValue()
	}

	reqQParams := make(map[string]string)
	qParams := req.GetQueryParams()
	for _, qp := range qParams {
		reqQParams[qp.GetKey()] = qp.GetValue()
	}

	reqBody := make(map[string]interface{})
	_ = json.Unmarshal([]byte(req.GetBody()), &reqBody)

	extReq := utils.ExternalRequest{URL: req.GetUrl(), Type: req.GetHttpMethod(), Body: reqBody, Headers: reqHeaders, Params: reqQParams, SslCert: req.GetSslCert(), CaCert: req.GetCaCert()}
	resp, err := extReq.DoExtReq()

	if err != nil {
		return nil, err
	}

	respObj := &mappingpb.ExecuteApiResponse{}
	respObj.StatusCode = strconv.Itoa(resp.StatusCode)

	var responseHeaders []byte
	responseHeaders, _ = json.Marshal(resp.Header)
	respObj.ResponseHeaders = string(responseHeaders)

	var bodyBytes []byte
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	respObj.ResponseBody = string(bodyBytes)

	return respObj, nil
}

func (m *MappingServer) ResolveExpressions(ctx context.Context, req *mappingpb.ResolveExpressionsRequest) (*mappingpb.ResolveExpressionsResponse, error) {

	mappingResp, err := m.GetMapping(ctx, &mappingpb.MappingRequest{ProjectId: req.GetProjectId(), RepositoryIds: []string{req.GetRepositoryId()}, ProjectVersion: req.GetProjectVersion(), IncludeEnvs: true})
	if err != nil {
		return nil, err
	}

	repo, ok := mappingResp.GetConfig().Repositories[req.GetRepositoryId()]

	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("repository(%s) not found for project(%s)", req.GetRepositoryId(), req.GetProjectId()))
	}

	envVarMap := make(map[string]string)
	for _, envVar := range mappingResp.GetConfig().GetEnvironmentVariables() {
		for env, data := range envVar.GetEnvironments() {
			if env == mappingResp.GetConfig().DefaultEnvironment {
				envVarMap[constants.EnvVarPrefix+envVar.GetId()+constants.EnvVarSuffix] = data
			}
		}
	}

	url, err := utils.GetResolvedUrl(repo, envVarMap, req.GetPrompts())
	if err != nil {
		return nil, err
	}
	fmt.Println(url)

	headers, err := utils.GetResolvedHeaders(repo, envVarMap, req.GetPrompts())
	if err != nil {
		return nil, err
	}

	reqBody := make(map[string]string, 0)
	if strings.ToLower(repo.GetBody().GetType()) == constants.BodyTypeForm {
		for _, b := range repo.GetBody().GetForm() {
			reqBody[b.GetKey()] = b.GetValue()
		}
	} else if strings.ToLower(repo.GetBody().GetType()) == constants.BodyTypeJson {
		for _, b := range repo.GetBody().GetJson().GetFormKeys() {
			reqBody[b.GetKey()] = req.GetPrompts()[b.GetKey()]
		}
	}

	b, _ := json.Marshal(reqBody)

	exReq := &mappingpb.ExecuteApiRequest{
		HttpMethod: repo.GetHttpMethod(),
		Url:        url,
		Headers:    headers,
		Body:       string(b),
	}

	logger.WithContext(ctx).Infof(fmt.Sprintf("sending message to execute endpoint: %v", exReq))
	apiResp, err := m.ExecuteAPI(ctx, exReq)
	if err != nil {
		return nil, err
	}
	logger.WithContext(ctx).Infof(fmt.Sprintf("getting data back from execute endpoint: %v", apiResp))

	expressions := make(map[string]*expressionpb.Expression)
	b, _ = json.Marshal(repo.GetExpressions())
	json.Unmarshal(b, &expressions)

	expressionMenus := make([]*expressionpb.Menu, 0)
	b, _ = json.Marshal(req.ExpressionMenu)
	json.Unmarshal(b, &expressionMenus)
	PopulateGroupChildren(&expressionMenus, repo.GetExpressionMenu())

	expressionReq := &expressionpb.EvalExpressionRequest{ExpressionMenu: expressionMenus, Expressions: expressions, Data: apiResp.ResponseBody}
	logger.WithContext(ctx).Infof(fmt.Sprintf("sending message to expression service: %v", expressionReq))

	b, _ = m.Rabbit.PrepareRMQMessage(ctx, rabbitmq.Action_Resolve_Expression, expressionReq)

	corrId := request.GetContextRequestID(ctx)

	// publish to expression queue
	err = m.Rabbit.PublishWithCallBack(rabbitmq.Expression_Queue_Name, m.ReplyQueue.Name, corrId, b, m.Channel)
	if err != nil {
		return nil, fmt.Errorf("failed to publish message to expression queue - %v", err)
	}

	// consume message
	msg, err := m.Rabbit.ConsumeRPCMessage(ctx, m.ReplyQueue.Name, rabbitmq.Action_Resolve_Expression, corrId, m.Channel)
	if err != nil {
		return nil, err
	}
	logger.WithContext(ctx).Infof(fmt.Sprintf("message from expression: %v", msg))
	return &mappingpb.ResolveExpressionsResponse{
		ExpressionResult: msg.Payload.(string),
	}, nil
}

func (m *MappingServer) RestoreMapping(ctx context.Context, req *mappingpb.RestoreMappingRequest) (*mappingpb.RestoreMappingResponse, error) {

	// un-delete the mapping
	rows, err := m.DB.Model(&mappingpb.Repositories{}).Set("deleted_at = null").Where("project_id = ?", req.GetProjectId()).Update()
	if err != nil || rows.RowsAffected() == 0 {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to restore mapping with project id(%s), %v", req.GetProjectId(), err))
	}
	return &mappingpb.RestoreMappingResponse{Status: "Mapping restored successfully."}, nil
}
