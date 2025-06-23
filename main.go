package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/v1nte/pubsub-go/server"
)

func main() {
	broker := server.NewBroker()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWS(broker, w, r)
	})

	fmt.Println("Server runnin in :9876/ws")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
