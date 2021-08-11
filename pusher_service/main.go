package main

import (
	"fmt"
	"github.com/guiaramos/go-meower/event"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalln(err)
	}

	// Connect to Nats
	hub := newHub()
	es, err := event.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		log.Fatal(err)
	}

	// Push messages to clients
	err = es.OnMeowCreated(func(m event.MeowCreatedMessage) {
		log.Printf("Meow received: %v\n", m)
		hub.broadcast(newMeowCreatedMessage(m.ID, m.Body, m.CreatedAt), nil)
	})
	if err != nil {
		log.Fatal(err)
	}

	event.SetEventStore(es)
	defer event.Close()

	// Run WebSocket server
	go hub.run()
	http.HandleFunc("/pusher", hub.handleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
