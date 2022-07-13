package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HttpPort  int
	KafkaPort int
	SMTPPort  int
}

func New() *Config {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	HttpPort, ok := viper.Get("HTTP_PORT").(int)
	if !ok {
		log.Println("Http port was not loaded from config.yaml")
	}

	KafkaPort, ok := viper.Get("KAFKA_PORT").(int)
	if !ok {
		log.Println("Http port was not loaded from config.yaml")
	}

	SMTPPort, ok := viper.Get("SMTP_PORT").(int)
	if !ok {
		log.Println("SMTP port was not loaded from config.yaml")
	}

	return &Config{
		HttpPort:  HttpPort,
		KafkaPort: KafkaPort,
		SMTPPort:  SMTPPort,
	}
}

func (c *Config) GetHttpPort() string {
	return fmt.Sprintf("%d", c.HttpPort)
}

func (c *Config) GetKafkaPort() string {
	return fmt.Sprintf("%d", c.KafkaPort)
}

func (c *Config) GetSMTPPort() string {
	return fmt.Sprintf("%d", c.SMTPPort)
}
