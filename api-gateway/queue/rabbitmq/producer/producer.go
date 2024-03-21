package producer

import (
	"github.com/streadway/amqp"
)

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQProducer(amqpURI string) (*RabbitMQProducer, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, err
	}

	// Open a channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQProducer) ProducerMessage(queueName string, message []byte) error {
	// Declare a queue
	_, err := r.channel.QueueDeclare(
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
	// Publish a message
	err = r.channel.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})
	return err
}

func (r *RabbitMQProducer) Close() {
	r.channel.Close()
	r.conn.Close()
}
