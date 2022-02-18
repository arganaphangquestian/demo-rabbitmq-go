package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

type Config struct {
	RABBITMQ_DEFAULT_HOST string
	RABBITMQ_DEFAULT_USER string
	RABBITMQ_DEFAULT_PASS string
}

var config Config

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Can't load .env")
	}
	config = Config{
		RABBITMQ_DEFAULT_HOST: os.Getenv("RABBITMQ_DEFAULT_HOST"),
		RABBITMQ_DEFAULT_USER: os.Getenv("RABBITMQ_DEFAULT_USER"),
		RABBITMQ_DEFAULT_PASS: os.Getenv("RABBITMQ_DEFAULT_PASS"),
	}
}

func main() {

	// args
	topic := flag.String("topic", "default", "topic name")
	flag.Parse()

	// Queue Connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.RABBITMQ_DEFAULT_USER, config.RABBITMQ_DEFAULT_PASS, config.RABBITMQ_DEFAULT_HOST))
	if err != nil {
		log.Fatal("Cannot connect to rabbitmq")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Cannot connect to channel")
	}
	defer ch.Close()

	// Declare Queue
	q, err := ch.QueueDeclare(
		*topic, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		log.Fatal("Cannot declare Queue")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("Failed consume queue")
	}

	quit := make(chan struct{})
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-quit
}
