package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Schema struct {
	Host   string            `mapstructure:"host"`
	Server map[string]string `mapstructure:"server"`
	Mongo  map[string]struct {
		Host    string        `mapstructure:"host"`
		Name    string        `mapstructure:"name"`
		User    string        `mapstructure:"user"`
		Pass    string        `mapstructure:"pass"`
		Timeout time.Duration `mapstructure:"timeout"`
	} `mapstructure:"mongo"`
	Redis map[string]struct {
		Host string `mapstructure:"host"`
		Name string `mapstructure:"name"`
		DB   int    `mapstructure:"database"`
		User string `mapstructure:"user"`
		Pass string `mapstructure:"pass"`
	} `mapstructure:"redis"`
	RabbitMQ string `mapstructure:"rabbit_connection"`
	API      struct {
		Extends map[string]string `mapstructure:"extends"`
	} `mapstructure:"api"`
	Consumer struct {
		FromExchange    string            `mapstructure:"from_exchange"`
		PushExchange    string            `mapstructure:"push_exchange"`
		CodExchange     string            `mapstructure:"cod_exchange"`
		MonitorExchange string            `mapstructure:"monitor_exchange"`
		Queues          map[string]string `mapstructure:"queues"`
	} `mapstructure:"consumer"`
	ClientConfig struct {
		TimeOut       string `mapstructure:"time_out"`
		TimeOutTenant int    `mapstructure:"time_out_tenant"`
		TimeClock     string `mapstructure:"time_clock"`
	} `mapstructure:"client_config"`
	ReleaseMode bool   `mapstructure:"release_mode"`
	SecretKey   string `mapstructure:"secret_key"`
}

var DataConfig Schema

func init() {
	config := viper.New()
	env := os.Getenv("GO_ENV")
	config.AddConfigPath(".")
	config.AddConfigPath("config/env")
	if env == "prod" {
		config.SetConfigName("prod")
	} else if env == "test" {
		config.SetConfigName("test")
	} else {
		config.SetConfigName("dev")
	}
	config.AutomaticEnv()
	err := config.ReadInConfig() // Find and read the config file
	if err != nil {              // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = config.Unmarshal(&DataConfig)
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetConfig() Schema {
	return DataConfig
}
