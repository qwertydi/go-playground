package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/qwertydi/go-challenge/internal/aggregator"
	"github.com/qwertydi/go-challenge/internal/util"
	"github.com/qwertydi/go-challenge/internal/wsclient"
	"log"
	"os"
	"os/signal"
	"time"
)

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value and a boolean indicating whether it was found.
// If the variable is not found, it returns the default value.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	sHost := GetEnv("SOCKET_HOST", "localhost")
	sPort := GetEnv("SOCKET_PORT", "5050")

	timeService := util.TimeService()
	dataService := aggregator.DataService(timeService)
	aggregateService := aggregator.AggregateService(dataService)

	handler := &wsclient.WebSocketHandlerImpl{
		DataServiceHandler: dataService,
		AggregateHandler:   aggregateService,
	}

	connection := "ws://" + sHost + ":" + sPort + "/ws"
	client := wsclient.NewWebSocketClient(connection, handler)

	err := client.Connect()
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}

	// Start listening for messages in a separate goroutine
	go client.Listen()

	// Send a message every 5 seconds
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for {
			<-ticker.C
			data := handler.AggregateHandler.AggregateData(time.Now())
			if len(data) != 0 {
				marshal, err := json.Marshal(data)
				if err != nil {
					return
				}
				if err := client.SendMessage(string(marshal)); err != nil {
					log.Println("Error sending message:", err)
				}
			}
		}
	}()

	// Wait for interrupt signal to gracefully close the connection
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	// Close the connection
	err = client.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Error during closing websocket:", err)
	}
	client.Handler.OnClose(websocket.CloseNormalClosure, "")
}
