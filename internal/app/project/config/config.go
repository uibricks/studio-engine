package config

import (
	"flag"
	"github.com/uibricks/studio-engine/internal/pkg/config"
	"io"
	"os"
	"sync"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = ""
	configFileUsage   = "this is config file path"
)

var (
	once         sync.Once
	cachedConfig AppConfig
)

type AppConfig struct {
	ServerConfig     config.ServerConfig            `mapstructure:"app"`
	DatabaseConfig   config.DBConfig         `mapstructure:"db"`
	RedisConfig      config.RedisConfig      `mapstructure:"db"`
	RabbitMqConfig   config.RabbitMqConfig   `mapstructure:"db"`
	PrometheusConfig config.PrometheusConfig `mapstructure:"db"`
}

func LoadConfig(reader io.Reader) (c AppConfig, err error) {

	// Todo: Make certain config configurable via env vars
	keysToEnvironmentVariables := map[string]string{}

	err = config.LoadConfig(reader, keysToEnvironmentVariables, &c)

	if err != nil {
		return c, err
	}

	return c, nil
}

func ProvideAppConfig() (c AppConfig, err error) {
	once.Do(func() {
		var configFile string
		flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
		flag.Parse()

		var configReader io.ReadCloser
		configReader, err = os.Open(configFile)
		defer configReader.Close() //nolint

		if err != nil {
			return
		}

		c, err = LoadConfig(configReader)
		if err != nil {
			return
		}

		cachedConfig = c
	})

	return cachedConfig, err
}
