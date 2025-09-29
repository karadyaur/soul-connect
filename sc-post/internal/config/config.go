package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	EnvType       string `mapstructure:"ENV_TYPE"`
	ServerPort    string `mapstructure:"SERVER_PORT"`
	WebappBaseUrl string `mapstructure:"WEBAPP_BASE_URL"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	KafkaBrokers  string `mapstructure:"KAFKA_BROKERS"`
	KafkaTopic    string `mapstructure:"POST_EVENTS_TOPIC"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + ".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
	}

	return
}
