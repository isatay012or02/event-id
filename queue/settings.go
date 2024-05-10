package queue

import "time"

// Settings настройки очереди
type Settings struct {
	Brokers     []string
	GroupID     string
	Topic       string
	MaxWaitTime time.Duration `json:"-"`
}

func (s Settings) validate() error {
	return nil
}
