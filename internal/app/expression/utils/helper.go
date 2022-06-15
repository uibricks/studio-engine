package utils

import (
	"fmt"
	"github.com/uibricks/studio-engine/internal/app/expression/constants"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

// generates an array of all possible case-combinations of a given string
func permutations(str string) []string {
	n := len(str)
	arr := make([]string, 0)
	str = strings.ToLower(str)
	max := 1 << n
	for i := 0; i < max; i++ {
		combination := make([]string, n)
		for j := range str {
			combination[j] = string(str[j])
		}
		for k := 0; k < n; k++ {
			if ((i >> k) & 1) == 1 {
				combination[k] = strings.ToUpper(string(str[k]))
			}
		}
		temp := ""
		for m := range combination {
			temp = temp + combination[m]
		}
		arr = append(arr, temp)

	}
	return arr
}

// adds additional dummy first argument to concat function
func addExtraArgToConcat(raw string) string {
	if strings.Contains(strings.ToLower(raw), "concat(") {
		re := regexp.MustCompile("(?i)concat" + "\\(")
		raw = re.ReplaceAllString(raw, "concat('', ")
	}
	return raw
}

// wrapDataWithKey converts array of data to map of objects with given key
func wrapDataWithKey(data interface{}, key string) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		d := reflect.ValueOf(data)
		tempData := make([]interface{}, 0)
		for i := 0; i < d.Len(); i++ {
			res := wrapDataWithKey(d.Index(i).Interface(), key)
			tempData = append(tempData, res)
		}
		return tempData
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		d := reflect.ValueOf(data)
		tempData := make(map[string]interface{})
		for _, k := range d.MapKeys() {
			tempData[k.String()] = wrapDataWithKey(d.MapIndex(k).Interface(), key)
		}
		return tempData
	} else {
		return map[string]interface{}{key: data}
	}
}

// mergeObjs converts two different array of objects to a single array by merging each objects at individual levels
func mergeObjs(obj1, obj2 interface{}) interface{} {

	d1 := reflect.ValueOf(obj1)
	d2 := reflect.ValueOf(obj2)
	if d1.Kind() == d2.Kind() && d2.Kind() == reflect.Map && len(d1.MapKeys()) == len(d2.MapKeys()) { //case - two fields of same level of a nested structure in nested-map format
		tempData := make(map[string]interface{})

		for i, k := range d1.MapKeys() {
			if d2.MapIndex(k).Kind() != reflect.Invalid {
				res := mergeObjs(d1.MapIndex(k).Interface(), d2.MapIndex(k).Interface())
				tempData[fmt.Sprintf("%s", k)] = res
			} else {
				for j, k2 := range d2.MapKeys() {
					if i == j {
						tempData[fmt.Sprintf("%s", k)] = d1.MapIndex(k).Interface()
						tempData[fmt.Sprintf("%s", k2)] = d2.MapIndex(k2).Interface()
						break
					}
				}
			}
		}
		return tempData
	} else if d1.Kind() == d2.Kind() && d2.Kind() == reflect.Slice && d1.Len() == d2.Len() { //case - two fields of same level of a nested structure in nested-slice format
		tempData := make([]interface{}, d1.Len())
		for i := 0; i < d1.Len(); i++ {
			currRes := mergeObjs(d1.Index(i).Interface(), d2.Index(i).Interface())
			tempData[i] = currRes
		}
		return tempData
	} else {
		var res interface{}
		if d1.Kind() != d2.Kind() { //case - two fields of diff. level of a nested structure of diff. formats
			if d1.Kind() == reflect.Slice {
				res = mergeUnequals(obj1, obj2)
			} else if d2.Kind() == reflect.Slice {
				res = mergeUnequals(obj2, obj1)
			}
			return res
		} else if d1.Kind() == d2.Kind() && d2.Kind() == reflect.Slice && d1.Len() != d2.Len() { //case - two fields of diff. level of a nested structure in nested-slice format
			var bigObjArr reflect.Value
			var smallObjArr reflect.Value
			if d1.Len() > d2.Len() {
				bigObjArr = d1
				smallObjArr = d2
			} else {
				bigObjArr = d2
				smallObjArr = d1
			}

			tempDataArr := make([]interface{}, bigObjArr.Len())

			for i:=0; i<bigObjArr.Len(); i++ {
				if i < smallObjArr.Len() {
					tempDataArr[i] = mergeObjs(bigObjArr.Index(i).Interface(), smallObjArr.Index(i).Interface())
				} else {
					tempDataArr[i] = bigObjArr.Index(i).Interface()
				}
			}
			return tempDataArr
		} else if d1.Kind() == d2.Kind() && d2.Kind() == reflect.Map && len(d1.MapKeys()) != len(d2.MapKeys()) { //case - two fields of diff. level of a nested structure in nested-map format
			tempData := make(map[string]interface{})
			for _,k := range d1.MapKeys() {
				tempData[k.String()] = d1.MapIndex(k).Interface()
			}
			for _,k := range d2.MapKeys() {
				tempData[k.String()] = d2.MapIndex(k).Interface()
			}
			return tempData
		}
	}
	return nil
}

func mergeUnequals(obj1 interface{}, obj2 interface{}) interface{} {

	if reflect.ValueOf(obj1).Kind() == reflect.Slice {
		d1 := reflect.ValueOf(obj1)
		tempDataArr := make([]interface{}, 0)
		for i:=0; i< d1.Len(); i++ {
			typeOfVal := reflect.ValueOf(d1.Index(i).Interface()).Kind()
			if typeOfVal == reflect.Slice {
				res := mergeUnequals(d1.Index(i).Interface(), obj2)
				tempDataArr = append(tempDataArr, res)
			} else if typeOfVal == reflect.Map {
				d2 := reflect.ValueOf(obj2)
				d1Maps := reflect.ValueOf(d1.Index(i).Interface())
				tempData := make(map[string]interface{})
				if d1Maps.Kind() == reflect.Map {
					for _, k := range d1Maps.MapKeys() {
						tempData[k.String()] = d1Maps.MapIndex(k).Interface()
					}
				}
				if d2.Kind() == reflect.Map {
					for _, k := range d2.MapKeys() {
						tempData[k.String()] = d2.MapIndex(k).Interface()
					}
				}
				tempDataArr = append(tempDataArr, tempData)
			}
		}
		return tempDataArr
	}
	return nil
}

func getNestedDepth(elem interface{}) int {
	d := reflect.ValueOf(elem)
	if d.Kind() == reflect.Slice && d.Len() > 0 {
		var max int
		tempArr := make([]int, d.Len())
		for i:=0; i<d.Len(); i++ {
			tempArr[i] = 1 + getNestedDepth(d.Index(i).Interface())
		}
		max = tempArr[0]
		for i := range tempArr {
			if tempArr[i] > max {
				max = tempArr[i]
			}
		}
		return max
	} else if d.Kind() == reflect.Map {
		var max int
		tempArr := make(map[string]int)
		for _,v := range d.MapKeys() {
			tempArr[v.String()] = 1 + getNestedDepth(d.MapIndex(v).Interface())
		}
		max = -1
		for i,_ := range tempArr {
			if tempArr[i] > max {
				max = tempArr[i]
			}
		}
		return max
	} else {
		return 0
	}
}

// flatten converts a nested structure to a flat structure
func flatten(data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		d := reflect.ValueOf(data)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, 0)
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for _, v := range tmpData {
			res := flatten(v)
			if reflect.ValueOf(res).Kind() == reflect.Slice {
				resArr := res.([]interface{})
				for j := range resArr {
					returnSlice = append(returnSlice, resArr[j])
				}
			} else {
				returnSlice = append(returnSlice, res)
			}
		}
		return returnSlice
	} else if reflect.ValueOf(data).Kind() == reflect.Map {
		d := reflect.ValueOf(data)
		tmpData := make(map[string]interface{})
		for _, k := range d.MapKeys() {
			typeOfValue := reflect.TypeOf(d.MapIndex(k).Interface()).Kind()
			if typeOfValue == reflect.Map || typeOfValue == reflect.Slice {
				tmpData[k.String()] = flatten(d.MapIndex(k).Interface())
			} else {
				tmpData[k.String()] = d.MapIndex(k).Interface()
			}

		}
		return tmpData
	}
	return data
}

// FindInArray checks if a string is present in the string-array or not
// common-module's 'FindInArray' is not being referenced because GopherJS encounters errors with some of its dependencies
func FindInArray(strArr []string, key string) (int, bool) {
	for i, str := range strArr {
		if str == key {
			return i, true
		}
	}
	return -1, false
}

// currentFuncName get caller's function name from stacktrace
func currentFuncName() string {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	parts := strings.Split(runtime.FuncForPC(pc[0]).Name(), "/")
	return strings.Title(strings.Split(parts[len(parts) - 1], ".")[1]) + "()"
}

func ConvertFormat(format string) string {
	goFormat := format
	if strings.Contains(goFormat, constants.YEAR1_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.YEAR1_FORMAT_UPPER, constants.GO_YYYY, -1)
	} else if strings.Contains(goFormat, constants.YEAR1_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.YEAR1_FORMAT_LOWER, constants.GO_YYYY, -1)
	} else if strings.Contains(goFormat, constants.YEAR2_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.YEAR2_FORMAT_UPPER, constants.GO_YY, -1)
	} else if strings.Contains(goFormat, constants.YEAR2_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.YEAR2_FORMAT_LOWER, constants.GO_YY, -1)
	}

	if strings.Contains(goFormat, constants.MONTH1_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.MONTH1_FORMAT_UPPER, constants.GO_MMMM, -1)
	} else if strings.Contains(goFormat, constants.MONTH1_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.MONTH1_FORMAT_LOWER, constants.GO_MMMM, -1)
	} else if strings.Contains(goFormat, constants.MONTH2_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.MONTH2_FORMAT_UPPER, constants.GO_MMM, -1)
	}   else if strings.Contains(goFormat, constants.MONTH2_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.MONTH2_FORMAT_LOWER, constants.GO_MMM, -1)
	}  else if strings.Contains(goFormat, constants.MONTH3_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.MONTH3_FORMAT_LOWER, constants.GO_MM, -1)
	} else if strings.Contains(goFormat, constants.MONTH3_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.MONTH3_FORMAT_UPPER, constants.GO_MM, -1)
	}

	if strings.Contains(goFormat, constants.DAY1_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.DAY1_FORMAT_LOWER, constants.GO_DDDD, -1)
	} else if strings.Contains(goFormat, constants.DAY2_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.DAY2_FORMAT_LOWER, constants.GO_DDD, -1)
	} else if strings.Contains(goFormat, constants.DAY3_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.DAY3_FORMAT_LOWER, constants.GO_DD, -1)
	} else if strings.Contains(goFormat, constants.DAY1_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.DAY1_FORMAT_UPPER, constants.GO_DDDD, -1)
	} else if strings.Contains(goFormat, constants.DAY2_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.DAY2_FORMAT_UPPER, constants.GO_DDD, -1)
	} else if strings.Contains(goFormat, constants.DAY3_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.DAY3_FORMAT_UPPER, constants.GO_DD, -1)
	}

	if strings.Contains(goFormat, constants.HOUR_FORMAT_UPPER) {
		goFormat = strings.Replace(goFormat, constants.HOUR_FORMAT_UPPER, constants.GO_HH, -1)
	} else if strings.Contains(goFormat, constants.HOUR_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.HOUR_FORMAT_LOWER, constants.GO_HH, -1)
	}


	if strings.Contains(goFormat, constants.MIN_FORMAT_UPPER){
		goFormat = strings.Replace(goFormat, constants.MIN_FORMAT_UPPER, constants.GO_MIN, -1)
	} else if strings.Contains(goFormat, constants.MIN_FORMAT_LOWER){
		goFormat = strings.Replace(goFormat, constants.MIN_FORMAT_LOWER, constants.GO_MIN, -1)
	}

	if strings.Contains(goFormat, constants.SEC_FORMAT_LOWER) {
		goFormat = strings.Replace(goFormat, constants.SEC_FORMAT_LOWER, constants.GO_SEC, -1)
	} else if strings.Contains(goFormat, constants.SEC_FORMAT_UPPER){
		goFormat = strings.Replace(goFormat, constants.SEC_FORMAT_UPPER, constants.GO_SEC, -1)
	}
	return goFormat
}