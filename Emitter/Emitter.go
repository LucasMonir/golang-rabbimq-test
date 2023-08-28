package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed Connecting to RBMQ server...")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open channel...")

	err = ch.ExchangeDeclare(
		"logs_direct",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to declare exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)

	err = ch.PublishWithContext(ctx, "logs_direct", severityFrom(os.Args), false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	failOnError(err, "Failed on publishing messages")

	log.Printf("[X] Message sent: %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "info"
	} else {
		s = os.Args[1]
	}
	return s
}