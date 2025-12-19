package service

import (
	"context"
	"encoding/json"
	"fmt"
	"warehouse-go/user-service/configs"


	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)
type EmailPayload struct {
	Email    	string `json:"email"`
	Password 	string `json:"password"`
	Type 	  	string `json:"type"`
	UserID   	uint   `json:"user_id"`
	Name    	string `json:"name"`
}

type RabbitMQServiceInterface interface {
	PublishEmail(ctx context.Context, payload EmailPayload) error
	Close() error
}

type rabbitMQService struct {
	conn *amqp.Connection
	ch *amqp.Channel
	config configs.Config
}


// ────────────────────────────────────────────────────────────────
// Close implements RabbitMQServiceInteface
// ────────────────────────────────────────────────────────────────
func (r *rabbitMQService) Close() error {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

// ────────────────────────────────────────────────────────────────
// PublishEmail implements RabbitMQServiceInteface
// ────────────────────────────────────────────────────────────────	
func (r *rabbitMQService) PublishEmail(ctx context.Context, payload EmailPayload) error {
		//Conver payload ke JSON
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal email payload: %v", err)
	}

	//Decalre queue if not exists
	queue, err := r.ch.QueueDeclare(
		"email_queue", //name
		true,		  //durable
		false,        //delete when unused
		false,        //exclusive
		false,        //no-wait
		nil,          //arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare email queue: %v", err)
	}

	//Publish ke email queue langsung (tanpa exchange)
	err = r.ch.Publish(
		"",		 //exchange
		queue.Name, //routing key (nama queue)
		false,	 //mandatory
		false,	 //immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish email message: %v", err)
	}

	return nil
}

func NewRabbitMQService(config configs.Config) (RabbitMQServiceInterface, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
    config.RabitMQ.Username,
    config.RabitMQ.Password,
    config.RabitMQ.Host,
    config.RabitMQ.Port,
))
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
		conn:   conn,
		ch:     ch,
		config: config,
	}, nil
}