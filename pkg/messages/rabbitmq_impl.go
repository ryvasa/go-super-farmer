package messages

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/ryvasa/go-super-farmer/pkg/env"
)

type RabbitMQImpl struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
	env        *env.Env
}

func NewRabbitMQ(env *env.Env) (RabbitMQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		env.RabbitMQ.User,
		env.RabbitMQ.Password,
		env.RabbitMQ.Host,
		env.RabbitMQ.Port,
	)

	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("error creating channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"report-exchange", // name
		"direct",          // type
		true,              // durable
		false,             // auto-deleted
		false,             // internal
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"verify-email-exchange", // name
		"direct",                // type
		true,                    // durable
		false,                   // auto-deleted
		false,                   // internal
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		return nil, err
	}

	queues := []string{"price-history-queue", "harvest-queue", "verify-email-queue"}
	for _, queueName := range queues {
		_, err = ch.QueueDeclare(
			queueName,
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			return nil, err
		}
	}

	err = ch.QueueBind(
		"price-history-queue", // queue name
		"price-history",       // routing key
		"report-exchange",     // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		"harvest-queue",   // queue name
		"harvest",         // routing key
		"report-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		"verify-email-queue",    // queue name
		"verify-email",          // routing key
		"verify-email-exchange", // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQImpl{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (r *RabbitMQImpl) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return r.Channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func (r *RabbitMQImpl) PublishJSON(ctx context.Context, exchange, routingKey string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return r.Channel.PublishWithContext(ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		})
}

func (r *RabbitMQImpl) DeclareQueue(name string) (amqp091.Queue, error) {
	return r.Channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQImpl) ConsumeMessages(queueName string) (<-chan amqp091.Delivery, error) {
	return r.Channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (r *RabbitMQImpl) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Connection != nil {
		r.Connection.Close()
	}
}
