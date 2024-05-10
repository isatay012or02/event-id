package queue

import (
	"context"
	"encoding/json"
	"git.homebank.kz/homebank-oauth/halykid-events/events"
	"git.homebank.kz/homebank-oauth/halykid-events/models"
	"git.homebank.kz/homebank-oauth/halykid-events/service"
	collector "git.homebank.kz/libs/go-prometheus-kafka"
	"github.com/prometheus/client_golang/prometheus"
	kafka "github.com/segmentio/kafka-go"
	"io"
	"sync"
)

type worker struct {
	reader  *kafka.Reader
	wg      sync.WaitGroup
	running bool
	log     service.AppLogger
}

func newWoker(l service.AppLogger, s Settings) *worker {

	return &worker{
		log: l,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  s.Brokers,
			GroupID:  s.GroupID,
			Topic:    s.Topic,
			MaxWait:  s.MaxWaitTime,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		}),
	}
}

func (w *worker) startWorker() {

	const method = "startWorker"

	w.wg.Add(1)
	w.running = true
	go func() {
		defer w.wg.Done()
		for w.running {
			msg, err := w.reader.ReadMessage(context.Background())
			if err == io.EOF {

				w.log.Error(method, "", "", -34266, err.Error(), "kafka reader is closed", nil)
				break
			}
			if err != nil {
				w.log.Error(method, "", "", -34267, err.Error(), "Read kafka error", nil)
				break
			}

			w.log.Info(method, "", "", "Read kafka message", "Read", string(msg.Value))

			ac := models.AuthEvent{}
			err = json.Unmarshal(msg.Value, &ac)
			if err != nil {
				w.log.Error(method, "", "", -34268, err.Error(), "Ошибка парсинга сообщение от kafka", map[string]string{"kafka message:": string(msg.Value)})
				break
			}
			go func() {
				// Call
				err = events.Default.SendAuthCode(ac)
				err = events.Hook.SendToWebApp(ac, "")
				// finish
			}()
		}
	}()
}

func Init(log service.AppLogger, s Settings) error {
	if err := s.validate(); err != nil {
		return err
	}
	worker := newWoker(log, s)
	// Added kakfa metrics
	readerCollector := collector.NewReaderCollector(worker.reader)
	prometheus.MustRegister(readerCollector)
	// Start
	go worker.startWorker()

	return nil
}

func (r *worker) Close() error {
	// TODO добавить вызов этого метода при получении сигнала о остановке сервиса
	r.running = false
	r.wg.Wait()
	return r.reader.Close()
}
