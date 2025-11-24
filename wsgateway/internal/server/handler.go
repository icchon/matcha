package server

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"time"
)

type NotificationPayload struct {
	ID          int64     `json:"id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	SenderID    string    `json:"sender_id"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func (g *Gateway) NotificationHandler(ctx context.Context, message *redis.Message) error {
	var notification NotificationPayload
	if err := json.Unmarshal([]byte(message.Payload), &notification); err != nil {
		return err
	}
	userID := notification.RecipientID

	g.mutex.RLock()
	conn, ok := g.connections[userID]
	g.mutex.RUnlock()

	if !ok {
		log.Printf("Push: User %s not connected (offline).", userID)
		return nil
	}

	data := ClientMessage{
		Type:    NotificationEvent,
		Payload: json.RawMessage(message.Payload),
	}
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Push failed for user %s: %v", userID, err)
		g.deregisterConnection(ctx, userID, conn)
	} else {
		log.Printf("Successfully pushed notification to user %s.", userID)
	}
	return nil
}

type MessagePayload struct {
	ID          int64     `json:"id"`
	SenderID    uuid.UUID `json:"sender_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
	SentAt      time.Time `json:"sent_at"`
}

func (g *Gateway) ChatMessageHandler(ctx context.Context, message *redis.Message) error {
	log.Printf("Received chat message: %s", message.Payload)
	var chatMsg MessagePayload
	if err := json.Unmarshal([]byte(message.Payload), &chatMsg); err != nil {
		return err
	}
	recipientID := chatMsg.RecipientID

	g.mutex.RLock()
	conn, ok := g.connections[recipientID]
	g.mutex.RUnlock()

	if !ok {
		log.Printf("Chat: User %s not connected (offline).", recipientID)
		return nil
	}

	data := ClientMessage{
		Type:    ChatEvent,
		Payload: json.RawMessage(message.Payload),
	}

	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Chat message push failed for user %s: %v", recipientID, err)
		g.deregisterConnection(ctx, recipientID, conn)
	} else {
		log.Printf("Successfully pushed chat message to user %s.", recipientID)
	}
	return nil
}

type AckPayload struct {
	UserID    uuid.UUID `json:"user_id"`
	MessageID int64     `json:"message_id"`
	Timestamp int64     `json:"timestamp"`
}

func (g *Gateway) AckHandler(ctx context.Context, message *redis.Message) error {
	var ack AckPayload
	if err := json.Unmarshal([]byte(message.Payload), &ack); err != nil {
		return err
	}
	userID := ack.UserID

	g.mutex.RLock()
	conn, ok := g.connections[userID]
	g.mutex.RUnlock()

	if !ok {
		log.Printf("Ack: User %s not connected (offline).", userID)
		return nil
	}
	data := ClientMessage{
		Type:    AckEvent,
		Payload: json.RawMessage(message.Payload),
	}
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Ack push failed for user %s: %v", userID, err)
		g.deregisterConnection(ctx, userID, conn)
	} else {
		log.Printf("Successfully pushed ack to user %s.", userID)
	}
	return nil
}

type PresencePayload struct {
	UserID      uuid.UUID `json:"user_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Status      string    `json:"status"`
}

func (g *Gateway) PresenceHandler(ctx context.Context, message *redis.Message) error {
	var presence PresencePayload
	if err := json.Unmarshal([]byte(message.Payload), &presence); err != nil {
		return err
	}
	recipientID := presence.RecipientID

	g.mutex.RLock()
	conn, ok := g.connections[recipientID]
	g.mutex.RUnlock()

	if !ok {
		log.Printf("Presence: User %s not connected (offline).", recipientID)
		return nil
	}
	data := ClientMessage{
		Type:    PresenceEvent,
		Payload: json.RawMessage(message.Payload),
	}
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Presence push failed for user %s: %v", recipientID, err)
		g.deregisterConnection(ctx, recipientID, conn)
	} else {
		log.Printf("Successfully pushed presence to user %s.", recipientID)
	}
	return nil
}

type ReadPayload struct {
	UserID      uuid.UUID `json:"user_id"`
	RecipientID uuid.UUID `json:"recipient_id"`
	Timestamp   int64     `json:"timestamp"`
}

func (g *Gateway) ReadHandler(ctx context.Context, message *redis.Message) error {
	var read ReadPayload
	if err := json.Unmarshal([]byte(message.Payload), &read); err != nil {
		return err
	}
	recipientID := read.RecipientID

	g.mutex.RLock()
	conn, ok := g.connections[recipientID]
	g.mutex.RUnlock()

	if !ok {
		log.Printf("Read: User %s not connected (offline).", recipientID)
		return nil
	}
	data := ClientMessage{
		Type:    ReadEvent,
		Payload: json.RawMessage(message.Payload),
	}
	if err := conn.WriteJSON(data); err != nil {
		log.Printf("Read push failed for user %s: %v", recipientID, err)
		g.deregisterConnection(ctx, recipientID, conn)
	} else {
		log.Printf("Successfully pushed read to user %s.", recipientID)
	}
	return nil
}
