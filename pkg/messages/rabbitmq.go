package rabbitmq

import (
	"fmt"

	"github.com/ryvasa/go-super-farmer/pkg/env"
)

func NewRabbitMQURL(env *env.Env) string {
    return fmt.Sprintf("amqp://%s:%s@%s:%s/",
        env.RabbitMQ.User,
        env.RabbitMQ.Password,
        env.RabbitMQ.Host,
        env.RabbitMQ.Port,
    )
}
