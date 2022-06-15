package utils

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// FindInArray - finds given string in the pointer array
func FindInArray(strArr *[]string, key string) (int, bool) {
	for i, str := range *strArr {
		if str == key {
			return i, true
		}
	}
	return -1, false
}

// StrPointerArrToStrArr - convers arr of pointers to a string array
func StrPointerArrToStrArr(arr []*string) []string {
	strArr := []string{}
	for _, str := range arr {
		strArr = append(strArr, *str)
	}
	return strArr
}

func MarshalToString(pb proto.Message) (string, error) {
	m := jsonpb.Marshaler{EmitDefaults: true}
	res, err := m.MarshalToString(pb)
	return res, err
}
