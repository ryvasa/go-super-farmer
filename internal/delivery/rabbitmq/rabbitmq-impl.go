package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

type publisher struct {
	conn *amqp091.Connection
}

func NewPublisher(env *env.Env) (Publisher, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		env.RabbitMQ.User,
		env.RabbitMQ.Password,
		env.RabbitMQ.Host,
		env.RabbitMQ.Port,
	)

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	return &publisher{
		conn: conn,
	}, nil
}

func (p *publisher) Publish(ctx context.Context, queueName string, body []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}
func (p *publisher) PublishJSON(ctx context.Context, queueName string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		})
}
