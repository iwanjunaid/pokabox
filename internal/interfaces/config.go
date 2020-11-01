package interfaces

type Config interface {
	GetOutboxName() string
	GetPollInterval() int
}
