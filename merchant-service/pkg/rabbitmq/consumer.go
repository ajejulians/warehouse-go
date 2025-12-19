package rabbitmq

import (
	"context"
	"encoding/json"
	"time"
	"warehouse-go/merchant-service/repository"

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
	merchantRepo 	repository.MerchantProductRepositoryInterface
}

func NewStockConsumer(url string, merchantRepo repository.MerchantProductRepositoryInterface) (*StockConsumer, error) {
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
		merchantRepo: merchantRepo,
	}, nil
}


func (s *StockConsumer) ConsumeStockReductionEvent(ctx context.Context) error {
	msgs, err := s.ch.Consume(
		"merchant_stock_events",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	
	if err != nil {
		log.Errorf("[StockConsumer] ConsumeStockReductionEvent - 1: %v", err)
		return err // Fixed: return err instead of nil
	}

	for {
		select {
		case <- ctx.Done():
			log.Info("Stopping stock consumer...")
			return nil
		case msg := <-msgs:
			go s.handleStockReductionEvent(msg)
		}
	}
}

func (sc *StockConsumer) handleStockReductionEvent(msg amqp.Delivery) error {
	defer msg.Ack(false)

	var event StockReducedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Errorf("[StockConsumer] handleStockReductionEvent - 1: %v", err)
		return err
	}

	for _, product := range event.Products {
		if err := sc.reduceStock(event.MerchantID, product.ProductID, product.Quantity); err != nil {
			log.Errorf("[StockConsumer] handleStockReductionEvent - 2: %v", err)
			continue
		}

		log.Infof("Successfully reduced stock for product %d by %d", product.ProductID, product.Quantity)
	}

	return nil
}

func (sc *StockConsumer) reduceStock(merchantID uint, productID uint, quantity int) error {
	err := sc.merchantRepo.ReduceStock(context.Background(), merchantID, productID, int64(quantity))
	if err != nil {
		log.Errorf("[StockConsumer] reduceStock - 1: %v", err)
		return err
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