package utils

import (
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/tidwall/gjson"
	"reflect"
	"strings"

	"github.com/uibricks/studio-engine/internal/app/expression/constants"
)

// GetExpressionVal is used to get all the expressions and their corresponding value
func GetExpressionVal(expressionMetaData map[string]ExpressionMD, expressions []Expression, data string) (resp map[string]interface{}, respErr error) {

	// to return error message in case of any error
	defer func() {
		if panicErr := recover(); panicErr != nil {
			fmt.Println("panic")
			fmt.Println(panicErr)
			resp = nil
			respErr = errors.New(fmt.Sprintf("error occured while resolving expression : %v", panicErr))
		}
	}()

	var mappedVal = make(map[string]interface{})
	responseData := make(map[string]interface{})

	// creating a copy of expressions by expanding groups
	// the original expressions will be used to reply back only groups and not children
	expressionsCopy := make([]Expression, 0)

	expandExpressions(expressions, &expressionsCopy, false, "")

	expressionMetaData = FilterExpressions(expressionMetaData, expressionsCopy)

	if len(expressionMetaData) == 0 {
		return nil, fmt.Errorf(constants.NO_EXPRESSIONS)
	}

	for k, expressionMD := range expressionMetaData {
		if err := AddExpressions(expressionMetaData, expressionMD, data, mappedVal, k, expressionsCopy); err != nil {
			return nil, err
		}
	}

	evaluateGroups(expressionMetaData, expressions, mappedVal, "")

	// To filter evaluated 'refs', 'nested refs' and children-expressions;only retain expressions requested via expressions
	for i := range expressions {
		expression := expressions[i]
		id := expression.ID
		if _, ok := mappedVal[id]; ok {
			if reflect.ValueOf(mappedVal[id]).Kind() == reflect.Slice {
				if len(mappedVal[id].([]interface{})) == 1 {
					mappedVal[id] = mappedVal[id].([]interface{})[0]
				}
			}
			responseData[id] = flatten(mappedVal[id])
		}
	}

	return responseData, nil
}

// FilterExpressions retains only those expressions in expression-menu for evaluation
func FilterExpressions(expressionMetaData map[string]ExpressionMD, expressions []Expression) map[string]ExpressionMD {
	toEvaluateExpressions := make(map[string]ExpressionMD)

	for i := range expressions {
		expression := expressions[i]
		id := expression.ID
		if exp, ok := expressionMetaData[id]; ok {
			toEvaluateExpressions[id] = exp
			includeNestedRefs(exp, expressionMetaData, &toEvaluateExpressions)
		}
	}
	return toEvaluateExpressions
}

// includeNestedRefs adds expression's nested refs in expression-menu for evaluation
func includeNestedRefs(expMd ExpressionMD, expressionMetaData map[string]ExpressionMD, toEvaluateExpressions *map[string]ExpressionMD) {
	if len(expMd.NestedRefs) > 0 {
		for j := range expMd.NestedRefs {
			id := expMd.NestedRefs[j]
			if exp, ok := expressionMetaData[id]; ok {
				(*toEvaluateExpressions)[id] = exp
			}
		}
	}
}

// AddExpressions is used to get the expression value from the response
func AddExpressions(expressionMetaData map[string]ExpressionMD, expressionMD ExpressionMD, resp string, mappedVal map[string]interface{}, k string, expressions []Expression) error {
	getValueFromResponse(expressionMD, resp, mappedVal)

	if err := addNestedRef(expressionMetaData, expressionMD, resp, mappedVal, expressions); err != nil {
		return err
	}

	if _, ok := mappedVal[k]; !ok {
		raw := expressionMD.Raw
		raw = addExtraArgToConcat(raw)
		if expValuate, err := govaluate.NewEvaluableExpressionWithFunctions(raw, functionMap); err != nil {
			return err
		} else if result, err := expValuate.Evaluate(mappedVal); err != nil {
			return err
		} else if typeOfVal := expressionMD.Type; typeOfVal != "" {

			// only child expressions of a group are wrapped within their expression id
			if isCh, _ := isExpressionChild(k,expressions); isCh {
				result = wrapDataWithKey(result, k)
				mappedVal[k] = result
				return nil
			}
			mappedVal[k] = result
		}
	}
	return nil
}

// getValueFromResponse get the parameter val from response based on parameter id and parameter path
func getValueFromResponse(expressionMD ExpressionMD, resp string, mappedVal map[string]interface{}) {
	refsMap := expressionMD.Refs

	for k, expressionRef := range refsMap {
		dataPath := strings.Join(expressionRef.Path, constants.PATH_SEPARATOR)
		respResult := gjson.Get(resp, dataPath)
		typeOfVal := expressionRef.Type
		if respResult.Exists() {
			convertTypeOfResponseData(typeOfVal, k, respResult, mappedVal)
		}
	}
}

// isExpressionChild determines if an expression belongs to a group
func isExpressionChild(expressionID string, expressions []Expression) (bool, string) {
	for i := range expressions {
		if expressions[i].ID == expressionID && expressions[i].isChild {
			return true, expressions[i].parentGroup
		}
	}
	return false, ""
}

// convertTypeOfResponseData converts the data to  provided data type
func convertTypeOfResponseData(typeOfVal, key string, respResult gjson.Result, mappedVal map[string]interface{}) {
	dataType := strings.ToLower(typeOfVal)
	if dataType == reflect.String.String() {
		mappedVal[key] =  respResult.String()
	} else if dataType == constants.RESPONSE_DATA_TYPE_FLOAT {
		mappedVal[key] =  respResult.Float()
	} else if dataType == constants.RESPONSE_DATA_TYPE_INTEGER {
		mappedVal[key] =  respResult.Int()
	} else if dataType == constants.RESPONSE_DATA_TYPE_BOOLEAN {
		mappedVal[key] = respResult.Bool()
	} else if strings.Contains(dataType, constants.RESPONSE_DATA_TYPE_ARRAY) {
		resultArray, _ := respResult.Value().([]interface{})
		mappedVal[key] =  resultArray
	} else {
		mappedVal[key] = fmt.Sprintf("%v", respResult)
	}
}

// addNestedRef adds all the nested value inside map so that it can be used during run time
func addNestedRef(expressionMetaData map[string]ExpressionMD, expressionMD ExpressionMD, resp string, mappedVal map[string]interface{}, expressions []Expression) error {
	nestedRefArr := expressionMD.NestedRefs
	for i := range nestedRefArr {
		key := fmt.Sprintf("%s", nestedRefArr[i])
		if _, ok := mappedVal[key]; !ok {
			if val, ok := expressionMetaData[key]; ok {
				return AddExpressions(expressionMetaData, val, resp, mappedVal, key, expressions)
			}
		}
	}
	return nil
}

// expandExpressions creates a copy of expression and their children in same expression array for evaluation
func expandExpressions(expressions []Expression, resExpressions *[]Expression, isChild bool, parent string) {

	for _,val := range expressions {
		if isChild {
			val.isChild = true
			val.parentGroup = parent
		}
		*resExpressions = append(*resExpressions, val)
		if len(val.Children) > 0 && strings.Contains(val.Type, constants.GROUP_EXPRESSION_TYPE) {
			expandExpressions(val.Children, resExpressions, true, val.ID)
		}
	}
}

// evaluateGroups populates the group expression as a map of its children expressions by combining them
func evaluateGroups(expressionMD map[string]ExpressionMD,expressions []Expression, mappedVal map[string]interface{}, currentGroupId string) {
	expChildMap := make(map[string][]string)
	for _,val := range expressions {
		if strings.Contains(strings.ToLower(val.Type), constants.GROUP_EXPRESSION_TYPE) && len(val.Children) > 0 {

			for _,ch := range val.Children {
				if exps, ok := expChildMap[val.ID]; ok {
					exps = append(exps, ch.ID)
					expChildMap[val.ID] = exps
				} else {
					expChildMap[val.ID] = []string{ch.ID}
				}
			}

			// making 1st child of a group as a non-group-type child
			if len(val.Children) > 0 && val.Children[0].Type == constants.GROUP_EXPRESSION_TYPE {
				for i := 1; i < len(val.Children); i++ {
					if val.Children[i].Type != constants.GROUP_EXPRESSION_TYPE {
						val.Children[0], val.Children[i] = val.Children[i], val.Children[0]
						break
					}
				}
			}
			// groups will get evaluated while backtracking
			evaluateGroups(expressionMD, val.Children, mappedVal, val.ID)
		}
		if len(currentGroupId) > 0 {
			if _, ok := mappedVal[currentGroupId]; ok {
				if val.Type != constants.GROUP_EXPRESSION_TYPE {
					mappedVal[currentGroupId] = mergeObjs(mappedVal[currentGroupId], mappedVal[val.ID])
				} else {
					if reflect.ValueOf(mappedVal[currentGroupId]).Kind() == reflect.Slice {
						lvl1 := getNestedDepth(mappedVal[val.ID])
						lvl2 := getNestedDepth(mappedVal[currentGroupId])
						res := mergeObjs(mappedVal[currentGroupId], mappedVal[val.ID])
						if lvl1 > lvl2 {
							res = nestedGroup(res,mappedVal[val.ID], val.ID, expChildMap)
						} else {
							res = nestedGroup(mappedVal[currentGroupId], res, val.ID, expChildMap)
						}
						mappedVal[currentGroupId] = res
					} else if reflect.ValueOf(mappedVal[currentGroupId]).Kind() == reflect.Map {
						tempData := make(map[string]interface{})
						d := reflect.ValueOf(mappedVal[currentGroupId])
						for _,k := range d.MapKeys() {
							tempData[k.String()] = d.MapIndex(k).Interface()
						}
						tempData[val.ID] = mappedVal[val.ID]
						mappedVal[currentGroupId] = tempData
					} else {
						mappedVal[currentGroupId] = []interface{}{mappedVal[currentGroupId], map[string]interface{}{val.ID: mappedVal[val.ID]}}
					}

				}
			} else {
				if val.Type == constants.GROUP_EXPRESSION_TYPE {
					mappedVal[currentGroupId] = map[string]interface{}{val.ID: mappedVal[val.ID]}
				} else {
					mappedVal[currentGroupId] = mappedVal[val.ID]
				}
			}
		}
	}
}

// nestedGroup wraps expression-results within a single object by their group-id
// obj1 represents before merge state
// obj2 represents after merge state
func nestedGroup(obj1 interface{}, obj2 interface{}, groupId string, expChildMap map[string][]string) interface{} {
	if reflect.ValueOf(obj1).Kind() == reflect.Slice && reflect.ValueOf(obj2).Kind() == reflect.Slice {
		d1 := reflect.ValueOf(obj1)
		d2 := reflect.ValueOf(obj2)
		tmpData := make([]interface{}, d1.Len())
		returnSlice := make([]interface{}, d1.Len())
		for i := 0; i < d1.Len(); i++ {
			tmpData[i] = d1.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i] = nestedGroup(v, d2.Index(i).Interface(), groupId, expChildMap)
		}
		return returnSlice
	} else if reflect.ValueOf(obj1).Kind() == reflect.Map && reflect.ValueOf(obj2).Kind() == reflect.Map {
		tempData := make(map[string]interface{})
		d1 := reflect.ValueOf(obj1)
		d2 := reflect.ValueOf(obj2)
		for _,k := range d2.MapKeys() {
			_, present := FindInArray(expChildMap[groupId], k.String())
			if present {
				tempData[k.String()] = d2.MapIndex(k).Interface()
			}
		}
		tempData = map[string]interface{}{groupId: tempData}
		for _,k := range d1.MapKeys() {
			_, present := FindInArray(expChildMap[groupId], k.String())
			if !present {
				tempData[k.String()] = d1.MapIndex(k).Interface()
			}

		}
		return tempData
	}
	return nil
}