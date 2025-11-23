package server

import (
	"context"
	"log"
	"net/http"

	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
)

type Gateway struct {
	rdb         *redis.Client
	connections map[uuid.UUID]*websocket.Conn
	mutex       sync.RWMutex
	wsUpgrader  websocket.Upgrader
}

type Channel string

const (
	NotificationChannel Channel = "notification_channnel"
	ChatChannnel        Channel = "chat_channel"
	AckChannel          Channel = "ack_channel"
)

type Event string

const (
	ChatEvent     Event = "chat_event"
	AckEvent      Event = "ack_event"
	PresenceEvent Event = "presence_event"
)

func NewGateway(rdb *redis.Client) *Gateway {
	return &Gateway{
		rdb:         rdb,
		connections: make(map[uuid.UUID]*websocket.Conn),
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (g *Gateway) handleConnections(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(UserIDContextKey).(uuid.UUID)
	if !ok {
		http.Error(w, "user_id not found in context", http.StatusUnauthorized)
		return
	}
	conn, err := g.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade failed for user %s: %v", userID, err)
		return
	}
	log.Printf("User %s connected.", userID)
	g.registerConnection(userID, conn)
	defer g.deregisterConnection(userID, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		g.routeIncomingMessage(r.Context(), userID, message)
	}
}

type ClientMessage struct {
	Type    Event           `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func (g *Gateway) routeIncomingMessage(ctx context.Context, userID uuid.UUID, rawMessage []byte) {
	var msg ClientMessage
	if err := json.Unmarshal(rawMessage, &msg); err != nil {
		log.Printf("Invalid message format from %s", userID)
		return
	}

	switch msg.Type {
	case "chat":
		g.rdb.Publish(ctx, string(ChatChannnel), rawMessage)
		log.Printf("Routing chat message from %s", userID)
	case "presence":
		g.rdb.Set(ctx, fmt.Sprintf("user:status:%s", userID), "online", 0)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (g *Gateway) registerConnection(userID uuid.UUID, conn *websocket.Conn) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.connections[userID] = conn
}

func (g *Gateway) deregisterConnection(userID uuid.UUID, conn *websocket.Conn) {
	conn.Close()

	g.mutex.Lock()
	defer g.mutex.Unlock()
	if existingConn, ok := g.connections[userID]; ok && existingConn == conn {
		delete(g.connections, userID)
		log.Printf("User %s disconnected and deregistered.", userID)
	}
}

func (g *Gateway) SubscribeChanel(ctx context.Context, channel Channel, handler func(ctx context.Context, message *redis.Message) error) error {
	pubsub := g.rdb.Subscribe(ctx, string(channel))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		if err := handler(ctx, msg); err != nil {
			log.Printf("Error handling message from channel %s: %v", channel, err)
		}
	}
	return nil
}
