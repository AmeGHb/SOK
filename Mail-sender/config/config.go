package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	HttpPort  int
	KafkaPort int
	MailPort  int
	MailHost  string
	MailFrom  string
}

func New() *Config {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}

	HttpPort, ok := viper.Get("HTTP_PORT").(int)
	if !ok {
		log.Println("Http port was not loaded from config.yaml")
	}

	KafkaPort, ok := viper.Get("KAFKA_PORT").(int)
	if !ok {
		log.Println("Http port was not loaded from config.yaml")
	}

	MailPort, ok := viper.Get("MAIL_PORT").(int)
	if !ok {
		log.Println("Mail port was not loaded from config.yaml")
	}

	MailHost, ok := viper.Get("MAIL_HOST").(string)
	if !ok {
		log.Println("Mail host was not loaded from config.yaml")
	}

	MailFrom, ok := viper.Get("MAIL_FROM").(string)
	if !ok {
		log.Println("Mail was not loaded from config.yaml")
	}

	return &Config{
		HttpPort:  HttpPort,
		KafkaPort: KafkaPort,
		MailPort:  MailPort,
		MailHost:  MailHost,
		MailFrom:  MailFrom,
	}
}

func (c *Config) GetHttpPort() string {
	return fmt.Sprintf("%d", c.HttpPort)
}

func (c *Config) GetKafkaPort() string {
	return fmt.Sprintf("%d", c.KafkaPort)
}

func (c *Config) GetMailPort() int {
	return c.MailPort
}

func (c *Config) GetMailHost() string {
	return c.MailHost
}

func (c *Config) GetMailFrom() string {
	return c.MailFrom
}
