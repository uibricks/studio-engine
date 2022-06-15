package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/uibricks/studio-engine/internal/app/expression/constants"
)

// contains mimics go strings package's contains functionality
func contains(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	switch strings.ToLower(args[0].(string)) {
	case constants.STRING_OP_CONTAINS:
		return strings.Contains(args[1].(string), args[2].(string)), nil

	case constants.STRING_OP_CONTAINS_ANY:
		return strings.ContainsAny(args[1].(string),  args[2].(string)), nil
	}

	return nil, fmt.Errorf(fmt.Sprintf(constants.INVALID_FUNCTION_IN_PARAM, args[0].(string), "Contains()"))
}

// count mimics go strings package's count functionality
func count(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	return strings.Count(args[0].(string), args[1].(string)), nil
}

// mimics go strings package's equal functionality
func equal(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	switch strings.ToLower(args[0].(string)) {
	case constants.STRING_OP_EQUAL_CS:
		return args[1].(string) == args[2].(string), nil
	case constants.STRING_OP_EQUAL_CIS:
		return strings.EqualFold(args[1].(string), args[2].(string)), nil
	}
	return nil, fmt.Errorf(fmt.Sprintf(constants.INVALID_FUNCTION_IN_PARAM, args[0].(string), "Equal()"))
}

// mimics go strings package's fields functionality
func fields(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 1 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	res := strings.Fields(args[0].(string))
	resArr := make([]interface{}, len(res))
	for i, v := range res {
		resArr[i] = v
	}
	return resArr, nil
}

// mimics go strings package's index functionality
func index(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	switch strings.ToLower(args[0].(string)) {
	case constants.STRING_OP_INDEX_ALL:
		resArr := make([]interface{}, 0)

		fullStr := args[1].(string)
		res := strings.Index(fullStr, args[2].(string))
		resArr = append(resArr, res)
		for ;res != -1; {
			fullStr = fullStr[res+1:]
			res = strings.Index(fullStr, args[2].(string))
			resArr = append(resArr, resArr[len(resArr)-1].(int) + res + 1)
		}
		resArr = resArr[:len(resArr)-1]

		return  resArr, nil
	case constants.STRING_OP_INDEX_LAST:
		return strings.LastIndex(args[1].(string), args[2].(string)), nil
	default:
		return strings.Index(args[1].(string), args[2].(string)), nil
	}
}

// mimics go strings package's repeat functionality
func repeat(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if !(reflect.ValueOf(args[0]).Kind() == reflect.String) || !(reflect.ValueOf(args[1]).Kind() == reflect.Float64) {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	return strings.Repeat(args[0].(string), int(args[1].(float64))), nil
}

// replace mimics go strings package's replace functionality
func replace(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 4 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if !(reflect.ValueOf(args[length-1]).Kind() == reflect.Float64) {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	return strings.Replace(args[0].(string), args[1].(string), args[2].(string), int(args[3].(float64))), nil
}

// split mimics go strings package's split functionality
func split(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	res := strings.Split(args[0].(string), args[1].(string))
	resArr := make([]interface{}, len(res))
	for i, v := range res {
		resArr[i] = v
	}
	return resArr, nil
}

// stringTransform mimics go strings package's toLower/toUpper/toTitle functionalities
func stringTransform(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	switch strings.ToLower(args[0].(string)) {
	case constants.STRING_OP_LOWER:
		return strings.ToLower(args[1].(string)), nil
	case constants.STRING_OP_UPPER:
		return strings.ToUpper(args[1].(string)), nil
	case constants.STRING_OP_TITLE:
		return strings.Title(args[1].(string)), nil
	}

	return nil, fmt.Errorf(fmt.Sprintf(constants.INVALID_FUNCTION_IN_PARAM, args[0].(string), currentFuncName()))
}