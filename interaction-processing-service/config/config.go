package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of application.
// The values read by viper from a config file or environment variables.
type Config struct {
	Server struct {
		Port string `mapstructure:"Port"`
		Host string `mapstructure:"Host"`
		Env  string `mapstructure:"Env"`
	}
	LogLevel string `mapstructure:"LogLevel"`
	DB       struct {
		Host     []string `mapstructure:"Host"`
		Port     int      `mapstructure:"Port"`
		Username string   `mapstructure:"Username"`
		Password string   `mapstructure:"Password"`
		Name     string   `mapstructure:"Name"`
	} `mapstructure:"Database"`
	Redis struct {
		Port     int    `mapstructure:"Port"`
		Host     string `mapstructure:"Host"`
		Password string `mapstructure:"Password"`
	} `mapstructure:"Redis"`
	Kafka struct {
		Brokers  []string `mapstructure:"Brokers"`
		ClientID string   `mapstructure:"ClientID"`
	} `mapstructure:"Kafka"`
}

// LoadConfig reads configuration from file or enviroment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
