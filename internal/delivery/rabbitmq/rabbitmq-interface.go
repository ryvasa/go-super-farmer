package rabbitmq

import "context"

type Publisher interface {
	Publish(ctx context.Context, queueName string, body []byte) error
	PublishJSON(ctx context.Context, queueName string, data interface{}) error
}
