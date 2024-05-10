package events

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	broadcastClientStatusCnt      *prometheus.CounterVec
	kafkaDurationInSec            *prometheus.HistogramVec
	sessionDurationInSec          *prometheus.HistogramVec
	broadcastL                    sync.Mutex
	defaultKafkaDurationBuckets   = []float64{0.001, 0.003, 0.005, 0.007, 0.01, 0.015, 0.02, 0.025, 0.05, 0.075, 0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.75, 1, 2, 3}
	defaultSessionDurationBuckets = []float64{1, 3, 5, 7, 10, 15, 20, 25, 50, 75, 100, 150, 200, 300, 400}
)

func (sender *impl) GetOrCreateActiveSessionsCnt() error {

	err := prometheus.Register(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "active_sessions",
			Help: "Количество активных сессии",
		},
		func() float64 { return float64(sender.m.Len()) },
	))

	if err != nil {
		return err
	}

	return err
}

func RegisterKafkaDurationHistogram(subSystem string, buckets []float64) error {
	if kafkaDurationInSec != nil {
		return nil
	}
	broadcastL.Lock()
	defer broadcastL.Unlock()
	if kafkaDurationInSec != nil {
		return nil
	}

	if buckets == nil || len(buckets) == 0 {
		buckets = defaultKafkaDurationBuckets
	}

	if subSystem == "" {
		return errors.New("register cache response durations histogram error: SubSystem not specified")
	}

	kafkaDurationInSec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subSystem,
		Name:      "request_duration_in_seconds_kafka",
		Help:      "Histogram для отображения длительности запросов в сокет сессиях",
		Buckets:   buckets,
	},
		[]string{"code", "clientID", "status"})

	return prometheus.Register(kafkaDurationInSec)
}

func RegisterSessionDurationHistogram(subSystem string, buckets []float64) error {
	if sessionDurationInSec != nil {
		return nil
	}
	broadcastL.Lock()
	defer broadcastL.Unlock()
	if sessionDurationInSec != nil {
		return nil
	}

	if buckets == nil || len(buckets) == 0 {
		buckets = defaultSessionDurationBuckets
	}

	if subSystem == "" {
		return errors.New("register cache response durations histogram error: SubSystem not specified")
	}

	sessionDurationInSec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subSystem,
		Name:      "request_duration_in_seconds_session",
		Help:      "Histogram для отображения длительности сокет сессии",
		Buckets:   buckets,
	}, []string{"clientID"})

	return prometheus.Register(sessionDurationInSec)
}

func BroadcastStatusInc(сode int, clientID, status string) {
	if broadcastClientStatusCnt == nil {
		return
	}
	broadcastClientStatusCnt.WithLabelValues(strconv.Itoa(сode), clientID, status).Inc()
}

func BroadcastStatus(сode int, clientID, status string, duration time.Duration) {
	if kafkaDurationInSec == nil {
		return
	}
	kafkaDurationInSec.WithLabelValues(strconv.Itoa(сode), clientID, status).Observe(duration.Seconds())
}

func BroadcastSession(clientID string, duration time.Duration) {
	if sessionDurationInSec == nil {
		return
	}
	sessionDurationInSec.WithLabelValues(clientID).Observe(duration.Seconds())
}

func RegisterMetrics() error {
	if broadcastClientStatusCnt != nil {
		return nil
	}

	broadcastL.Lock()
	defer broadcastL.Unlock()

	broadcastClientStatusCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "halykidEvents",
			Name:      "ClientsStatus",
			Help:      "счетчик количества стартов в сокет сессиях",
		},
		[]string{"code", "clientID", "status"},
	)

	return prometheus.Register(broadcastClientStatusCnt)
}
