package controllers

import (
	"net/http"

	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
	"github.com/gin-gonic/gin"
)

func StatusOK(c *gin.Context, resp interface{}) {

	c.JSON(http.StatusOK, resp)
}

func StatusBadRequest(c *gin.Context, log service.AppLogger, method string, code int, message, customMessage string, data interface{}) {

	log.Error(method, "", "", code, message, customMessage, nil)

	c.JSON(http.StatusBadRequest, models.CodeMessage{
		Code:    code,
		Message: customMessage,
	})
}
