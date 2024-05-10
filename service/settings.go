package service

type LoggerSettings struct {
	Component string
	MinLevel  string
	Writer    struct {
		Brokers []string
		Topic   string
	}
}
