package constants

import "time"

const(
	DefaultQueueExpiration time.Duration = 300
)

const (
	INVALID_PARAM_ERR = "invalid number of parameters for `%s`"
	INVALID_PARAM_TYPE_ERR = "expression contains one or more invalid type of parameters for `%s`"
	INVALID_FUNCTION_IN_PARAM = "expression contains invalid function `%s` for '%s'"
	INVALID_TYPE_ERR  = "invalid type of parameters for `%s`"
	NOT_FOUND = "element not found for '%s'"
	NO_EXPRESSIONS = "no expressions to evaluate"
)

const (
	DATE_TIME_FORMAT       = "01-02-2006 15:04:05"
	RESPONSE_DATA_TYPE_FLOAT = "number"
	RESPONSE_DATA_TYPE_INTEGER = "integer"
	RESPONSE_DATA_TYPE_BOOLEAN = "boolean"
	RESPONSE_DATA_TYPE_ARRAY = "arr"
)

const  (
	YEAR1_FORMAT_UPPER = "YYYY"
	YEAR1_FORMAT_LOWER = "yyyy"
	GO_YYYY = "2006"

	YEAR2_FORMAT_UPPER = "YY"
	YEAR2_FORMAT_LOWER = "yy"
	GO_YY = "06"

	MONTH1_FORMAT_UPPER = "MMMM"
	MONTH1_FORMAT_LOWER = "mmmm"
	GO_MMMM = "January"

	MONTH2_FORMAT_UPPER = "MMM"
	MONTH2_FORMAT_LOWER = "mmm"
	GO_MMM = "Jan"

	MONTH3_FORMAT_UPPER = "MM"
	MONTH3_FORMAT_LOWER = "mm"
	GO_MM = "01"

	DAY1_FORMAT_UPPER = "DDDD"
	DAY1_FORMAT_LOWER = "dddd"
	GO_DDDD = "Monday"

	DAY2_FORMAT_UPPER = "DDD"
	DAY2_FORMAT_LOWER = "ddd"
	GO_DDD = "Mon"

	DAY3_FORMAT_UPPER = "DD"
	DAY3_FORMAT_LOWER = "dd"
	GO_DD = "02"

	HOUR_FORMAT_UPPER = "HH"
	HOUR_FORMAT_LOWER = "hh"
	GO_HH = "15"

	MIN_FORMAT_UPPER = "MM"
	MIN_FORMAT_LOWER = "mm"
	GO_MIN = "04"

	SEC_FORMAT_UPPER = "SS"
	SEC_FORMAT_LOWER = "ss"
	GO_SEC = "15"

	MIN_FORMAT_IDENTIFIER = "min"
)

const (
	STRING_OP_LOWER = "lower"
	STRING_OP_UPPER = "upper"
	STRING_OP_TITLE = "title"
	STRING_OP_INDEX_ALL = "all"
	STRING_OP_INDEX_LAST = "last"
	STRING_OP_EQUAL_CS = "case-sensitive"
	STRING_OP_EQUAL_CIS = "case-insensitive"
	STRING_OP_CONTAINS = "contains"
	STRING_OP_CONTAINS_ANY = "containsany"
)

const (
	PATH_SEPARATOR = "."
	GROUP_EXPRESSION_TYPE = "group"
)

const (
	FUNCTION_CONCAT = "concat"
	FUNCTION_ENDSWITH = "endsWith"
	FUNCTION_STARTSWITH = "startsWith"
	FUNCTION_TRIM = "trim"
	FUNCTION_CONTAINS = "contains"
	FUNCTION_COUNT = "count"
	FUNCTION_EQUAL = "equal"
	FUNCTION_FIELDS = "fields"
	FUNCTION_INDEX = "index"
	FUNCTION_REPEAT = "repeat"
	FUNCTION_REPLACE = "replace"
	FUNCTION_SPLIT = "split"
	FUNCTION_STRING_TRANSFORM = "stringTransform"
	FUNCTION_REGEX = "regex"
	FUNCTION_CONCAT_DELIMITER = "concatDelimiter"
	FUNCTION_APPEND_TO_ARRAY = "appendToArray"
	FUNCTION_COUNT_ELEMENTS = "countElements"
	FUNCTION_INDEX_ARRAY = "indexArray"
	FUNCTION_INSERT = "insert"
	FUNCTION_POP = "pop"
	FUNCTION_REMOVE = "remove"
	FUNCTION_REVERSE = "reverse"
	FUNCTION_SORT_ARRAY = "sortArray"
	FUNCTION_EXTEND = "extend"
	FUNCTION_GET_AT = "getAt"
	FUNCTION_CURRENT_DATE = "currentDate"
	FUNCTION_CASE = "case"
	FUNCTION_ADD_DAYS = "addToDate"
	FUNCTION_IF = "if"
	FUNCTION_FORMAT_DATE = "formatDate"
)