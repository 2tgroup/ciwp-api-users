package rabbitMQ

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"g.ghn.vn/callback-partner/receive-order/config/machine"

	"github.com/manucorporat/try"
	"github.com/streadway/amqp"
)

// Consumer holds all infromation
// about the RabbitMQ connection
// This setup does limit a consumer
// to one exchange. This should not be
// an issue. Having to connect to multiple
// exchanges means something else is
// structured improperly.
type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string // Name that consumer identifies itself to the server with
	uri          string // uri of the rabbitmq server
	exchange     string // exchange that we will bind to
	exchangeType string // topic, direct, etc...

	lastRecoverTime int64
	//track service current status
	currentStatus atomic.Value
}

const RECOVER_INTERVAL_TIME = 6 * 60

// NewConsumer returns a Consumer struct that has been initialized properly
// essentially don't touch conn, channel, or done and you can create Consumer manually
func newConsumer(consumerTag, uri, exchange, exchangeType string) *Consumer {
	name, err := os.Hostname()
	if err != nil {
		name = "_sim"
	}
	consumer := &Consumer{
		consumerTag:     fmt.Sprintf("%s-%s", consumerTag, name),
		uri:             uri,
		exchange:        exchange,
		exchangeType:    exchangeType,
		done:            make(chan error),
		lastRecoverTime: time.Now().Unix(),
	}
	consumer.currentStatus.Store(true)
	return consumer
}

func maxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func RunConsumer(consumerTag, exchange, exchangeType, queueName, routingKey string, handler func([]byte) bool) {

	config := envConfig.DataConfig

	consumer := newConsumer(consumerTag, config.RabbitMQ, exchange, exchangeType)

	if err := consumer.Connect(); err != nil {
		fmt.Println(err, fmt.Sprintf("[%s]connect error", consumerTag))
	}

	deliveries, err := consumer.AnnounceQueue(queueName, routingKey)

	fmt.Println(err, fmt.Sprintf("[%s]Error when calling AnnounceQueue()", consumerTag))

	consumer.Handle(deliveries, handler, maxParallelism(), queueName, routingKey)
}

// ReConnect is called in places where NotifyClose() channel is called
// wait 30 seconds before trying to reconnect. Any shorter amount of time
// will  likely destroy the error log while waiting for servers to come
// back online. This requires two parameters which is just to satisfy
// the AccounceQueue call and allows greater flexability
func (c *Consumer) ReConnect(queueName, routingKey string, retryTime int) (<-chan amqp.Delivery, error) {

	c.Close()

	time.Sleep(time.Duration(15+rand.Intn(60)+2*retryTime) * time.Second)

	fmt.Println("Try ReConnect with times:", retryTime)

	if err := c.Connect(); err != nil {
		return nil, err
	}

	deliveries, err := c.AnnounceQueue(queueName, routingKey)

	if err != nil {
		return deliveries, errors.New("Couldn't connect")
	}

	return deliveries, nil
}

// Connect to RabbitMQ server
func (c *Consumer) Connect() error {

	var err error

	fmt.Println("dialing: ", c.uri)

	c.conn, err = amqp.Dial(c.uri)

	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Waits here for the channel to be closed
		fmt.Println("closing: ", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Let Handle know it's not time to reconnect
		c.done <- errors.New("Channel Closed")
	}()

	fmt.Println("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	fmt.Println("got Channel, declaring Exchange ", c.exchange)

	if err = c.channel.ExchangeDeclare(
		c.exchange,     // name of the exchange
		c.exchangeType, // type
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

// AnnounceQueue sets the queue that will be listened to for this
// connection...
func (c *Consumer) AnnounceQueue(queueName, routingKey string) (<-chan amqp.Delivery, error) {

	args := make(amqp.Table)

	args["x-message-ttl"] = int32(900000)

	fmt.Println("declared Exchange, declaring Queue:", queueName)

	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		args,      // arguments
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	fmt.Println(fmt.Sprintf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, routingKey))

	// Qos determines the amount of messages that the queue will pass to you before
	// it waits for you to ack them. This will slow down queue consumption but
	// give you more certainty that all messages are being processed. As load increases
	// I would reccomend upping the about of Threads and Processors the go process
	// uses before changing this although you will eventually need to reach some
	// balance between threads, procs, and Qos.
	err = c.channel.Qos(8, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		routingKey, // routingKey
		c.exchange, // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	fmt.Println("Queue bound to Exchange, starting Consume consumer tag:", c.consumerTag)
	deliveries, err := c.channel.Consume(
		queue.Name,    // name
		c.consumerTag, // consumerTag,
		false,         // noAck
		false,         // exclusive
		false,         // noLocal
		false,         // noWait
		nil,           // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}
	return deliveries, nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
		c.channel = nil
	}
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Consumer) Handle(
	deliveries <-chan amqp.Delivery,
	fn func([]byte) bool,
	threads int,
	queue string,
	routingKey string) {

	var err error
	for {
		fmt.Println("Enter for busy loop with thread:", threads)
		for i := 0; i < threads; i++ {
			go func() {
				fmt.Println("Enter go with thread with deliveries", deliveries)
				for msg := range deliveries {
					//fmt.Println("Enter deliver")
					ret := false
					try.This(func() {
						body := msg.Body[:]
						ret = fn(body)
					}).Finally(func() {
						if ret == true {
							msg.Ack(false)
							currentTime := time.Now().Unix()
							if currentTime-c.lastRecoverTime > RECOVER_INTERVAL_TIME && !c.currentStatus.Load().(bool) {
								fmt.Println("Try to Recover Unack Messages!")
								c.currentStatus.Store(true)
								c.lastRecoverTime = currentTime
								c.channel.Recover(true)
							}
						} else {
							// this really a litter dangerous. if the worker is panic very quickly,
							// it will ddos our sentry server......plz, add [retry-ttl] in header.
							//msg.Nack(false, true)
							c.currentStatus.Store(false)
						}
					}).Catch(func(e try.E) {
						fmt.Println(e)
					})
				}
			}()
		}

		// Go into reconnect loop when
		// c.done is passed non nil values
		if <-c.done != nil {
			c.currentStatus.Store(false)
			retryTime := 1
			for {
				deliveries, err = c.ReConnect(queue, routingKey, retryTime)
				if err != nil {
					fmt.Println(err, "Reconnecting Error")
					retryTime += 1
				} else {
					break
				}
			}
		}
		fmt.Println("Reconnected!!!")
	}
}
