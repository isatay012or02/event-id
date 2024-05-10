package controllers

import (
	"errors"
	"sync"
	"time"

	redissats "git.homebank.kz/libs/go-prometheus-go-redis.v8-stats"
	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheReturnedDataCnt        *prometheus.CounterVec
	cacheL                      sync.Mutex
	cacheDurationInSec          *prometheus.HistogramVec
	defaultCacheDurationBuckets []float64 = []float64{0.001, 0.003, 0.005, 0.007, 0.01, 0.015, 0.02, 0.025, 0.05, 0.075, 0.1, 0.15, 0.2, 0.3, 0.4, 0.5, 0.75, 1, 2, 3}
)

// Название проекта которое будет отображаться в метриках КОНФИГ
const ApplicationName string = "hb_halykid_events"

func addMetrics(cache *redis.Client, s *Settings) error {
	err := NewCacheMetrics(s)

	if err != nil {
		return err
	}

	if s.Cache.Metrics.Enabled {
		redisProm := redissats.NewStatsCollector(cache, s.Cache.Metrics.Label, "", "")
		prometheus.MustRegister(redisProm)
	}

	return nil
}

// Init Инициализирует метрики
func NewCacheMetrics(s *Settings) error {
	var err error

	if s.Cache.Metrics.Enabled {
		err = registerCacheReturnedCounter(s.Cache.Metrics.Label)
		if err != nil {
			return errors.New("cache returned counter")
		}
		err = registerCacheDurationHistogramm(s.Cache.Metrics.Label, s.Cache.Metrics.DurationBuckets)
		if err != nil {
			return errors.New("cache response durations")
		}
	}

	return nil
}

func registerCacheReturnedCounter(subSystem string) error {
	if cacheReturnedDataCnt != nil {
		return nil
	}

	cacheL.Lock()
	defer cacheL.Unlock()

	if cacheReturnedDataCnt != nil {
		return nil
	}

	if subSystem == "" {
		return errors.New("register cache returned response counter error: SubSytem not specified")
	}

	cacheReturnedDataCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subSystem,
			Name:      "get_count_request_from_cache",
			Help:      "счетчик количества запросов в кэш",
		},
		[]string{"isSuccess"},
	)
	return prometheus.Register(cacheReturnedDataCnt)
}

func registerCacheDurationHistogramm(subSystem string, buckets []float64) error {
	if cacheDurationInSec != nil {
		return nil
	}
	cacheL.Lock()
	defer cacheL.Unlock()
	if cacheDurationInSec != nil {
		return nil
	}

	if buckets == nil || len(buckets) == 0 {
		buckets = defaultCacheDurationBuckets
	}

	if subSystem == "" {
		return errors.New("register cache response durations histogram error: SubSytem not specified")
	}

	cacheDurationInSec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subSystem,
		Name:      "request_duration_in_seconds_cache",
		Help:      "Histogram для отображения длительности запросов в кэш",
		Buckets:   buckets,
	},
		[]string{"command", "status"})

	return prometheus.Register(cacheDurationInSec)
}

func CacheDurationObserve(d time.Duration, command string, err error) {
	if cacheDurationInSec == nil {
		return
	}
	status := "ok"
	if err != nil {
		status = "error"
	}
	cacheDurationInSec.WithLabelValues(command, status).Observe(d.Seconds())
}
