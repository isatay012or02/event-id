package models

// CodeMessage описывает базовый ответ code+message
type CodeMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
