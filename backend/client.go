package backend

import (
	"bytes"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// Maximum message size allowed from chat server
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	server *Server
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel for outbound messages.
	outboundMessageChan chan []byte
}

// readMessages reads messages from the websocket connection
func (c *Client) readMessages() {
	// On exit, unregister client and close connection
	defer func() {
		c.server.unregisterChan <- c
		c.conn.Close()
	}()
	// Set read limit size for messages (1 KB here)
	c.conn.SetReadLimit(maxMessageSize)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		c.server.messageChan <- message
	}
}

// writeMessages writes messages to the websocket connection.
func (c *Client) writeMessages() {
	// On exit, unregister client and close connection
	defer func() {
		c.server.unregisterChan <- c
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.outboundMessageChan:
			if !ok {
				// The server closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.outboundMessageChan)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.outboundMessageChan)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

// ServeWebSocket handles websocket requests to the chat server.
func ServeWebSocket(server *Server, w http.ResponseWriter, r *http.Request) {
	// HTTP -> WebSocket upgrade
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}
	client := &Client{server: server, conn: conn, outboundMessageChan: make(chan []byte, 1024)}
	// Send client information to registerChannel to register client with the server
	client.server.registerChan <- client

	// Start reader and writer goroutines
	go client.writeMessages()
	go client.readMessages()
}
