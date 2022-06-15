package config

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"pwd"`
	Db       string `mapstructure:"db"`
}
