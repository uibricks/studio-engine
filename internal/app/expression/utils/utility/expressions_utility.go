package main

import (
	"encoding/json"
	"fmt"
	"github.com/uibricks/studio-engine/internal/app/expression/utils"

	"github.com/goplusjs/gopherjs/js"
)

func main() {
	js.Module.Get("exports").Set("evalExpressions", evalExpressions)
	js.Module.Get("exports").Set("evalExpression", evalExpression)
}

//evalExpression is a helper function fon conversion of go file to js code
//it returns a map and error
func evalExpressions(expressionMetaDataMap map[string]interface{}, expressionsMap []map[string]interface{}, data string) (map[string]interface{}, string) {
	fmt.Println("evalExpressions starts")
	var expressionMetaData map[string]utils.ExpressionMD
	var expressions []utils.Expression
	expressionMetaDataStr, _ := json.Marshal(expressionMetaDataMap)
	expressionsStr, _ := json.Marshal(expressionsMap)
	if err := json.Unmarshal(expressionMetaDataStr, &expressionMetaData); err != nil {
		return nil, err.Error()
	}
	if err := json.Unmarshal(expressionsStr, &expressions); err != nil {
		return nil, err.Error()
	}
	res, err := utils.GetExpressionVal(expressionMetaData, expressions, data)
	//fmt.Println("res : ", res)
	fmt.Println("err : ", err)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return res, errMsg
}

func evalExpression(expressionMetaDataMap map[string]interface{}, expressionsMap map[string]interface{}, data string) (map[string]interface{}, string) {
	fmt.Println("evalExpression starts")
	var expressionsMetaData map[string]utils.ExpressionMD
	var expressionMetaData map[string]utils.ExpressionMD
	expressions := make([]utils.Expression, 0)
	expressionsMetaDataStr, _ := json.Marshal(expressionMetaDataMap)
	expressionMetaDataStr, _ := json.Marshal(expressionsMap)
	if err := json.Unmarshal(expressionsMetaDataStr, &expressionsMetaData); err != nil {
		return nil, err.Error()
	}
	if err := json.Unmarshal(expressionMetaDataStr, &expressionMetaData); err != nil {
		return nil, err.Error()
	}

	for k,v := range expressionMetaData {
		expressions = append(expressions, utils.Expression{ID: k})
		expressionsMetaData[k] = v
	}

	res, err := utils.GetExpressionVal(expressionsMetaData, expressions, data)
	//fmt.Println("res : ", res)
	fmt.Println("err : ", err)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return res, errMsg
}