package controllers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"git.homebank.kz/homebank-oauth/halykid-events/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var pingResult = gin.H{
	"result": 1,
}

// Ping endpoint ping
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, pingResult)
}

type Controller interface {
	Connect(*gin.Context)
	RedirectURL(*gin.Context)
	UserAuthStatus(*gin.Context)
}

type impl struct {
	redis    *redis.Client
	settings Settings
	log      service.AppLogger
}

func (controller *impl) CacheGET(key string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), controller.settings.Cache.Connect.Timeout)
	start := time.Now()
	defer cancel()
	val, err := controller.redis.Get(ctx, key).Result()
	CacheDurationObserve(time.Since(start), "get", err)
	return val, err
}

func (controller *impl) GetClientIDByCache(aid string) (val string, err error) {
	keyCheckAID := fmt.Sprintf(controller.settings.Cache.KeyCheckAID, aid)
	val, err = controller.CacheGET(keyCheckAID)
	if err != nil {
		return val, err
	}

	return val, nil
}

func (controller *impl) CacheCheckByAID(key string) (bool, error) {

	ctx, cancel := context.WithTimeout(context.Background(), controller.settings.Cache.Connect.Timeout)

	start := time.Now()
	defer cancel()

	keyOfCheck := fmt.Sprintf(controller.settings.Cache.KeyCheckAID, key)
	val, err := controller.redis.Exists(ctx, keyOfCheck).Result()
	CacheDurationObserve(time.Since(start), "KeyCheckAIDExist", err)
	if val == 0 {
		return false, nil
	}

	counterKey := fmt.Sprintf(controller.settings.Cache.KeyCheckAIDCount, key)
	expirationTime := time.Hour
	val, err = controller.redis.Exists(ctx, counterKey).Result()
	metricDuration := time.Since(start)
	CacheDurationObserve(metricDuration, "KeyCheckAIDCountExist", err)
	if err != nil {
		return false, err
	}

	if val == 0 {
		counterVal := 1
		err = controller.redis.Set(ctx, counterKey, counterVal, expirationTime).Err()
		CacheDurationObserve(metricDuration, "KeyCheckAIDCountSet", err)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	return true, nil
}

func NewController(appLog service.AppLogger, settings Settings) (Controller, error) {
	r := redis.NewClient(&redis.Options{
		Addr:     settings.Cache.Connect.Addr,
		DB:       settings.Cache.Connect.DB,
		Password: settings.Cache.Connect.Password,
		PoolSize: settings.Cache.Connect.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), settings.Cache.Connect.Timeout)
	defer cancel()
	if err := r.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	controllers := &impl{
		settings: settings,
		redis:    r,
		log:      appLog,
	}

	addMetrics(r, &settings)

	return controllers, nil
}
