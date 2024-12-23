package handler

type RabbitMQHandler interface {
	ConsumerHandler() error
}
