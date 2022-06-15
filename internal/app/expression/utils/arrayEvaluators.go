package utils

import (
	"fmt"
	"github.com/uibricks/studio-engine/internal/app/expression/constants"
	"reflect"
	"sort"
)

// extend combines multiple arrays
//
// Example 1 :
// args[0] is a string => 'Amb. Sheela Verma'
// args[1] is a string => 'Mr. Ranjit Verma'
// args[2] is a slice of strings => ['Male', 'Male']
// result => ['Amb. Sheela Verma', 'Mr. Ranjit Verma', 'Male', 'Male']
//
// Example 2 :
// args[0] is a slice of strings => ['A', 'B']
// args[1] is a is a slice of strings => ['C', 'D']
// result => ['A', 'B', 'C', 'D']
func extend(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length < 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	resArr := make([]interface{}, 0)
	for _, val := range args {
		if reflect.ValueOf(val).Kind() == reflect.Slice {
			resArr = append(resArr, val.([]interface{})...)
		} else {
			resArr = append(resArr, val)
		}
	}
	return resArr, nil
}

// appendToArray mimics array-append functionality
//
// Example:
// args[0] is a number => 10
// args[1] is a slice of numbers => [5,6,7]
// result => [5,6,7,10]
func appendToArray(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if !(reflect.ValueOf(args[1]).Kind() == reflect.Slice) {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	args[1] = append(args[1].([]interface{}), args[0])
	return args[1], nil
}

// countArrayElements returns length of input array
//
// Example:
// args is a slice => [1,2,3.4,'data']
// result => 4
func countArrayElements(args ...interface{}) (interface{}, error) {
	return len(args), nil
}

// indexArray returns position of first-argument [element] in second-argument [array]
func indexArray(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if reflect.ValueOf(args[1]).Kind() != reflect.Slice {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	for i, val := range args[1].([]interface{}) {
		if val == args[0] {
			return i, nil
		}
	}
	return nil, fmt.Errorf(fmt.Sprintf(constants.NOT_FOUND, currentFuncName()))
}

// insert inserts given element [1st argument] at given position [2nd argument] into array [3rd argument]
func insert(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if reflect.ValueOf(args[1]).Kind() != reflect.Float64 || reflect.ValueOf(args[2]).Kind() != reflect.Slice {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	index := int(args[1].(float64))
	arr := args[2].([]interface{})
	arr = append(arr, 0)
	copy(arr[index+1:], arr[index:])
	arr[index] = args[0]

	return arr, nil
}

// pop removes element from array [2nd argument] at given position [1st argument]
func pop(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if !(reflect.ValueOf(args[1]).Kind() == reflect.Slice) || !(reflect.ValueOf(args[0]).Kind() == reflect.Float64) {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	index := int(args[0].(float64))
	arr := args[1].([]interface{})
	copy(arr[index:], arr[index+1:])
	arr[len(arr)-1] = ""
	arr = arr[:len(arr)-1]

	return arr, nil
}

// remove removes given element [1st argument] from array [2nd argument] by its value
func remove(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if reflect.ValueOf(args[1]).Kind() != reflect.Slice || reflect.ValueOf(args[0]).Kind() != reflect.String {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	index := -1
	arr := args[1].([]interface{})
	for i, val := range arr {
		if val == args[0] {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, fmt.Errorf(fmt.Sprintf(constants.NOT_FOUND, currentFuncName()))
	}
	return append(arr[:index], arr[index+1:]...), nil
}

// reverse reverses input array
func reverse(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length == 0 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	for i, j := 0, len(args)-1; i < j; i, j = i+1, j-1 {
		args[i], args[j] = args[j], args[i]
	}

	return args, nil
}

// getAt gets the element at given index
func getAt(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if reflect.ValueOf(args[1]).Kind() != reflect.Slice || reflect.ValueOf(args[0]).Kind() != reflect.Float64 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	return args[1].([]interface{})[int(args[0].(float64))], nil
}

// sortArray sorts string/int/float arrays
func sortArray(args ...interface{}) (interface{}, error) {
	length := len(args)

	if reflect.ValueOf(args[0]).Kind() == reflect.String {
		arr := make([]string, length)
		resArr := make([]interface{}, length)
		for i, val := range args {
			arr[i] = val.(string)
		}
		sort.Strings(arr)
		for i, val := range arr {
			resArr[i] = val
		}
		return resArr, nil
	} else if reflect.ValueOf(args[0]).Kind() == reflect.Float64 {
		arr := make([]float64, length)
		resArr := make([]interface{}, length)
		for i, val := range args {
			arr[i] = val.(float64)
		}
		sort.Float64s(arr)
		for i, val := range arr {
			resArr[i] = val
		}
		return resArr, nil
	}

	return args, nil
}
