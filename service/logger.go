package service

import (
	kgolw "git.homebank.kz/libs/logger-go/kafkago"
	logger "git.homebank.kz/libs/logger-go/logger"
	loggergo "git.homebank.kz/libs/logger-go/logger"
)

type appLogger struct {
	log logger.Logger
}

type AppLogger interface {
	Error(action, clientID, userID string, code int, message, customMessage string, data interface{})
	Warn(action, clientID, userID string, code int, message, customMessage string, data interface{})
	Info(action, clientID, userID string, message, customMessage string, data interface{})
	ServerInfo(action, message string)
	ServerError(action, message, customMessage string)
}

func NewAppLogger(cnf *LoggerSettings) (AppLogger, error) {

	lg, err := logger.New(logger.Settings{
		Logger: logger.LoggerSettings{
			Component: cnf.Component,
			OnFinish:  logger.PrintErrorInStdout(),
			MinLevel:  cnf.MinLevel,
		},
		Writer: &kgolw.Settings{
			Brokers: cnf.Writer.Brokers,
			Topic:   cnf.Writer.Topic,
		},
	})
	if err != nil {
		return nil, err
	}

	return &appLogger{
		log: lg,
	}, nil
}

func (r *appLogger) Warn(action, clientID, userID string, code int, message, customMessage string, data interface{}) {
	go r.log.Warn(loggergo.Message{
		Action:        action,
		Code:          &code,
		Message:       &message,
		CustomMessage: &customMessage,
		Data:          data,
		ClientID:      &clientID,
		UserID:        &userID,
	})
}

func (r *appLogger) Error(action, clientID, userID string, code int, message, customMessage string, data interface{}) {
	go r.log.Warn(loggergo.Message{
		Action:        action,
		Code:          &code,
		Message:       &message,
		CustomMessage: &customMessage,
		Data:          data,
		ClientID:      &clientID,
		UserID:        &userID,
	})
}

func (r *appLogger) Info(action, clientID, userID string, message, customMessage string, data interface{}) {
	go r.log.Info(loggergo.Message{
		Action:        action,
		Message:       &message,
		CustomMessage: &customMessage,
		Data:          data,
		ClientID:      &clientID,
		UserID:        &userID,
	})
}

func (r *appLogger) ServerInfo(action, message string) {
	go r.log.Info(loggergo.Message{
		Action:        action,
		Message:       &message,
		CustomMessage: nil,
		Data:          nil,
		ClientID:      nil,
		UserID:        nil,
	})
}

func (r *appLogger) ServerError(action, message, customMessage string) {
	go r.log.Info(loggergo.Message{
		Action:        action,
		Message:       &message,
		CustomMessage: &customMessage,
		Data:          nil,
		ClientID:      nil,
		UserID:        nil,
	})
}
