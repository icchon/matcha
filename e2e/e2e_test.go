package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	apiBaseURL = "http://localhost/api/v1"
	wsBaseURL  = "ws://localhost"
	wsPath     = "/ws"
)

// --- Structs for API interaction ---

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	IsVerified   bool      `json:"is_verified"`
	AuthMethod   string    `json:"auth_method"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

// --- Structs for WebSocket messages ---

// Generic message format for both client-sent and server-received
type WebsocketMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Payloads for client-sent messages
type ChatSendPayload struct {
	RecipientID uuid.UUID `json:"recipient_id"`
	Content     string    `json:"content"`
}
type ReadSendPayload struct {
	RecipientID uuid.UUID `json:"recipient_id"`
	Timestamp   int64     `json:"timestamp"`
}

// Payloads for server-received messages
type AckReceivePayload struct {
	MessageID int64 `json:"message_id"`
}
type ChatReceivePayload struct {
	ID       int64     `json:"id"`
	SenderID uuid.UUID `json:"sender_id"`
	Content  string    `json:"content"`
	SentAt   time.Time `json:"sent_at"`
}
type NotificationReceivePayload struct {
	SenderID string `json:"sender_id"`
	Type     string `json:"type"`
}
type PresenceReceivePayload struct {
	UserID uuid.UUID `json:"user_id"`
	Status string    `json:"status"`
}
type ReadReceivePayload struct {
	UserID    uuid.UUID `json:"user_id"`
	Timestamp int64     `json:"timestamp"`
}

// --- Test Subjects ---

type TestUser struct {
	ID    uuid.UUID
	Email string
	Pass  string
	Token string
	Conn  *websocket.Conn
	// Use a channel to receive messages concurrently and avoid race conditions
	Messages chan WebsocketMessage
	wg       sync.WaitGroup
	mu       sync.Mutex
}

// --- Helper Functions ---

func (u *TestUser) Cleanup() {
	u.mu.Lock()
	defer u.mu.Unlock()
	if u.Conn != nil {
		u.Conn.Close()
		u.wg.Wait() // Wait for reader goroutine to finish
	}
}

// Concurrently reads all messages from the websocket and puts them on a channel
func (u *TestUser) listen() {
	defer u.wg.Done()
	if u.Conn == nil {
		return
	}
	for {
		var msg WebsocketMessage
		err := u.Conn.ReadJSON(&msg)
		if err != nil {
			// This is expected when the connection is closed
			return
		}
		u.Messages <- msg
	}
}

func newUser(t *testing.T) *TestUser {
	user := &TestUser{
		Email: fmt.Sprintf("user-%s@example.com", uuid.NewString()),
		Pass:  "password123",
	}

	// Signup
	signupReqBody, _ := json.Marshal(SignupRequest{Email: user.Email, Password: user.Pass})
	signupResp, err := http.Post(apiBaseURL+"/auth/signup", "application/json", bytes.NewBuffer(signupReqBody))
	if err == nil {
		signupResp.Body.Close()
	}

	// Login
	loginReqBody, _ := json.Marshal(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{Email: user.Email, Password: user.Pass})
	loginResp, err := http.Post(apiBaseURL+"/auth/login", "application/json", bytes.NewBuffer(loginReqBody))
	require.NoError(t, err, "Login request failed for %s", user.Email)
	defer loginResp.Body.Close()
	require.Equal(t, http.StatusOK, loginResp.StatusCode, "Login failed with non-200 status for %s", user.Email)

	var loginData LoginResponse
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	require.NoError(t, err, "Failed to decode login response for %s", user.Email)

	user.ID = loginData.UserID
	user.Token = loginData.AccessToken
	return user
}

func (u *TestUser) connect(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer "+u.Token)

	wsURL, _ := url.Parse(wsBaseURL + wsPath)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL.String(), header)
	require.NoError(t, err, "Failed to connect user %s to WebSocket", u.Email)

	u.mu.Lock()
	u.Conn = conn
	u.Messages = make(chan WebsocketMessage, 10) // Buffered channel
	u.mu.Unlock()

	u.wg.Add(1)
	go u.listen()
}

func (u *TestUser) findMessage(t *testing.T, eventType string, timeout time.Duration) (WebsocketMessage, bool) {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		select {
		case msg := <-u.Messages:
			if msg.Type == eventType {
				return msg, true
			}
			log.Printf("User %s received other message type: %s (expected %s)", u.ID, msg.Type, eventType)
		case <-timer.C:
			t.Errorf("Timeout waiting for message type '%s' for user %s", eventType, u.ID)
			return WebsocketMessage{}, false
		}
	}
}

func likeUser(t *testing.T, liker *TestUser, liked *TestUser) {
	client := &http.Client{}
	url := fmt.Sprintf("%s/users/%s/like", apiBaseURL, liked.ID)
	req, err := http.NewRequest("POST", url, nil)
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+liker.Token)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

// Drains any stale messages from previous test runs.
func (u *TestUser) drainMessages() {
	for {
		select {
		case msg := <-u.Messages:
			log.Printf("Draining stale message for user %s: type %s", u.ID, msg.Type)
		case <-time.After(100 * time.Millisecond): // Timeout to prevent blocking forever
			return
		}
	}
}

func TestE2EWebSocketFlows(t *testing.T) {
	// --- Setup ---
	userA := newUser(t)
	userB := newUser(t)
	userC := newUser(t)

	defer userA.Cleanup()
	defer userB.Cleanup()
	defer userC.Cleanup()

	log.Printf("User A: %s", userA.ID)
	log.Printf("User B: %s", userB.ID)
	log.Printf("User C: %s", userC.ID)

	// --- Tests ---

	t.Run("Chat and Ack", func(t *testing.T) {
		userA.connect(t)
		userB.connect(t)
		userA.drainMessages()
		userB.drainMessages()

		// 1. User A sends a message to User B
		messageContent := "Hello User B from the test!"
		chatPayload, _ := json.Marshal(ChatSendPayload{
			RecipientID: userB.ID,
			Content:     messageContent,
		})
		msgToSend := WebsocketMessage{Type: "chat_event", Payload: chatPayload}
		require.NoError(t, userA.Conn.WriteJSON(msgToSend))

		// 2. User A should receive an ACK for their sent message
		ackMsg, found := userA.findMessage(t, "ack_event", 5*time.Second)
		require.True(t, found, "User A did not receive ACK")
		var ackPayload AckReceivePayload
		require.NoError(t, json.Unmarshal(ackMsg.Payload, &ackPayload))
		assert.NotZero(t, ackPayload.MessageID, "ACK payload has zero message ID")
		log.Printf("SUCCESS: User A received ACK for message")

		// 3. User B should receive the chat message
		chatMsg, found := userB.findMessage(t, "chat_event", 5*time.Second)
		require.True(t, found, "User B did not receive chat message")
		var receivedChatPayload ChatReceivePayload
		require.NoError(t, json.Unmarshal(chatMsg.Payload, &receivedChatPayload))

		assert.Equal(t, messageContent, receivedChatPayload.Content)
		assert.Equal(t, userA.ID, receivedChatPayload.SenderID)
		log.Printf("SUCCESS: User B received chat message from User A")
	})

	t.Run("Notifications", func(t *testing.T) {
		userC.connect(t)
		userC.drainMessages()

		// 1. User A (via HTTP) "likes" User C
		likeUser(t, userA, userC)

		// 2. User C should receive a 'like' notification via WebSocket
		notifMsg, found := userC.findMessage(t, "notification_event", 5*time.Second)
		require.True(t, found, "User C did not receive notification")

		var notifPayload NotificationReceivePayload
		require.NoError(t, json.Unmarshal(notifMsg.Payload, &notifPayload))
		assert.Equal(t, "like", notifPayload.Type)
		assert.Equal(t, userA.ID.String(), notifPayload.SenderID)
		log.Printf("SUCCESS: User C received 'like' notification from User A")
	})

	t.Run("Read Receipts", func(t *testing.T) {
		// Uses UserA and UserB from previous test's connection
		if userA.Conn == nil {
			userA.connect(t)
		}
		if userB.Conn == nil {
			userB.connect(t)
		}
		userA.drainMessages()
		userB.drainMessages()

		// 1. User B sends a `read` event for their chat with User A
		readPayload, _ := json.Marshal(ReadSendPayload{
			RecipientID: userA.ID,
			Timestamp:   time.Now().UnixMilli(),
		})
		msgToSend := WebsocketMessage{Type: "read_event", Payload: readPayload}
		require.NoError(t, userB.Conn.WriteJSON(msgToSend))

		// 2. User A should receive the read receipt
		readMsg, found := userA.findMessage(t, "read_event", 5*time.Second)
		require.True(t, found, "User A did not receive read receipt")
		var readReceivePayload ReadReceivePayload
		require.NoError(t, json.Unmarshal(readMsg.Payload, &readReceivePayload))
		assert.Equal(t, userB.ID, readReceivePayload.UserID)
		log.Printf("SUCCESS: User A received read receipt from User B")
	})

	t.Run("Presence", func(t *testing.T) {
		// Ensure user A is connected and clean
		if userA.Conn == nil {
			userA.connect(t)
		}
		userA.drainMessages()

		// 1. Create a match between A and B so they receive each other's presence
		likeUser(t, userA, userB)
		likeUser(t, userB, userA)
		log.Printf("Users A and B are now matched.")

		// 2. Disconnect User B to generate "offline" event
		userB.Cleanup()
		log.Printf("User B disconnected. Draining messages for User A to clear offline event.")
		// Drain User A's channel to remove the 'offline' message before checking for 'online'
		userA.drainMessages()

		// 3. Reconnect User B to generate "online" event
		time.Sleep(1 * time.Second) // Give time for disconnect to be fully processed on the backend
		userB.connect(t)

		// 4. User A should now receive the "online" presence update for User B
		presenceMsg, found := userA.findMessage(t, "presence_event", 5*time.Second)
		require.True(t, found, "User A did not receive presence event for User B")

		var presencePayload PresenceReceivePayload
		require.NoError(t, json.Unmarshal(presenceMsg.Payload, &presencePayload))

		assert.Equal(t, userB.ID, presencePayload.UserID)
		assert.Equal(t, "online", presencePayload.Status)
		log.Printf("SUCCESS: User A received 'online' presence from User B")
	})
}
