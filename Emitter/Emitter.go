package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var severity string

func makeRandom() *rand.Rand {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	return random
}

func makeMessages(severity string) string {
	errorMsg := []string{"Ooops, program crashed!", "The execution stopped!"}
	infoMsg := []string{"Excecuton running properly!", "Program working fine!"}

	random := makeRandom()

	if severity == "error" {
		return errorMsg[random.Intn(len(errorMsg))]
	}

	return infoMsg[random.Intn(len(infoMsg))]
}

func makeSeverity() string {
	severities := []string{"info", "error"}
	random := makeRandom()

	return severities[random.Intn(len(severities))]
}

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
	severity = severityFrom(os.Args)
	body := bodyFrom(os.Args)

	err = ch.PublishWithContext(ctx, "logs_direct", severity, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	failOnError(err, "Failed on publishing messages")

	log.Printf("[X] Message sent: %s", body)
}

func bodyFrom(args []string) string {
	var body string

	if (len(args) < 3) || os.Args[2] == "" {
		body = makeMessages(severity)
		print("a")
	} else {
		body = strings.Join(args[2:], " ")
	}

	return body
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = makeSeverity()
	} else {
		s = os.Args[1]
	}

	return s
}
