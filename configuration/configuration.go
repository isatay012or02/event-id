package configuration

import (
	"fmt"

	"git.homebank.kz/homebank-oauth/halykid-events/queue"
	"git.homebank.kz/homebank-oauth/halykid-events/server"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
	"github.com/spf13/viper"
)

var Config Configuration

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	viper.Unmarshal(&Config)
}

// Config конфигурация приложения
type Configuration struct {
	WEBServer server.Settings
	Queue     queue.Settings
	Logger    service.LoggerSettings
}
