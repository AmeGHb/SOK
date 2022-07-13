package config

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HttpPort  int
	KafkaPort int
	*DatabaseConfig
}

type DatabaseConfig struct {
	Url               string
	Attempts, Seconds int
	Delay             time.Duration
}

func New() *Config {

	/*
		confPath, getConfErr := getConfigDirPath()
		if getConfErr != nil {
			log.Fatalf("Configuration directory was not found. Error message: %v", getConfErr)
		}

		viper.AddConfigPath(confPath)
	*/
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

	DbUrl, ok := viper.Get("DATABASE_URL").(string)
	if !ok {
		log.Println("Database URL was not loaded from config.yaml")
	}

	DBAttempts, ok := viper.Get("DATABASE_ATTEMPTS").(int)
	if !ok {
		log.Println("Database attempts value was not loaded from config.yaml")
	}

	DBSeconds, ok := viper.Get("DATABASE_SECONDS").(int)
	if !ok {
		log.Println("Database seconds value was not loaded from config.yaml")
	}

	DBDelay, ok := viper.Get("DATABASE_DELAY").(int)
	if !ok {
		log.Println("Database delay value was not loaded from config.yaml")
	}

	DBDelayTime := time.Duration(DBDelay)

	return &Config{
		HttpPort:  HttpPort,
		KafkaPort: KafkaPort,
		DatabaseConfig: &DatabaseConfig{
			Url:      DbUrl,
			Attempts: DBAttempts,
			Seconds:  DBSeconds,
			Delay:    DBDelayTime,
		},
	}
}

func getConfigDirPath() (string, error) {

	dir_path, err := filepath.Abs("./")
	if err != nil {
		log.Fatalln(err)
	}

	stringIndex := strings.Index(dir_path, "transaction")

	if stringIndex == -1 {
		return "", errors.New("Config path - '" + dir_path + "' was not found.")
	}

	return dir_path[:stringIndex] + "transaction/", nil
}

func (c *Config) GetHttpPort() string {
	return fmt.Sprintf("%d", c.HttpPort)
}

func (c *Config) GetKafkaPort() string {
	return fmt.Sprintf("%d", c.KafkaPort)
}

func (c *Config) GetDatabaseURL() string {
	return c.DatabaseConfig.Url
}

func (c *Config) GetDatabaseAttempts() string {
	return fmt.Sprintf("%d", c.DatabaseConfig.Attempts)
}

func (c *Config) GetDatabaseSeconds() string {
	return fmt.Sprintf("%d", c.DatabaseConfig.Seconds)
}

func (c *Config) GetDatabaseDelay() string {
	return fmt.Sprintf("%d", c.DatabaseConfig.Delay)
}
