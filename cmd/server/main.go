package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

var CONNECTION_STRING = "amqp://guest:guest@localhost:5672/"

func main() {
	fmt.Println("Starting Peril server...")

	conn, err := amqp.Dial(CONNECTION_STRING)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	fmt.Println("Opened channel")
	defer channel.Close()

	pubsub.PublishJSON(channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt) // provide channel and signals to listen to, otherwise all signals will be relayed
	<-ch                            // block until signal is received
	fmt.Println("Received interrupt signal, exiting...")
}
