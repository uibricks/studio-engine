package config

import (
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
)

func LoadConfig(reader io.Reader, bindings map[string]string, appConfig interface{}) error {

	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := bind(bindings)
	if err != nil {
		return err
	}

	if err = viper.ReadConfig(reader); err != nil {
		return errors.Wrap(err, "Failed to load app config file")
	}

	if err = viper.Unmarshal(&appConfig); err != nil {
		return errors.Wrap(err, "Unable to parse app config file")
	}

	return nil
}

func bind(keysToEnvironmentVariables map[string]string) error {
	var bindErrors error
	for key, environmentVariable := range keysToEnvironmentVariables {
		if err := viper.BindEnv(key, environmentVariable); err != nil {
			bindErrors = multierror.Append(bindErrors, err)
		}
	}
	return bindErrors
}
