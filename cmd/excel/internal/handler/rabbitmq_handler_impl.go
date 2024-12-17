package handler

import (
	"log"

	"github.com/ryvasa/go-super-farmer/cmd/excel/internal/usecase"
	"github.com/ryvasa/go-super-farmer/pkg/messages"
)

type RabbitMQHandlerImpl struct {
	rabbitMQUsecase usecase.RabbitMQUsecase
	excelService    usecase.ExcelInterface
	rabbitMQ        messages.RabbitMQ
}

func NewRabbitMQHandler(rabbitMQUsecase usecase.RabbitMQUsecase, excelSvc usecase.ExcelInterface, rabbitMQ messages.RabbitMQ) RabbitMQHandler {
	return &RabbitMQHandlerImpl{rabbitMQUsecase, excelSvc, rabbitMQ}
}

func (h *RabbitMQHandlerImpl) ConsumerHandler() error {
	prices, err := h.rabbitMQ.ConsumeMessages("price-history-queue")
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range prices {
			if err := h.rabbitMQUsecase.HandlePriceHistoryMessage(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	harvest, err := h.rabbitMQ.ConsumeMessages("harvest-queue")
	if err != nil {
		log.Fatal(err)
	}

	forever = make(chan bool)

	go func() {
		for d := range harvest {
			if err := h.rabbitMQUsecase.HandleHarvestMessage(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}
