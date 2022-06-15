package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/uibricks/studio-engine/internal/app/mapping/constants"
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	mappingpb "github.com/uibricks/studio-engine/internal/pkg/proto/mapping"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"net/url"
	"strings"
)

func GetEmptyConfig() *mappingpb.Config {
	return &mappingpb.Config{
		Repositories: map[string]*mappingpb.Repository{},
	}
}

func CastStringToMapping(val string) (*mappingpb.Repositories, error) {
	mapping := &mappingpb.Repositories{
		Config: GetEmptyConfig(),
	}
	if len(val) > 0 {
		err := jsonpb.UnmarshalString(val, mapping)
		if err != nil {
			return nil, err
		}
	}

	if mapping.GetConfig().Repositories == nil {
		mapping.GetConfig().Repositories = make(map[string]*mappingpb.Repository)
	}

	return mapping, nil
}

func MarshalToString(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{}
	res, err := m.MarshalToString(pb)
	return res, err
}

func UpdateRepoNames(repos map[string]string, menu []*mappingpb.Menu) {
	if menu != nil {
		for _, item := range menu {
			if _, found := repos[item.Id]; found {
				repos[item.Id] = item.Name
			}
			if item.Children != nil && len(item.Children) > 0 {
				UpdateRepoNames(repos, item.Children)
			}
		}
	}
}

func MapToArrayDependencies(dependencies map[string]string) []*mappingpb.Dependency {
	arr := make([]*mappingpb.Dependency, 0)
	for key, val := range dependencies {
		arr = append(arr, &mappingpb.Dependency{
			Id:   key,
			Name: val,
		})
	}
	return arr
}

func DeleteRepoFromMenu(repoId string, menu []*mappingpb.Menu, deleted *bool) []*mappingpb.Menu {
	if menu != nil && len(repoId) > 0 {
		for i, item := range menu {
			if item.Id == repoId {
				menu = append(menu[:i], menu[(i+1):]...)
				*deleted = true
				return menu
			}
			item.Children = DeleteRepoFromMenu(repoId, item.Children, deleted)
			if *deleted {
				break
			}
		}
	}
	return menu
}

func GetResolvedUrl(repo *mappingpb.Repository, envVarMap, prompts map[string]string) (string, error) {

	url := repo.GetUrl()
	for _, qp := range repo.GetQueryParams() {

		if !qp.GetChecked() {
			continue
		}

		if qp.GetPrompt() {
			if val, ok := prompts[qp.GetKey()]; ok {
				if qp.GetIncludeInPath() {
					if !strings.HasSuffix(url, constants.UrlSuffix) {
						url += constants.UrlSuffix
					}
					url += val
					continue
				}
				if !strings.Contains(url, constants.QueryParamsSeparator) {
					url += constants.QueryParamsSeparator + qp.GetKey() + constants.QueryParamValueSeparator + val
				} else {
					url += constants.QueryParamAppend + qp.GetKey() + constants.QueryParamValueSeparator + val
				}
			} else {
				return "", status.Error(codes.InvalidArgument, fmt.Sprintf("no value for url prompt : %s", qp.GetKey()))
			}
		} else {
			if qp.GetIncludeInPath() {
				if !strings.HasSuffix(url, constants.UrlSuffix) {
					url += constants.UrlSuffix
				}
				url += qp.GetValue()
				continue
			}
			if !strings.Contains(url, constants.QueryParamsSeparator) {
				url += constants.QueryParamsSeparator + qp.GetKey() + constants.QueryParamValueSeparator + qp.GetValue()
			} else {
				url += constants.QueryParamAppend + qp.GetKey() + constants.QueryParamValueSeparator + qp.GetValue()
			}
		}

	}

	for k,v := range envVarMap {
		if strings.Contains(url, k) {
			url = strings.Replace(url, k, v, -1)
		}
	}

	if strings.Contains(url, constants.EnvVarPrefix) {
		return "", status.Error(codes.InvalidArgument , "missing environment variables for url")
	}

	return url, nil
}

func GetResolvedHeaders(repo *mappingpb.Repository, envVarMap, prompts map[string]string) ([]*mappingpb.Headers, error) {
	resHeaders := repo.GetHeaders()

	for _, hdr := range repo.GetHeaders() {
		if !hdr.GetChecked() {
			continue
		}

		if hdr.GetPrompt() {
			if val, ok := prompts[hdr.GetKey()]; ok {
				hdr.Value = val
			} else {
				return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("no value for header prompt : %s", hdr.GetKey()))
			}
		} else {
			if strings.Contains(hdr.GetValue(), constants.EnvVarPrefix) && strings.Contains(hdr.GetValue(), constants.EnvVarSuffix) {
				if val, ok := envVarMap[hdr.GetValue()]; ok {
					hdr.Value = val
				} else {
					return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("no environment variable for header : %s", hdr.GetKey()))
				}
			}
		}
	}
	return resHeaders, nil
}

func ResponseStringToJSON(v interface{}) (interface{}, error) {

	var dResp interface{}
	err := json.Unmarshal([]byte(v.(string)), &dResp)
	if err != nil {
		return nil, err
	}

	return dResp, nil

}

type ExternalRequest struct {
	Name    string                 `json:"name"`
	URL     string                 `json:"url"`
	Type    string                 `json:"type"`
	Headers map[string]string      `json:"headers"`
	Params  map[string]string      `json:"params"`
	Body    map[string]interface{} `json:"body"`
	SslCert []string `json:"sslCert"`
	CaCert string `json:"caCert"`

	ReqCtx context.Context
}

func (r *ExternalRequest) DoExtReq() (*http.Response, error) {

	body := strings.NewReader(r.getRequestBody())
	req, err := http.NewRequest(r.Type, r.URL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range r.Headers {
		req.Header.Add(k, v)
	}

	q := req.URL.Query()
	for k, v := range r.Params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	var client *http.Client

	if len(r.SslCert) == 2 || r.CaCert != "" {
		client = r.GetHTTPClient()
	} else {
		transport := &http.Transport{TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		}}
		client = &http.Client{Transport: transport}
		//client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Sugar.Errorf("failed to execute api : - %v", err)
	}

	return resp, err
}

func (r *ExternalRequest) getRequestBody() string {

	var contentType string
	for k, v := range r.Headers {
		if strings.ToLower(k) == constants.HeaderKeyContentType {
			contentType = strings.ToLower(v)
			break
		}
	}

	switch contentType {
	case constants.ContentTypeJson:
		j, _ := json.Marshal(r.Body)
		return string(j)
	default:
		kvs := getQueryKVs(r.Body, "", false)
		return strings.Join(kvs, "&")
	}
}


func getQueryKVs(bMap map[string]interface{}, parentKey string, numericKeys bool) []string {
	var kvs = make([]string, 0)
	for key, value := range bMap {
		key = url.QueryEscape(key)
		if parentKey != "" {
			if numericKeys {
				key = fmt.Sprintf("%s[]", parentKey)
			} else {
				key = fmt.Sprintf("%s[%s]", parentKey, key)
			}
		}

		switch value.(type) {
		case []interface{}:
			tMap := make(map[string]interface{})
			for i, tVal := range value.([]interface{}) {
				tMap[fmt.Sprintf("%d", i)] = tVal
			}
			kvs = append(kvs, getQueryKVs(tMap, key, true)...)
		case map[string]interface{}:
			kvs = append(kvs, getQueryKVs(value.(map[string]interface{}), key, false)...)
		default:
			valStr := url.QueryEscape(fmt.Sprintf("%v", value))
			kvs = append(kvs, fmt.Sprintf("%s=%s", key, valStr))
		}
	}

	return kvs
}

func (r *ExternalRequest) GetHTTPClient() *http.Client {
	if (r.SslCert[0] == "" || r.SslCert[1] == "") && r.CaCert == "" {
		return http.DefaultClient
	}

	var tlsConfig = new(tls.Config)

	if r.SslCert[0] != "" && r.SslCert[1] != "" {
		cert, _ := tls.X509KeyPair([]byte(r.SslCert[0]), []byte(r.SslCert[1]))
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if r.CaCert != "" {
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM([]byte(r.CaCert))
		tlsConfig.RootCAs = caCertPool
		//tlsConfig.InsecureSkipVerify = true
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return &http.Client{Transport: transport}
}
