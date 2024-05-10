package controllers

import (
	"fmt"
	"time"
)

type Settings struct {
	Cache              Redis
	Status             Status
	Kafka              Kafka
	SocketLifeDuration struct {
		Duration string
		TTL      time.Duration `json:"-"`
	}
}

type Status struct {
	Start   string
	Cancel  string
	Success string
}

type Kafka struct {
	Metrics Metrics
}

type Redis struct {
	Connect          ConnectType
	Metrics          Metrics
	KeyCheckAID      string
	KeyCheckAIDCount string
}

type ConnectType struct {
	Addr     string
	DB       int
	Password string
	PoolSize int
	Timeout  time.Duration
}

type Metrics struct {
	Enabled                 bool
	Label                   string
	DurationBuckets         []float64
	DurationBucketsToSocket []float64
}

// parse Парсит сырые данные в конфиге в готовые значения
func (c *Settings) Parse() error {
	var err error

	c.SocketLifeDuration.TTL, err = time.ParseDuration(c.SocketLifeDuration.Duration)
	if err != nil {
		return fmt.Errorf("parse user access token error: %v", err)
	}

	return nil
}
