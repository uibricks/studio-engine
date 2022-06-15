package utils

import (
	"fmt"
)

func GetParentIDFilter(parentId int32) string {
	if parentId > 0 {
		return fmt.Sprintf("parent_id=%d", parentId)
	}
	return "parent_id is null"
}
