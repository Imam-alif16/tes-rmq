package configs

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() (*amqp091.Connection, error) {
	rabbitMQURI := EnvRabbitMQURI()
	if rabbitMQURI == "" {
		return nil, fmt.Errorf("RABBITMQ_URL not set in environment variables")
	}

	mqConn, err := amqp091.Dial(rabbitMQURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	fmt.Println("Connected to RabbitMQ")

	return mqConn, nil
}

func ChannelRabbitMQ(mqConn *amqp091.Connection) (*amqp091.Channel, error) {
	channel, err := mqConn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	fmt.Println("Channel created")

	return channel, nil
}

func ExchangeDeclareRabbitMQ(channel *amqp091.Channel, exchangeName, exchangeType string) error {
	err := channel.ExchangeDeclare(
		exchangeName, // Name of the exchange
		exchangeType, // Type of the exchange: "direct", "fanout", "topic", etc.
		false,        // Durable
		false,        // AutoDelete
		false,        // Internal
		false,        // NoWait
		nil,          // Arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a exchange: %v", err)
	}

	fmt.Printf("Exchange %s declared\n", exchangeName)

	return nil
}

func QueueDeclareRabbitMQ(channel *amqp091.Channel, queueName string) (*amqp091.Queue, error) {
	queue, err := channel.QueueDeclare(
		queueName,						// name
		false,								// durable
		false,								// delete when unused
		false,								// exclusive
		false,								// no-wait
		nil,									// arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %v", err)
	}

	fmt.Printf("Queue %s declared\n", queueName)

	return &queue, nil
}

func QueueBindRabbitMQ(channel *amqp091.Channel, queueName, routingKey, exchangeName string) error {
	err := channel.QueueBind(
		queueName,        // queue name
		routingKey, // routing key
		exchangeName, // exchange name
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		return fmt.Errorf("failed to bind a queue: %v", err)
	}

	fmt.Println("Queue Bound")

	return nil
}