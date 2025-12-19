package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type StockReducedEvent struct {
	MerchantID uint                       `json:"merchant_id"`
	Products   []StockReducedEventProduct `json:"products"`
	OrderID    string                     `json:"order_id"`
	Timestamp  time.Time                  `json:"timestamp"` // Fixed typo: Timestap -> Timestamp
}

type StockReducedEventProduct struct {
	ProductID	uint `json:"product_id"`
	Quantity 	int  `json:"quantity"`
}

type StockConsumer struct {
	conn 			*amqp.Connection
	ch 				*amqp.Channel
}

func NewStockConsumer(url string) (*StockConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 1: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 2: %v", err)
		return nil, err
	}

	// Fixed: business_events (consistent spelling)
	err = ch.ExchangeDeclare(
		"business_events",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 3: %v", err)
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"merchant_stock_events",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 4: %v", err)
		return nil, err
	}

	// Fixed: business_events (consistent spelling)
	err = ch.QueueBind(
		q.Name,
		"merchant.stock.*",
		"business_events",
		false,
		nil,
	)

	if err != nil {
		log.Errorf("[StockConsumer] NewStockConsumer - 5: %v", err)
		return nil, err
	}

	return &StockConsumer{
		conn: conn,
		ch: ch,
	}, nil
}

func (r *RabbitMQService) PublishStockReducedEvent(ctx context.Context, event StockReducedEvent) error{
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = r.ch.Publish(
		"business_events",		//exchange
		"merchant.stock.reduced", //routing-key, lebih spesifik
		false,	//mandatory
		false,	//immidiate
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
			DeliveryMode: amqp.Persistent,
			Timestamp: time.Now(),
		})
	if err != nil{
		return fmt.Errorf("failed to publish event: %v", err)
	}

	return nil
}

func (sc *StockConsumer) Close() error {
	if sc.ch != nil {
		sc.ch.Close()
	}
	if sc.conn != nil {
		return sc.conn.Close()
	}
	return nil
}