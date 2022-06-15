package config

type DBConfig struct {
	Url    string `mapstructure:"url"`
	Schema string `mapstructure:"schema"`
}