package rabbitmq

import (
	"context"
	"encoding/json"
	"warehouse-go/notification-service/configs"
	"warehouse-go/notification-service/pkg/email"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQServiceInterface interface {
	ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error
	Close() error
}

type rabbitMQService struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  configs.Config
}

// Close implements RabbitMQServiceInterface.
func (r *rabbitMQService) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// ConsumeEmail implements RabbitMQServiceInterface.
func (r *rabbitMQService) ConsumeEmail(ctx context.Context, emailService email.EmailServiceInterface) error {
	msgs, err := r.channel.Consume(
		"email_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[RabbitMQService] ConsumeEmail - 1: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case <- ctx.Done():
				log.Errorf("Email consumer context cancelled")
				return
			case msg := <- msgs:
				var emailPayload email.EmailPayload
				if err := json.Unmarshal(msg.Body, &emailPayload); err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 2: %v", err)
					msg.Nack(false, false)
					continue
				} 
				//Process email based on type
				var err error
				switch emailPayload.Type {
				case "welcome", "welcome_email":
					err = emailService.SendWelcomeEmail(ctx, emailPayload)
				default:
					log.Errorf("[RabbitMQService] ConsumeEmail - 3: %s", "unknown email type")
					msg.Nack(false, false)
					continue
				}

				if err != nil {
					log.Errorf("[RabbitMQService] ConsumeEmail - 4: %v", err)
					msg.Nack(false, true) //requeue
				} else {
					log.Infof("[RabbitMQService] ConsumeEmail - 5: %s", "Email sent successfully")
					msg.Ack(false)
				}
			}
		}
	}()

	return nil 
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := amqp.Dial(config.RabitMQ.URL())
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 1: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[RabbitMQService] NewRabbitMQService - 2: %v", err)
		return nil, err
	}

	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		config:  config,
	}, nil
}
