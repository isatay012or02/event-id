package server

import (
	"fmt"
	"time"

	"git.homebank.kz/homebank-oauth/halykid-events/controllers"
	ginmiddlewares "git.homebank.kz/libs/gin-middlewares/v2"
)

const DEFAULT_SERVER_SHUTDOWNTIMEOUT time.Duration = time.Duration(5) * time.Second

// GINSettings настройки GIN
type GINSettings struct {
	// UseRecovery нужно ли использовать родной middleware для перехвата паник
	UseRecovery bool
	// UseLogger нужно ли использовать родной middleware для access log-ов
	UseLogger bool
	// ReleaseMode необходимость запустить GIN в releae mode
	ReleaseMode bool
}

// Settings Настройки веб сервера
type Settings struct {
	// Port порт веб сервлера
	Port int
	// ShutdownTimeout таймаут остановки сервера для завершения запросов "в полете"
	ShutdownTimeout time.Duration
	// GIN настройки GIN
	GIN GINSettings
	// Controllers настройки обработчиков http запросов
	Controllers      controllers.Settings
	CorsSettings     ginmiddlewares.CORSSettings
	ProfilingEnabled bool
}

func (s *Settings) validate() error {
	if s.Port <= 0 {
		return fmt.Errorf("порт веб сервера должен быть положительным, но в конфигурации указано значение %d", s.Port)
	}

	if s.ShutdownTimeout <= 0 {
		s.ShutdownTimeout = DEFAULT_SERVER_SHUTDOWNTIMEOUT
	}

	return nil
}
