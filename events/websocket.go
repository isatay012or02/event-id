package events

import (
	"encoding/json"
	"fmt"
	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var Default Sender
var locker = &sync.Mutex{}

type Sender interface {
	SendAuthCode(models.AuthEvent) error
	HandleRequest(w http.ResponseWriter, r *http.Request, aid string, duration time.Duration)
	HandleDisconnect(aid, clientID string, start time.Time)
}

type impl struct {
	m   *melody.Melody
	log service.AppLogger
}

func (sender *impl) SendAuthCode(authEvent models.AuthEvent) error {

	const method = "SendAuthCode"

	countSession := sender.m.Len()
	sender.log.Info(method, "", "", fmt.Sprintf("Get message before Broadcast, the count of active session is: %s",
		strconv.Itoa(countSession)), "Before broadcast", authEvent)

	if countSession == 0 {
		sender.log.Warn(method, "", "", -34269, "Get message before Broadcast, the count of active session is: 0", "Before broadcast", authEvent)
	}

	authEventBytes, _ := json.Marshal(authEvent)

	sender.m.BroadcastFilter(authEventBytes, func(s *melody.Session) bool {

		BroadcastStatusInc(0, authEvent.ClientID, authEvent.Status)

		timeStart := time.Now()
		sessionAID, ok := s.Get("aid")
		sender.log.Info(method, "", "", "Get message in Broadcast", "Start", gin.H{"authEvent": authEvent, "sessionAID": sessionAID})

		if !ok {
			sender.log.Warn(method, "", "", -34269, "Не получен параметр aid в при открытии socket сессии", "END", authEvent)
			BroadcastStatus(-34269, authEvent.ClientID, authEvent.Status, time.Since(timeStart))
			return false
		}

		aid, isStr := sessionAID.(string)
		if !isStr {
			sender.log.Error(method, "", "", -34270, "Ошибка преобразование параметра aid в тип string", "END", authEvent)
			BroadcastStatus(-34270, authEvent.ClientID, authEvent.Status, time.Since(timeStart))
			return false
		}
		if aid == authEvent.AID {
			sender.log.Info(method, "", "", "Совпадение параметра aid из socket сессии со значением из aid из сообщения с kafka", "END", authEvent)
			BroadcastStatus(0, authEvent.ClientID, authEvent.Status, time.Since(timeStart))
			return true
		}

		sender.log.Warn(method, "", "", -34271, "Не совпадение параметра из socket сессии с значением из aid из сообщения с kafka", "END", gin.H{"authEvent": authEvent, "sessionAID": sessionAID})
		BroadcastStatus(-34271, authEvent.ClientID, authEvent.Status, time.Since(timeStart))
		return false

	})
	return nil
}

func (sender *impl) HandleRequest(w http.ResponseWriter, r *http.Request, aid string, duration time.Duration) {
	const method = "HandleRequest"

	sender.log.Info(method, "", aid, fmt.Sprintf("starting new session for aid: %s", aid), "used session id as userID", nil)

	sender.m.HandleConnect(func(s *melody.Session) {
		go startSessionTimer(s, duration)
	})

	sender.m.HandleRequestWithKeys(w, r, map[string]interface{}{
		"aid": aid,
	})
}

func (sender *impl) HandleDisconnect(aid, clientID string, start time.Time) {
	const method = "HandleRequest"

	sender.m.HandleDisconnect(func(s *melody.Session) {
		sender.log.Info(method, "", "", fmt.Sprintf("Произашло отключение сессии по AID: %s", aid), "INFO", nil)
		resetSessionTimer(s)
	})

	duration := time.Since(start)
	BroadcastSession(clientID, duration)
}

func startSessionTimer(s *melody.Session, duration time.Duration) {
	timer := time.NewTimer(duration)

	select {
	case <-timer.C:
		s.Close()
	}
}

func resetSessionTimer(s *melody.Session) {
	s.Close()
}

func NewSender(appLog service.AppLogger, isDefault ...bool) (Sender, error) {
	sender := &impl{
		m:   melody.New(),
		log: appLog,
	}

	d := false
	if len(isDefault) > 0 {
		d = isDefault[0]
	}

	if Default == nil || d {
		locker.Lock()
		defer locker.Unlock()
		if Default == nil || d {
			Default = sender
		}
	}

	sender.GetOrCreateActiveSessionsCnt()

	return sender, nil
}

func Init(log service.AppLogger) error {
	_, err := NewSender(log, true)

	return err
}
