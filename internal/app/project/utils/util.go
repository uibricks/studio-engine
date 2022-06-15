package utils

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	projectpb "github.com/uibricks/studio-engine/internal/pkg/proto/project"
	"github.com/uibricks/studio-engine/internal/pkg/utils"
)

func RetainLatestProjects(rows []*projectpb.Object) []*projectpb.Object {
	retainedProjects := make([]*projectpb.Object, 0)
	projectNames := make([]string, 0)

	for i := 0; i < len(rows); i=i+1 {
		if idx, ok := utils.FindInArray(&projectNames, rows[i].GetLuid()); ok {
			retainedProjects[idx] = rows[i]
		} else {
			retainedProjects = append(retainedProjects, rows[i])
			projectNames = append(projectNames, rows[i].GetLuid())
		}
	}
	return retainedProjects
}

func MarshalToString(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{}
	res, err := m.MarshalToString(pb)
	return res, err
}

func UpdateCompNames(comps map[string]string, components []*projectpb.Component) {
	if components != nil {
		for _, item := range components {
			if _, found := comps[item.Id]; found {
				comps[item.Id] = item.Name
			}
			if item.Children != nil && len(item.Children) > 0 {
				UpdateCompNames(comps, item.Children)
			}
		}
	}
}

func ConvertComponentMapToArr(componentMap map[string]*projectpb.Component) []*projectpb.Component {
	componentArr := make([]*projectpb.Component, 0)

	for _,v := range componentMap {
		componentArr = append(componentArr, v)
	}

	return componentArr
}

func ConvertComponentArrToMap(componentArr []*projectpb.Component) map[string]*projectpb.Component {
	componentMap := make(map[string]*projectpb.Component)

	for _,v := range componentArr {
		componentMap[v.Id] = v
	}

	return componentMap
}

func MapToArrayDependencies(dependencies map[string]string) []*projectpb.ComponentDependency {
	arr := make([]*projectpb.ComponentDependency, 0)
	for key, val := range dependencies {
		arr = append(arr, &projectpb.ComponentDependency{
			Id:   key,
			Name: val,
		})
	}
	return arr
}

func DeleteComponentFromComponentArr(compId string, components []*projectpb.Component, deleted *bool) []*projectpb.Component {
	if components != nil && len(compId) > 0 {
		for i, item := range components {
			if item.Id == compId {
				components = append(components[:i], components[(i+1):]...)
				*deleted = true
				return components
			}
			item.Children = DeleteComponentFromComponentArr(compId, item.Children, deleted)
			if *deleted {
				break
			}
		}
	}
	return components
}

