package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		"game_logs.*",
		pubsub.QueueTypeDurable,
	)

	for {
		input := gamelogic.GetInput()
		switch input[0] {
		case "pause":

			fmt.Println("Pausing the game...")
			pubsub.PublishJSON(channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
		case "resume":
			fmt.Println("Resuming the game...")
			pubsub.PublishJSON(channel,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
		case "quit":
			fmt.Println("Quitting the game...")
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command, please try again.")
		}
	}

}
