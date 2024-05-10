package controllers

import (
	"encoding/json"
	"fmt"
	"git.homebank.kz/homebank-oauth/halykid-events/events"
	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"github.com/gin-gonic/gin"
	"time"
)

// Connect обработчик endpoint-а подключения по web socket
func (controller *impl) Connect(c *gin.Context) {

	const method = "Connect"
	aid := c.Query("aid")
	val, err := controller.GetClientIDByCache(aid)
	if err != nil {
		msg := "Ошибка при получении clientID из кэша"
		StatusBadRequest(c, controller.log, method, -34317, msg, msg, gin.H{"aid": aid})
		events.BroadcastStatusInc(-34317, "Не определен", "Error")
		return
	}

	var data *models.QRInfo
	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		msg := "Ошибка преоброзования из байта в структуру"
		StatusBadRequest(c, controller.log, method, -34318, msg, msg, gin.H{"aid": aid})
		events.BroadcastStatusInc(-34318, "Не определен", "Error")
		return
	}

	if len(aid) < 1 {
		msg := "Не передан идентификар авторизации"
		StatusBadRequest(c, controller.log, method, -34260, msg, msg, gin.H{"aid": aid})
		events.BroadcastStatusInc(-34260, data.ClientID, "Error")
		return
	}

	if exist, err := controller.CacheCheckByAID(aid); !exist {
		msg := "Невалидное значение параметра AID"
		errMsg := msg
		if err != nil {
			errMsg = fmt.Sprintf("Описание ошибки: %s", err.Error())
		}

		StatusBadRequest(c, controller.log, method, -34272, errMsg, msg, gin.H{"aid": aid})
		events.BroadcastStatusInc(-34272, data.ClientID, "Error")
		return
	}

	start := time.Now()

	events.Default.HandleRequest(c.Writer, c.Request, aid, controller.settings.SocketLifeDuration.TTL)

	events.Default.HandleDisconnect(aid, data.ClientID, start)
}
