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
	message := flag.String("message", "", "message data")
	flag.Parse()

	// Queue Connection
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", config.RABBITMQ_DEFAULT_USER, config.RABBITMQ_DEFAULT_PASS, config.RABBITMQ_DEFAULT_HOST))
	if err != nil {
		log.Fatal("Cannot connect to rabbitmq")
	}
	defer conn.Close()

	// create channel to queue
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Cannot create channel to rabbitmq")
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

	// Send Data
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(*message),
		})

	if err != nil {
		log.Fatal("Data failed to send")
	}
	fmt.Println("Message sent")
}
