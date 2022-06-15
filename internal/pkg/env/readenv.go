package env

import (
	"github.com/uibricks/studio-engine/internal/pkg/logger"
	"os"
	"strconv"
)

func Get(name string) string {
	v := os.Getenv(name)
	return v
}

func GetBool(name string) bool {
	v := os.Getenv(name)
	return v != "" && v == "true"
}

func MustGet(name string) string {
	v := os.Getenv(name)
	if v == "" {
		logger.Sugar.Panicf("Environment variable is missing - %s", name)
	}
	return v
}

func MustGetInt(name string) int {
	v := os.Getenv(name)
	if v == "" {
		logger.Sugar.Panicf("Environment variable is missing - %s", name)
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		logger.Sugar.Panicf("Environment variable, not a valid value for - %s", name)
	}
	return n
}
