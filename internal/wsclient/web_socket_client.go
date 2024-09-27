package wsclient

import (
	"github.com/gorilla/websocket"
	"github.com/qwertydi/go-challenge/internal/aggregator"
	"log"
)

// WebSocketHandler interface defines the methods for handling WebSocket events
type WebSocketHandler interface {
	OnConnect(conn *websocket.Conn)
	OnMessage(messageType int, message []byte)
	OnError(err error)
	OnClose(code int, text string)
}

// WebSocketClient struct to manage the connection and handler
type WebSocketClient struct {
	URL         string
	Conn        *websocket.Conn
	Handler     WebSocketHandler
	DataService aggregator.DataServiceHandler
}

// WebSocketHandlerImpl provides a basic implementation of WebSocketHandler
type WebSocketHandlerImpl struct {
	DataServiceHandler aggregator.DataServiceHandler
	AggregateHandler   aggregator.AggregateServiceHandler
}

func (h *WebSocketHandlerImpl) OnConnect(conn *websocket.Conn) {
	log.Println("Connected to server")
}

func (h *WebSocketHandlerImpl) OnMessage(messageType int, message []byte) {
	// slog.Printf("Received message: %s\n", string(message))

	h.DataServiceHandler.ProcessMessage(message)

}

func (h *WebSocketHandlerImpl) OnError(err error) {
	log.Printf("Error: %v\n", err)
}

func (h *WebSocketHandlerImpl) OnClose(code int, text string) {
	log.Printf("Connection closed: %d %s\n", code, text)
}

// Connect establishes a connection to the WebSocket server
func (c *WebSocketClient) Connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.URL, nil)
	if err != nil {
		return err
	}
	c.Conn = conn
	c.Handler.OnConnect(conn)
	return nil
}

// Listen starts listening for messages from the server
func (c *WebSocketClient) Listen() {
	defer c.Conn.Close()

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			c.Handler.OnError(err)
			return
		}
		c.Handler.OnMessage(messageType, message)
	}
}

// SendMessage sends a message to the server
func (c *WebSocketClient) SendMessage(message string) error {
	log.Printf("Sending message: %s\n", message)

	return c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

// NewWebSocketClient creates a new WebSocketClient
func NewWebSocketClient(url string, handler WebSocketHandler) *WebSocketClient {
	return &WebSocketClient{
		URL:     url,
		Handler: handler,
	}
}
