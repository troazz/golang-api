package stores

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection

// Init rabbitmq connection
func Init() error {
	var err error
	// Initialize the package level "conn" variable that represents the connection the the rabbitmq server
	for {
		conn, err = amqp.Dial("amqp://rabbitmq:rabbitmq@rabbitmq:5672/")
		if err != nil {
			log.Println(err)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	return nil
}

// Queue news request to rabbitmq
func Queue(qName string, data interface{}) error {
	err := Init()
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.New("RMQ: Failed to open channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		qName, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.New("RMQ: Failed to declare queue")
	}

	body, err := json.Marshal(data)
	if err != nil {
		return errors.New("RMQ: Failed to encode json object")
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return errors.New("RMQ: Failed to publish")
	}

	log.Printf("Queued: %s\n", data)

	return nil
}

// Subscribe channel to require new News
func Subscribe(qName string) (<-chan amqp.Delivery, func(), error) {
	err := Init()
	if err != nil {
		return nil, nil, err
	}

	// create a channel through which we publish
	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}
	// assert that the queue exists (creates a queue if it doesn't)
	q, err := ch.QueueDeclare(
		qName, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	c, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	// return the created channel
	return c, func() {
		conn.Close()
		ch.Close()
	}, err
}
