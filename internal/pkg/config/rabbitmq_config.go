package config

type RabbitMqConfig struct {
	Host                  string `mapstructure:"host"`
	Port                  string `mapstructure:"port"`
	User                  string `mapstructure:"user"`
	Password              string `mapstructure:"pwd"`
	OrchestratorQueueName string `mapstructure:"orchestrator_queue_name"`
	MappingQueueName      string `mapstructure:"mapping_queue_name"`
	QueuePrefix           string `mapstructure:"queue_prefix"`
}
