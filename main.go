package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/v1nte/pubsub-go/database"
	"github.com/v1nte/pubsub-go/logger"
	"github.com/v1nte/pubsub-go/server"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatal("Could not connect to DB")
	}

	defer database.Close()

	if err := logger.Init(); err != nil {
		log.Fatal("Logger couldn't init")
	}

	logger.Log.Info("App Started")
	defer logger.Log.Info("App shutdown")
	defer logger.Log.Sync()

	broker := server.NewBroker()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.HandleWS(broker, w, r)
	})

	fmt.Println("Server runnin in :9876/ws")
	log.Fatal(http.ListenAndServe(":9876", nil))
}
