package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/uibricks/studio-engine/internal/app/expression/constants"
)

// currentDate returns the current date
func currentDate(_ ...interface{}) (interface{}, error) {
	return time.Now().Format(constants.DATE_TIME_FORMAT), nil
}

// addToDate add days into  date
func addToDate(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	currentDateTime := args[0]
	offsetUnit := args[1]
	offset := args[2]

	if reflect.ValueOf(currentDateTime).Kind() != reflect.ValueOf(time.Now()).Kind() || reflect.ValueOf(offsetUnit).Kind() != reflect.String || reflect.ValueOf(offset).Kind() != reflect.Float64 {
		errMsg := fmt.Sprintf(constants.INVALID_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	switch strings.ToLower(offsetUnit.(string)) {

	case constants.DAY3_FORMAT_LOWER:
		return currentDateTime.(time.Time).AddDate(0,0, int(offset.(float64))).String(), nil
	case constants.MONTH3_FORMAT_LOWER:
		return currentDateTime.(time.Time).AddDate(0,int(offset.(float64)), 0).String(), nil
	case constants.YEAR2_FORMAT_LOWER:
		return currentDateTime.(time.Time).AddDate(int(offset.(float64)),0, 0).String(), nil
	case constants.HOUR_FORMAT_LOWER:
		t := time.Duration(int(offset.(float64)))
		return currentDateTime.(time.Time).Add(time.Hour * t).String(), nil
	case constants.MIN_FORMAT_IDENTIFIER:
		t := time.Duration(int(offset.(float64)))
		return currentDateTime.(time.Time).Add(time.Minute * t).String(), nil
	case constants.SEC_FORMAT_LOWER:
		t := time.Duration(int(offset.(float64)))
		return currentDateTime.(time.Time).Add(time.Second * t).String(), nil
	default:
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)

	}
}

func formatDate(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	input := args[0]
	layout := args[1]
	format := args[2]

	if reflect.ValueOf(input).Kind() != reflect.String || reflect.ValueOf(layout).Kind() != reflect.String || reflect.ValueOf(format).Kind() != reflect.String {
		errMsg := fmt.Sprintf(constants.INVALID_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	t, _ := time.Parse(ConvertFormat(layout.(string)), input.(string))
	goFormat := ConvertFormat(format.(string))
	resTime, _ := time.Parse(goFormat, t.Format(goFormat))
	return resTime, nil
}

// caseFunc is used to implement a switch case conditional
func caseFunc(args ...interface{}) (interface{}, error) {
	length := len(args)
	if (length-1)%2 != 0 || length == 1 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	i := 0
	for i < length-1 {
		if val, ok := args[i].(bool); val && ok {
			return args[i+1], nil
		}
		i += 2
	}
	return args[length-1], nil
}

// ifFunc is used to implement if else conditional
func ifFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}
	if val, ok := args[0].(bool); val && ok {
		return args[1], nil
	} else if ok {
		return args[2], nil
	}
	errMsg := fmt.Sprintf(constants.INVALID_TYPE_ERR, currentFuncName())
	return nil, fmt.Errorf(errMsg)
}

// concat is used to combine multiple string/array/both into one string/array
// limitation : dummy-first argument [it is being ignored],
// dummy-argument is being added to raw of the expression if concat() is encountered before go-evaluate invocation
// https://github.com/Knetic/govaluate/issues/91
// https://github.com/Knetic/govaluate/issues/89
//
// Example 1 :
// args[0] is a dummy argument => ''
// args[1] is a slice of strings => ['Amb. Sheela Verma', 'Mr. Ranjit Verma', 'Girika Verma']
// args[2] is a string => " - "
// args[3] is a slice of strings => ['Male', 'Male', 'Female']
// result => ['Amb. Sheela Verma - Male', 'Mr. Ranjit Verma - Male', 'Girika Verma - Female']
//
// Example 2 :
// args[0] is a dummy argument => ''
// args[1] is a string => "This is an"
// args[2] is a string => " example."
// result => "This is an example."
func concat(args ...interface{}) (interface{}, error) {

	// ignoring first argument here
	args = append(args[:0], args[1:]...)
	length := len(args)

	if length < 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	var result interface{}
	for i := range args {
		if i > 0 {
			result = rConcat(result, args[i])
		} else {
			result = args[i]
		}
	}
	return result, nil
}

// rConcat is a helper function for concat
func rConcat(res interface{}, args interface{}) interface{} {

	if reflect.ValueOf(args).Kind() == reflect.Slice && reflect.ValueOf(res).Kind() == reflect.Slice {
		d1 := reflect.ValueOf(args)
		d2 := reflect.ValueOf(res)
		returnSlice := make([]interface{}, d2.Len())
		for i := 0; i < d2.Len(); i++ {
			returnSlice[i] = rConcat(d2.Index(i).Interface(), d1.Index(i).Interface())
		}
		return returnSlice
	} else if reflect.ValueOf(res).Kind() == reflect.Slice {
		d := reflect.ValueOf(res)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i] = rConcat(v, args)
		}
		return returnSlice
	} else if reflect.ValueOf(args).Kind() == reflect.Slice {
		d := reflect.ValueOf(args)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i] = rConcat(res, v)
		}
		return returnSlice
	} else {
		return fmt.Sprintf("%s", res) + fmt.Sprintf("%s", args)
	}
}

// concatDelimiter mimics go strings package's 'join' function
//
// Example:
// args[0] is a string separator => ','
// args[1] is a slice of strings => ['None', 'Glazed Sugar', 'Powdered Sugar', 'Chocolate with Sprinkles', 'Chocolate', 'Maple']
// result => 'None,Glazed,Sugar,Powdered Sugar,Chocolate with Sprinkles,Chocolate,Maple'
func concatDelimiter(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length < 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if reflect.ValueOf(args[0]).Kind() != reflect.String {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_TYPE_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	sep := args[0].(string)
	var result interface{}
	for i:=1; i<len(args); i++ {
		if i-1 > 0 {
			args[i] = rConcat(sep, args[i])
			result = rConcat(result, args[i])
		} else {
			result = args[i]
		}
	}
	return result, nil
}

// regex is used to implement regex function - for both array and strings
//
// Example 1 :
// args[0] is pattern => '^D'
// args[1] is slice of strings => ['Amb. Sheela Verma', 'Mr. Ranjit Verma', 'Dhanpati Verma']
// result => [false, false, true]
//
// Example 2 :
// args[0] is pattern => '^D'
// args[1] is a string => 'Dhanpati Verma'
// result => true
func regex(args ...interface{}) (interface{}, error) {

	if len(args) != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	pattern := args[0]

	if reflect.ValueOf(args[1]).Kind() == reflect.Slice {
		resArr := make([]interface{}, 0)
		for _, val := range args[1].([]interface{}) {
			res, err := rRegex(pattern.(string), val)
			if err != nil {
				resArr = append(resArr, false)
			}
			resArr = append(resArr, res)
		}
		return resArr, nil
	}

	r, err := regexp.Compile(pattern.(string))
	if err != nil {
		return nil, err
	}

	return r.MatchString(args[1].(string)), nil
}

// rRegex is a helper function for regex
func rRegex(pattern string, args interface{}) (interface{}, error) {

	if reflect.ValueOf(args).Kind() == reflect.Slice {
		d := reflect.ValueOf(args)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		var err error
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i], err = rRegex(pattern, v)
		}
		if err != nil {
			return nil, err
		}
		return returnSlice, err
	} else if reflect.ValueOf(args).Kind() == reflect.Map {
		d := reflect.ValueOf(args)
		tmpData := make(map[string]interface{})
		var err error
		for _, k := range d.MapKeys() {
			typeOfValue := reflect.TypeOf(d.MapIndex(k).Interface()).Kind()
			if typeOfValue == reflect.Map || typeOfValue == reflect.Slice {
				tmpData[k.String()], err = rRegex(pattern, d.MapIndex(k).Interface())
			} else {
				tmpData[k.String()] = d.MapIndex(k).Interface()
			}
		}
		if err != nil {
			return nil, err
		}
		return tmpData, nil
	} else {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		return r.MatchString(fmt.Sprintf("%s", args)), nil
	}
}

// trim function trims the string with space
//
// Example 1 :
// args is slice of strings => [' Amb. Sheela Verma ', ' Mr. Ranjit Verma', 'Dhanpati Verma ']
// result => ['Amb. Sheela Verma', 'Mr. Ranjit Verma', 'Dhanpati Verma']
//
// Example 2 :
// args is a string => ' Dhanpati Verma '
// result => 'Dhanpati Verma'
func trim(args ...interface{}) (interface{}, error) {
	length := len(args)

	if length == 0 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	if length > 1 {
		resArr := make([]interface{}, 0)
		for _, val := range args {
			res, err := rTrim(val)
			if err != nil {
				return nil, err
			} else {
				resArr = append(resArr, res)
			}
		}
		return resArr, nil
	}

	originalString := fmt.Sprintf("%v", args[0])
	trimmedString := strings.Trim(originalString, " ")

	return trimmedString, nil
}

// rTrim is a helper function for trim
func rTrim(args interface{}) (interface{}, error) {

	if reflect.ValueOf(args).Kind() == reflect.Slice {
		d := reflect.ValueOf(args)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		var err error
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i], err = rTrim(v)
		}
		if err != nil {
			return nil, err
		}
		return returnSlice, err
	} else if reflect.ValueOf(args).Kind() == reflect.Map {
		d := reflect.ValueOf(args)
		tmpData := make(map[string]interface{})
		var err error
		for _, k := range d.MapKeys() {
			typeOfValue := reflect.TypeOf(d.MapIndex(k).Interface()).Kind()
			if typeOfValue == reflect.Map || typeOfValue == reflect.Slice {
				tmpData[k.String()], err = rTrim(d.MapIndex(k).Interface())
			} else {
				tmpData[k.String()] = d.MapIndex(k).Interface()
			}
		}
		if err != nil {
			return nil, err
		}
		return tmpData, nil
	} else {
		return strings.Trim(fmt.Sprintf("%s", args), " "), nil
	}
}

func startsWith(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	pattern := args[0]

	if reflect.ValueOf(args[1]).Kind() == reflect.Slice {
		resArr := make([]interface{}, 0)
		for _, val := range args[1].([]interface{}) {
			res, err := rStartsWith(pattern.(string), val)
			if err != nil {
				resArr = append(resArr, false)
			}
			resArr = append(resArr, res)

		}
		return resArr, nil
	}

	return strings.HasPrefix(args[1].(string), pattern.(string)), nil
}

func rStartsWith(pattern string, args interface{}) (interface{}, error) {
	fmt.Println(pattern, args)

	if reflect.ValueOf(args).Kind() == reflect.Slice {
		d := reflect.ValueOf(args)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		var err error
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i], err = rStartsWith(pattern, v)
		}
		if err != nil {
			return nil, err
		}
		return returnSlice, err
	} else {
		return strings.HasPrefix(args.(string), pattern), nil
	}
}

func endsWith(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		errMsg := fmt.Sprintf(constants.INVALID_PARAM_ERR, currentFuncName())
		return nil, fmt.Errorf(errMsg)
	}

	pattern := args[0]

	if reflect.ValueOf(args[1]).Kind() == reflect.Slice {
		resArr := make([]interface{}, 0)
		for _, val := range args[1].([]interface{}) {
			res, err := rEndsWith(pattern.(string), val)
			if err != nil {
				resArr = append(resArr, false)
			}
			resArr = append(resArr, res)
		}
		return resArr, nil
	}

	return strings.HasSuffix(args[1].(string), pattern.(string)), nil
}

func rEndsWith(pattern string, args interface{}) (interface{}, error) {

	if reflect.ValueOf(args).Kind() == reflect.Slice {
		d := reflect.ValueOf(args)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		var err error
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i], err = rEndsWith(pattern, v)
		}
		if err != nil {
			return nil, err
		}
		return returnSlice, err
	} else {
		return strings.HasSuffix(args.(string), pattern), nil
	}
}