package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Missatge per difusió
type Message struct {
	Channel string      `json:"channel"` // "general" o "workcenter:{id}"
	Data    interface{} `json:"data"`
}

type Hub struct {
	clients    map[*websocketClient]bool
	register   chan *websocketClient
	unregister chan *websocketClient
	broadcast  chan Message
	mu         sync.Mutex
}

// struct que envolta cada connexió per saber a quin canal està
type websocketClient struct {
	conn    *websocket.Conn
	channel string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewHub() *Hub {
	h := &Hub{
		clients:    make(map[*websocketClient]bool),
		register:   make(chan *websocketClient),
		unregister: make(chan *websocketClient),
		broadcast:  make(chan Message, 100),
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.conn.Close()
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			data, err := json.Marshal(msg.Data)
			if err != nil {
				log.Printf("error marshalling message: %v", err)
				continue
			}
			h.mu.Lock()
			for c := range h.clients {
				if c.channel == msg.Channel || msg.Channel == "general" {
					if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
						c.conn.Close()
						delete(h.clients, c)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// Upgrade a websocket i registra el client
func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request, channel string, initialData interface{}) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &websocketClient{conn: conn, channel: channel}
	h.register <- client

	// envia estat inicial
	if initialData != nil {
		if data, err := json.Marshal(initialData); err == nil {
			conn.WriteMessage(websocket.TextMessage, data)
		}
	}

	// llegir per detectar tancament
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.unregister <- client
			break
		}
	}
}

// difondre missatge a un canal concret o a "general"
func (h *Hub) Broadcast(channel string, msg interface{}) {
	select {
	case h.broadcast <- Message{Channel: channel, Data: msg}:
	default:
		log.Println("broadcast channel full, dropping message")
	}
}
