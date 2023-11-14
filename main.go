package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	msgbroker "github.com/carlosm27/go-rabbitmq-demo/pkg/msg-broker"
	"github.com/carlosm27/go-rabbitmq-demo/pkg/msg-broker/rabbitmq"
)

type MsgBroker interface {
	PublishMessage(ctx context.Context, queueSettings msgbroker.QueueSettings, msgBody []byte) error
}

var msgBroker MsgBroker

func main() {

	mux := http.NewServeMux()
	port := ":8000"

	// Create the rabbitMQ message broker implementation
	msgBroker, onClose := rabbitmq.MustGetNewRabbitMqBroker("amqp://guest:guest@localhost:5672/")
	defer onClose()

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte("Hello this is a home page"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Declare the queue settings
		queueSettings := msgbroker.QueueSettings{Name: "hello"}

		// publish the rabbitMQ message
		body := "Hello World!"
		err := msgBroker.PublishMessage(ctx, queueSettings, []byte(body))
		if err != nil {
			panic(fmt.Errorf("failed to publish a message: %v", err))
		}
		log.Printf(" [x] Sent %s\n", body)
	}))

	slog.Info("Listening on ", "port", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		slog.Warn("Problem starting the server", "error", err)
	}

}
