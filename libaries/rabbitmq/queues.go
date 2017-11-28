package rabbitMQ

import (
	"errors"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

//RabbitMQConnect check connect and hold
type RabbitMQConnect struct {
	Conn         *amqp.Connection
	Channel      *amqp.Channel
	ExchangeName string
	ExchangeType string
	Done         chan error
	Uri          string
}

//RabbitConnect connect to RabbitMQ server
func (c *RabbitMQConnect) RabbitConnect() error {

	var err error

	fmt.Println("dialing: ", c.Uri)

	c.Conn, err = amqp.Dial(c.Uri)

	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		fmt.Println("closing: ", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.Done <- errors.New("Channel Closed")
	}()

	fmt.Println("Connecting ok...")

	c.Channel, err = c.Conn.Channel()

	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	return nil
}

//RabbitClose is close all connect
func (c *RabbitMQConnect) RabbitClose() {
	if c.Channel != nil {
		c.Channel.Close()
		c.Channel = nil
	}
	if c.Conn != nil {
		c.Conn.Close()
		c.Conn = nil
	}
}

//RabbitReConnect reconect to rabbitMQ
func (c *RabbitMQConnect) RabbitReConnect() error {

	c.RabbitClose()

	time.Sleep(time.Duration(15) * time.Second)

	fmt.Println("Try ReConnect....")

	err := c.RabbitConnect()

	return err
}

//RabbitCreateExchange use declear exchange
func (c *RabbitMQConnect) RabbitCreateExchange() error {
	if err := c.Channel.ExchangeDeclare(
		c.ExchangeName, // name of the exchange
		c.ExchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}
	return nil
}

//RabbitCreateQueue use create queue
func (c *RabbitMQConnect) RabbitCreateQueue(queueName, routingKey string) error {

	fmt.Println("declared Exchange:", c.ExchangeName, "declaring Queue:", queueName)

	// unbind before bind
	if err := c.Channel.QueueUnbind(
		queueName,      // name of the queue
		routingKey,     // routingKey
		c.ExchangeName, // sourceExchange
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Error Queue Unbind: %s", err)
	}

	queue, err := c.Channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)

	if err != nil {
		return fmt.Errorf("Error Queue Declare: %s", err)
	}

	fmt.Println(fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey))

	/* err = c.Channel.Qos(1, 0, false)
	if err != nil {
		return fmt.Errorf("Error setting qos: %s", err)
	} */

	if err = c.Channel.QueueBind(
		queue.Name,     // name of the queue
		routingKey,     // routingKey
		c.ExchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("Error Queue Bind: %s", err)
	}
	return nil
}
