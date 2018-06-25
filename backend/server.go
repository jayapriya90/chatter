package backend

import (
	log "github.com/sirupsen/logrus"
)

// Server maintains the set of active clients and broadcasts messages to the
// clients.
type Server struct {
	// Maintains client to chatroom mapping
	clients map[*Client]string
	// Message to be broadcasted to every other client in the chatroom
	messageChan chan []byte
	// Clients send registration requests (along with client information)
	// in this channel
	registerChan chan *Client
	// Clients send unregistration requests in this channel
	// via explicit disconnect(logout) or socket disconnect (closing the tab etc)
	unregisterChan chan *Client
}

func NewServer() *Server {
	return &Server{
		messageChan:    make(chan []byte),
		registerChan:   make(chan *Client),
		unregisterChan: make(chan *Client),
		clients:        make(map[*Client]string),
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.registerChan:
			s.clients[client] = "globalroom"
			log.Debugf("Registered client. Total clients: %d", len(s.clients))
		case client := <-s.unregisterChan:
			_, exists := s.clients[client]
			if exists {
				delete(s.clients, client)
				log.Debugf("Unregistered client. Total clients: %d", len(s.clients))
				close(client.outboundMessageChan)
			}
		case message := <-s.messageChan:
			log.Debugf("Broadcasting message '%s' to %d clients", message, len(s.clients))
			for client := range s.clients {
				select {
				case client.outboundMessageChan <- message:
				default:
					close(client.outboundMessageChan)
					delete(s.clients, client)
				}
			}
		}
	}
}
