package controllers

import (
	"encoding/json"
	"fmt"
	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"github.com/gin-gonic/gin"
)

// RedirectURL обработчик endpoint-а получения redirect_url по AID
func (controller *impl) RedirectURL(c *gin.Context) {
	const method = "RedirectURL"
	aid := c.Query("aid")

	if aid == "" {
		StatusBadRequest(c, controller.log, method, -34258, "Не передан идентификар авторизации", "END", nil)
		return
	}

	redirectURL, err := controller.CacheGET(fmt.Sprintf("EXTENSION:ID:AID:%s", aid))
	if err != nil {
		StatusBadRequest(c, controller.log, method, -34259, err.Error(), "customMessage", gin.H{"aid": aid})
		return
	}
	var oauthStatus models.UserAuthStatus
	errToUnmarshal := json.Unmarshal([]byte(redirectURL), &oauthStatus)
	if errToUnmarshal != nil {
		return
	}
	if err != nil {
		StatusBadRequest(c, controller.log, method, -34315, err.Error(), "customMessage", gin.H{"aid": aid})
		return
	}

	StatusOK(c, models.RedirectURLSuccessResponse(oauthStatus.RedirectURL))
}

func (controller *impl) UserAuthStatus(c *gin.Context) {
	const method = "UserStatus"
	aid := c.Query("aid")

	redirectURL, err := controller.CacheGET(fmt.Sprintf("EXTENSION:ID:AID:%s", aid))
	if err != nil {
		if err.Error() == "redis: nil" {
			UserAuthStatus := models.UserAuthStatus{RedirectURL: "", Status: "InProgress"}
			StatusOK(c, models.UserAuthStatusSuccessResponse(UserAuthStatus))
			return
		}
		StatusBadRequest(c, controller.log, method, -34314, err.Error(), "customMessage", gin.H{"aid": aid})
		return
	}

	var oauthStatus models.UserAuthStatus
	err = json.Unmarshal([]byte(redirectURL), &oauthStatus)
	if err != nil {
		StatusBadRequest(c, controller.log, method, -34291, err.Error(), "customMessage", gin.H{"aid": aid})
		return
	}

	UserAuthStatus := models.UserAuthStatus{RedirectURL: oauthStatus.RedirectURL, Status: oauthStatus.Status}

	StatusOK(c, models.UserAuthStatusSuccessResponse(UserAuthStatus))
}
