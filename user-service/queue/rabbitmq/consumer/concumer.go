package consumer

import (
	"github.com/streadway/amqp"
)

type RabbitMQConsumer interface {
	ConsumeMessage(handler func(message []byte)) error
	Close() error
}

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewRabbitMQConsumer(brokerURL, queueName string) (RabbitMQConsumer, error) {
	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Consumer{
		conn:    conn,
		channel: channel,
		queue:   &queue,
	}, nil
}

func (c *Consumer) ConsumeMessage(handler func(message []byte)) error {
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler(msg.Body)
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}

	err = c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// package consumer

// import (
// 	"context"

// 	"github.com/segmentio/kafka-go"
// )

// type ConcumerKafka interface {
// 	ConcumerMessage(handler func(message []byte)) error
// 	Close() error
// }

// type Consumer struct {
// 	reader *kafka.Reader
// }

// func NewKafkaConsumerInit(brokers []string, topic string, groupId string) (ConcumerKafka, error) {
// 	reader := kafka.NewReader(kafka.ReaderConfig{
// 		Brokers: brokers,
// 		Topic:   topic,
// 		GroupID: groupId,
// 	})
// 	return &Consumer{reader: reader}, nil
// }

// func (c *Consumer) ConcumerMessage(handler func(message []byte)) error {

// 	for {
// 		m, err := c.reader.ReadMessage(context.Background())
// 		if err != nil {
// 			return err
// 		}

// 		handler(m.Value)
// 	}
// }

// func (c *Consumer) Close() error {
// 	return c.reader.Close()

// }
