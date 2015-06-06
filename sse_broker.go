// Heavily based on the excellent resource from Bamboo
// https://www.new-bamboo.co.uk/blog/2014/05/13/writing-a-server-sent-events-server-in-go/

package main

import (
	"fmt"
	"log"
	"net/http"
)

// A Broker:
// * holds open client connections
// * listens for incoming events on its Notifier channel
// * broadcasts event data to all registered connections
type Broker struct {
	// Events are pushed to this channel
	// Example: broker.Notifier <- []byte(string)
	Notifier chan []byte

	// A channel that takes and registers newly open HTTP connections
	newClients chan chan []byte

	// Will be notified when a client connection closes
	closingClients chan chan []byte

	// Main registry for client connections
	clients map[chan []byte]bool
}

func (broker *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// Make sure our response writer supports the http.Flusher interface
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Brokers
	// connections registry
	messageChan := make(chan []byte)

	// Signal the broker that we have a new connection
	broker.newClients <- messageChan

	// If for any reason our connection dies, we want to make sure our broker is
	// notified so we can clean up
	defer func() {
		broker.closingClients <- messageChan
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		broker.closingClients <- messageChan
	}()

	// block waiting for messages broadcast on this connections messageChan
	for {
		fmt.Fprintf(rw, "data: %s\n\n", <-messageChan)

		// Flush the data immediately instead of buffering it for later
		flusher.Flush()
	}
}

func (broker *Broker) listen() {
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			log.Printf("Client added. %d register clients\n", len(broker.clients))

		case s := <-broker.closingClients:
			delete(broker.clients, s)
			log.Printf("Removed client. %d register clients\n", len(broker.clients))

		// We received a new event, send to all connected clients
		case event := <-broker.Notifier:
			for clientMessageChan, _ := range broker.clients {
				clientMessageChan <- event
			}
		}
	}
}

func NewServer() (broker *Broker) {
	broker = &Broker{
		Notifier:       make(chan []byte, 1),
		newClients:     make(chan chan []byte),
		closingClients: make(chan chan []byte),
		clients:        make(map[chan []byte]bool),
	}

	go broker.listen()

	return
}
