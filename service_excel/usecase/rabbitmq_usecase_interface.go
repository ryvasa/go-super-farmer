package usecase

type RabbitMQUsecase interface {
	HandlePriceHistoryMessage(msgBody []byte) error
	HandleHarvestMessage(msgBody []byte) error
}
