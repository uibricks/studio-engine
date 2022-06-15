package utils

import (
	"github.com/Knetic/govaluate"
	"github.com/uibricks/studio-engine/internal/app/expression/constants"
)

type Path struct {
	Name      string
}

// To store all custom functions
var functionMap = make(map[string]govaluate.ExpressionFunction)

type BaseExpressions struct {
	Expressions   []Expression            `json:"expressions"`
	ExpressionsMD map[string]ExpressionMD `json:"expressionMetadata"`
}

type Expression struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Children []Expression `json:"children"`
	isChild bool
	parentGroup string
}

type ExpressionMD struct {
	Id string `json:"id"`
	Name       string                   `json:"name"`
	Type       string                   `json:"type"`
	Refs       map[string]ExpressionRef `json:"refs"`
	NestedRefs []string                 `json:"nestedRefs"`
	Raw        string                   `json:"raw"`
}

type ExpressionRef struct {
	Type string   `json:"type"`
	Path []string `json:"path"`
}

func init() {

	funcNames := map[string]func(...interface{}) (interface{}, error) {
		constants.FUNCTION_CONCAT : concat,
		constants.FUNCTION_ENDSWITH : endsWith,
		constants.FUNCTION_STARTSWITH : startsWith,
		constants.FUNCTION_TRIM : trim,
		constants.FUNCTION_CONTAINS : contains,
		constants.FUNCTION_COUNT : count,
		constants.FUNCTION_EQUAL : equal,
		constants.FUNCTION_FIELDS : fields,
		constants.FUNCTION_INDEX : index,
		constants.FUNCTION_REPEAT : repeat,
		constants.FUNCTION_REPLACE: replace,
		constants.FUNCTION_SPLIT : split,
		constants.FUNCTION_STRING_TRANSFORM: stringTransform,
		constants.FUNCTION_REGEX : regex,
		constants.FUNCTION_CONCAT_DELIMITER : concatDelimiter,
		constants.FUNCTION_APPEND_TO_ARRAY : appendToArray,
		constants.FUNCTION_COUNT_ELEMENTS : countArrayElements,
		constants.FUNCTION_INDEX_ARRAY : indexArray,
		constants.FUNCTION_INSERT : insert,
		constants.FUNCTION_POP : pop,
		constants.FUNCTION_REMOVE : remove,
		constants.FUNCTION_REVERSE : reverse,
		constants.FUNCTION_SORT_ARRAY : sortArray,
		constants.FUNCTION_EXTEND : extend,
		constants.FUNCTION_GET_AT : getAt,
		constants.FUNCTION_CURRENT_DATE : currentDate,
		constants.FUNCTION_CASE : caseFunc,
		constants.FUNCTION_ADD_DAYS : addToDate,
		constants.FUNCTION_IF: ifFunc,
		constants.FUNCTION_FORMAT_DATE: formatDate,
	}

	for k,v := range funcNames {
		initFunctionMap(k, v)
	}
}


func initFunctionMap(funcName string, fn func(args ...interface{}) (interface{}, error)) {
	funcNames := permutations(funcName)
	for i := range funcNames {
		functionMap[funcNames[i]] = fn
	}
}
