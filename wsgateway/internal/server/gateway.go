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
	NotificationChannel     Channel = "notification_channel"
	ChatIncomingChannel     Channel = "chat_incoming"
	ChatOutgoingChannel     Channel = "chat_outgoing"
	ReadIncomingChannel     Channel = "read_incoming"
	ReadOutgoingChannel     Channel = "read_outgoing"
	AckChannel              Channel = "ack_channel"
	PresenceIncomingChannel Channel = "presence_incoming"
	PresenceOutgoingChannel Channel = "presence_outgoing"
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
		if err := g.rdb.Publish(ctx, string(ChatIncomingChannel), string(payloadBytes)).Err(); err != nil {
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
		if err := g.rdb.Publish(ctx, string(ReadIncomingChannel), string(payloadBytes)).Err(); err != nil {
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
	g.rdb.Publish(ctx, string(PresenceIncomingChannel), fmt.Sprintf(`{"user_id":"%s","status":"online"}`, userID))
	log.Printf("User %s registered.", userID)
}

func (g *Gateway) deregisterConnection(ctx context.Context, userID uuid.UUID, conn *websocket.Conn) {
	conn.Close()

	g.mutex.Lock()
	defer g.mutex.Unlock()
	if existingConn, ok := g.connections[userID]; ok && existingConn == conn {
		delete(g.connections, userID)
		g.rdb.Del(ctx, fmt.Sprintf("user:status:%s", userID))
		g.rdb.Publish(ctx, string(PresenceIncomingChannel), fmt.Sprintf(`{"user_id":"%s","status":"offline"}`, userID))
		log.Printf("User %s disconnected and deregistered.", userID)
	}
}

func (g *Gateway) SubscribeChannel(ctx context.Context, channel Channel, handler func(ctx context.Context, msg *redis.Message) error) error {
	pubsub := g.rdb.Subscribe(ctx, string(channel))
	ch := pubsub.Channel()

	go func() {
		defer pubsub.Close()
		for {
			select {
			case <-ctx.Done():
				log.Printf("Context cancelled. Stopping subscription for channel: %s", channel)
				return

			case msg, ok := <-ch:
				if !ok {
					log.Printf("Redis channel closed for %s.", channel)
					return
				}

				if err := handler(ctx, msg); err != nil {
					log.Printf("Error handling message from channel %s: %v", channel, err)
				}
			}
		}
	}()
	return nil
}
