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
	ReadChannel         Channel = "read_channel"
	AckChannel          Channel = "ack_channel"
	PresenceChannel     Channel = "presence_channel"
)

type Event string

const (
	ChatEvent         Event = "chat_event"
	ReadEvent         Event = "read_event"
	AckEvent          Event = "ack_event"
	PresenceEvent     Event = "presence_event"
	NotificationEvent Event = "notification_event"
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
	g.registerConnection(r.Context(), userID, conn)
	defer g.deregisterConnection(r.Context(), userID, conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		g.routeIncomingMessage(r.Context(), userID, message)
	}
}

// client -> websocket: chat

// server -> redis -> websocket: notification ack chat presence

// websocket -> redis -> server: chat presence
// websocket -> client: notification ack chat presence

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
	case ChatEvent:
		var payload MessagePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("Invalid chat message payload from %s", userID)
			return
		}
		payload.SenderID = userID
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal chat message payload from %s", userID)
			return
		}
		if err := g.rdb.Publish(ctx, string(ChatChannnel), string(payloadBytes)).Err(); err != nil {
			log.Printf("Failed to publish chat message from %s: %v", userID, err)
			return
		}
		log.Printf("Routing chat message from %s", userID)
	case ReadEvent:
		var payload ReadPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			log.Printf("Invalid read message payload from %s", userID)
			return
		}
		payload.UserID = userID
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Failed to marshal read message payload from %s", userID)
			return
		}
		if err := g.rdb.Publish(ctx, string(ReadChannel), string(payloadBytes)).Err(); err != nil {
			log.Printf("Failed to publish read message from %s: %v", userID, err)
			return
		}
		log.Printf("Routing read message from %s", userID)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

func (g *Gateway) registerConnection(ctx context.Context, userID uuid.UUID, conn *websocket.Conn) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.connections[userID] = conn
	g.rdb.Set(ctx, fmt.Sprintf("user:status:%s", userID), "online", 0)
	g.rdb.Publish(ctx, string(PresenceEvent), fmt.Sprintf(`{"user_id":"%s","status":"online"}`, userID))
}

func (g *Gateway) deregisterConnection(ctx context.Context, userID uuid.UUID, conn *websocket.Conn) {
	conn.Close()

	g.mutex.Lock()
	defer g.mutex.Unlock()
	if existingConn, ok := g.connections[userID]; ok && existingConn == conn {
		delete(g.connections, userID)
		g.rdb.Del(ctx, fmt.Sprintf("user:status:%s", userID))
		g.rdb.Publish(ctx, string(PresenceEvent), fmt.Sprintf(`{"user_id":"%s","status":"offline"}`, userID))
		log.Printf("User %s disconnected and deregistered.", userID)
	}
}

func (g *Gateway) SubscribeChanel(ctx context.Context, channel Channel, handler func(ctx context.Context, message *redis.Message) error) {
	pubsub := g.rdb.Subscribe(ctx, string(channel))
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		if err := handler(ctx, msg); err != nil {
			log.Printf("Error handling message from channel %s: %v", channel, err)
		}
	}
}
