package main

import (
	"database/sql" // Required for sql.NullString etc.

	"fmt"
	"log"
	"math/rand" // Required for random numbers
	"os"
	"time" // Added for time.Now()

	"github.com/go-faker/faker/v4"                        // Added for faker data
	"github.com/google/uuid"                              // Added for uuid.New()
	"github.com/icchon/matcha/api/internal/domain/entity" // Added for entity.User, entity.Tag
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"io/ioutil" // Added for reading HTTP response body
	"net/http"  // Added for making HTTP requests
)

const (
	numUsers = 100
	numTags  = 10
)

// import (
// 	"database/sql" // Required for sql.NullString etc.
// 	"fmt"
// 	"io/ioutil" // Added for reading HTTP response body
// 	"log"
// 	"math/rand" // Required for random numbers
// 	"net/http" // Added for making HTTP requests
// 	"os"
// 	"time" // Added for time.Now()

// 	"github.com/go-faker/faker/v4" // Added for faker data
// 	"github.com/google/uuid" // Added for uuid.New()
// 	"github.com/icchon/matcha/api/internal/domain/entity" // Added for entity.User, entity.Tag
// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

func getDogImageUrl() string {
	resp, err := http.Get("http://filesrv/dog")
	if err != nil {
		log.Printf("Failed to fetch dog image URL: %v. Using placeholder.", err)
		return "https://example.com/placeholder-dog.jpg" // Fallback
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read dog image URL response body: %v. Using placeholder.", err)
		return "https://example.com/placeholder-dog.jpg" // Fallback
	}
	return string(body)
}

func createViews(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	viewCount := 0
	numViewsToCreate := numUsers * 10 // Create about ten times as many views as users

	for i := 0; i < numViewsToCreate; i++ {
		viewer := users[rand.Intn(len(users))]
		viewed := users[rand.Intn(len(users))]

		if viewer.ID == viewed.ID {
			continue // A user cannot view themselves
		}

		view := entity.View{
			ViewerID: viewer.ID,
			ViewedID: viewed.ID,
			ViewTime: time.Now().Add(time.Duration(-rand.Intn(720)) * time.Hour), // Views within the last month
		}

		// Check if this view already exists within a short timeframe to avoid spamming
		// This is a simplification; a real app might have more complex logic
		var existingViewerID string
		err := db.QueryRow("SELECT viewer_id FROM views WHERE viewer_id = $1 AND viewed_id = $2 AND view_time > $3", viewer.ID, viewed.ID, time.Now().Add(-1*time.Hour)).Scan(&existingViewerID)
		if err == nil { // View already exists recently
			continue
		} else if err != sql.ErrNoRows { // Some other error
			log.Fatalf("Failed to check for existing view: %v", err)
		}

		tx.NamedExec(`INSERT INTO views (viewer_id, viewed_id, view_time) VALUES (:viewer_id, :viewed_id, :view_time)`, view)
		viewCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert views: %v", err)
	}
	log.Printf("Created %d view records\n", viewCount)
}

func createVerificationTokens(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	vtCount := 0
	numVTsToCreate := numUsers / 5 // Create VTs for about 20% of users

	for i := 0; i < numVTsToCreate; i++ {
		user := users[rand.Intn(len(users))]
		verificationToken := entity.VerificationToken{
			Token:     faker.UUIDHyphenated(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24), // Expires in 24 hours
		}

		tx.NamedExec(`INSERT INTO verification_tokens (token, user_id, expires_at)
            VALUES (:token, :user_id, :expires_at)`, verificationToken)
		vtCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert verification tokens: %v", err)
	}
	log.Printf("Created %d verification token records\n", vtCount)
}

func createRefreshTokens(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	rtCount := 0
	numRTsToCreate := numUsers / 5 // Create RTs for about 20% of users

	for i := 0; i < numRTsToCreate; i++ {
		user := users[rand.Intn(len(users))]
		refreshToken := entity.RefreshToken{
			TokenHash: faker.UUIDHyphenated(), // This would be a hash in a real app
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 30), // Expires in 30 days
			CreatedAt: time.Now().Add(time.Duration(-rand.Intn(720)) * time.Hour),
			Revoked:   rand.Float32() < 0.1, // 10% chance of being revoked
		}

		tx.NamedExec(`INSERT INTO refresh_tokens (token_hash, user_id, expires_at, created_at, revoked)
            VALUES (:token_hash, :user_id, :expires_at, :created_at, :revoked)`, refreshToken)
		rtCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert refresh tokens: %v", err)
	}
	log.Printf("Created %d refresh token records\n", rtCount)
}

func createPasswordResets(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	prCount := 0
	numPRsToCreate := numUsers / 10 // Create PRs for about 10% of users

	for i := 0; i < numPRsToCreate; i++ {
		user := users[rand.Intn(len(users))]
		passwordReset := entity.PasswordReset{
			UserID:    user.ID,
			Token:     faker.UUIDHyphenated(),
			ExpiresAt: time.Now().Add(time.Hour * 24), // Expires in 24 hours
		}

		if _, err := tx.NamedExec(`INSERT INTO password_resets (user_id, token, expires_at)
            VALUES (:user_id, :token, :expires_at) ON CONFLICT (user_id) DO NOTHING`, passwordReset); err != nil {
			log.Fatalf("Error inserting password reset for user %s: %v", user.ID, err)
		}
		prCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert password resets: %v", err)
	}
	log.Printf("Created %d password reset records\n", prCount)
}

func createNotifications(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	notificationCount := 0
	numNotificationsToCreate := numUsers * 5

	notificationTypes := []entity.NotificationType{
		entity.NotifLike,
		entity.NotifMatch,
		entity.NotifMessage,
		entity.NotifView,
		entity.NotifUnlike,
	}

	for i := 0; i < numNotificationsToCreate; i++ {
		recipient := users[rand.Intn(len(users))]
		sender := users[rand.Intn(len(users))]

		senderID := sql.NullString{Valid: false}
		if rand.Float32() < 0.8 { // 80% chance to have a sender
			senderID = sql.NullString{String: sender.ID.String(), Valid: true}
		}

		notification := entity.Notification{
			RecipientID: recipient.ID,
			SenderID:    senderID,
			Type:        notificationTypes[rand.Intn(len(notificationTypes))],
			IsRead:      sql.NullBool{Bool: rand.Float32() < 0.6, Valid: true}, // 60% chance of being read
			CreatedAt:   time.Now().Add(time.Duration(-rand.Intn(720)) * time.Hour),
		}

		tx.NamedExec(`INSERT INTO notifications (recipient_id, sender_id, type, is_read, created_at)
            VALUES (:recipient_id, :sender_id, :type, :is_read, :created_at)`, notification)
		notificationCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert notifications: %v", err)
	}
	log.Printf("Created %d notification records\n", notificationCount)
}

func createMessages(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	messageCount := 0
	numMessagesToCreate := numUsers * 5 // Create about five times as many messages as users

	// Fetch existing connections to ensure messages are exchanged between connected users
	var connections []entity.Connection
	err := db.Select(&connections, "SELECT user1_id, user2_id FROM connections")
	if err != nil {
		log.Fatalf("Failed to fetch existing connections: %v", err)
	}

	if len(connections) == 0 {
		log.Println("No connections found, skipping message seeding.")
		return
	}

	for i := 0; i < numMessagesToCreate; i++ {
		// Pick a random connection
		conn := connections[rand.Intn(len(connections))]

		senderID := conn.User1ID
		recipientID := conn.User2ID

		// Randomly swap sender and recipient to make it more realistic
		if rand.Float32() < 0.5 {
			senderID, recipientID = recipientID, senderID
		}

		message := entity.Message{
			SenderID:    senderID,
			RecipientID: recipientID,
			Content:     faker.Sentence(),
			SentAt:      time.Now().Add(time.Duration(-rand.Intn(720)) * time.Hour), // Messages sent within the last month
			IsRead:      sql.NullBool{Bool: rand.Float32() < 0.7, Valid: true},      // 70% chance of being read
		}

		tx.NamedExec(`INSERT INTO messages (sender_id, recipient_id, content, sent_at, is_read)
            VALUES (:sender_id, :recipient_id, :content, :sent_at, :is_read)`, message)
		messageCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert messages: %v", err)
	}
	log.Printf("Created %d message records\n", messageCount)
}

func createLikes(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	likeCount := 0
	numLikesToCreate := numUsers * 4 // Create about four times as many likes as users

	for i := 0; i < numLikesToCreate; i++ {
		liker := users[rand.Intn(len(users))]
		liked := users[rand.Intn(len(users))]

		if liker.ID == liked.ID {
			continue // A user cannot like themselves
		}

		like := entity.Like{
			LikerID:   liker.ID,
			LikedID:   liked.ID,
			CreatedAt: time.Now(),
		}

		// Check if this like already exists to avoid duplicates
		// Removed: ON CONFLICT DO NOTHING handles this


		if _, err := tx.NamedExec(`INSERT INTO likes (liker_id, liked_id, created_at) VALUES (:liker_id, :liked_id, :created_at) ON CONFLICT (liker_id, liked_id) DO NOTHING`, like); err != nil {
			log.Fatalf("Error inserting like for liker %s and liked %s: %v", liker.ID, liked.ID, err)
		}
		likeCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert likes: %v", err)
	}
	log.Printf("Created %d like records\n", likeCount)
}

func createConnections(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	defer tx.Rollback() // Ensure transaction is rolled back on error or early exit
	connectionCount := 0
	numConnectionsToCreate := numUsers * 3 // Create about three times as many connections as users

	for i := 0; i < numConnectionsToCreate; i++ {
		user1 := users[rand.Intn(len(users))]
		user2 := users[rand.Intn(len(users))]

		if user1.ID == user2.ID {
			continue // A user cannot connect to themselves
		}

		// Ensure consistent ordering for connection (user1_id always less than user2_id)
		// This prevents duplicate connections like (A,B) and (B,A)
		var u1ID, u2ID uuid.UUID
		if user1.ID.String() < user2.ID.String() {
			u1ID = user1.ID
			u2ID = user2.ID
		} else {
			u1ID = user2.ID
			u2ID = user1.ID
		}

		connection := entity.Connection{
			User1ID:   u1ID,
			User2ID:   u2ID,
			CreatedAt: time.Now(),
		}

		if _, err := tx.NamedExec(`INSERT INTO connections (user1_id, user2_id, created_at) VALUES (:user1_id, :user2_id, :created_at) ON CONFLICT (user1_id, user2_id) DO NOTHING`, connection); err != nil {
			// If there's an error, log it and terminate.
			// The ON CONFLICT clause handles duplicates, so other errors are critical.
			log.Fatalf("Error inserting connection for users %s and %s: %v", u1ID, u2ID, err)
		} else {
			connectionCount++
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert connections: %v", err)
	}
	log.Printf("Created %d connection records\n", connectionCount)
}

func createBlocks(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	defer tx.Rollback() // Ensure transaction is rolled back on error or early exit
	blockCount := 0
	numBlocksToCreate := numUsers * 2 // Create about twice as many blocks as users

	for i := 0; i < numBlocksToCreate; i++ {
		blocker := users[rand.Intn(len(users))]
		blocked := users[rand.Intn(len(users))]

		if blocker.ID == blocked.ID {
			continue // A user cannot block themselves
		}

		block := entity.Block{
			BlockerID: blocker.ID,
			BlockedID: blocked.ID,
		}

		if _, err := tx.NamedExec(`INSERT INTO blocks (blocker_id, blocked_id) VALUES (:blocker_id, :blocked_id) ON CONFLICT (blocker_id, blocked_id) DO NOTHING`, block); err != nil {
			// if pqErr, ok := err.(*pq.Error); ok && pqErr.SQLState() == "23505" {
			// 	log.Printf("Ignored duplicate block for users %s and %s", blocker.ID, blocked.ID)
			// } else {
			// 	log.Fatalf("Error inserting block for users %s and %s: %v", blocker.ID, blocked.ID, err)
			// }
		} else {
			blockCount++
		}
		blockCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert blocks: %v", err)
	}
	log.Printf("Created %d block records\n", blockCount)
}

func createPictures(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	pictureCount := 0
	for _, user := range users {
		numPictures := rand.Intn(3) + 1 // 1 to 3 pictures per user
		for i := 0; i < numPictures; i++ {
			isProfilePic := sql.NullBool{Bool: false, Valid: true}
			if i == 0 { // First picture can be a profile pic
				isProfilePic.Bool = true
			}
			picture := entity.Picture{
				UserID:       user.ID,
				URL:          getDogImageUrl(),
				IsProfilePic: isProfilePic,
				CreatedAt:    time.Now(),
			}
			tx.NamedExec(`INSERT INTO pictures (user_id, url, is_profile_pic, created_at)
                VALUES (:user_id, :url, :is_profile_pic, :created_at)`, picture)
			pictureCount++
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert pictures: %v", err)
	}
	log.Printf("Created %d picture records\n", pictureCount)
}

func createUserTags(db *sqlx.DB, users []entity.User, tags []entity.Tag) {
	tx := db.MustBegin()
	defer tx.Rollback() // Ensure transaction is rolled back on error or early exit
	userTagCount := 0
	for _, user := range users { // Outer loop: for each user
		numTagsToAssign := rand.Intn(5) + 1 // Re-introduce this line
		// Create a shuffled list of tag IDs to ensure unique assignment
		shuffledTags := make([]entity.Tag, len(tags))
		copy(shuffledTags, tags)
		rand.Shuffle(len(shuffledTags), func(i, j int) {
			shuffledTags[i], shuffledTags[j] = shuffledTags[j], shuffledTags[i]
		})

		for i := 0; i < numTagsToAssign; i++ {
			if i >= len(shuffledTags) { // Ensure we don't try to assign more unique tags than available
				break
			}
			tag := shuffledTags[i]

			userTag := entity.UserTag{
				UserID: user.ID,
				TagID:  tag.ID,
			}
			if _, err := tx.NamedExec(`INSERT INTO user_tags (user_id, tag_id) VALUES (:user_id, :tag_id)`, userTag); err != nil {
				log.Printf("Error inserting user tag for user %s, tag %d: %v", user.ID, tag.ID, err)
				// Do not tx.Rollback() here, let the defer do it on function exit or let tx.Commit() fail
				continue // Skip to the next tag assignment
			}
			userTagCount++
		}
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert user tags: %v", err)
	}
	log.Printf("Created %d user tag records\n", userTagCount)
}

func createUserProfiles(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	for _, user := range users {
		birthday := time.Now().AddDate(-rand.Intn(40)-18, rand.Intn(12)-6, rand.Intn(28)-14) // 18-58 years old
		genders := []string{"male", "female", "other"}
		sexualPreferences := []string{"heterosexual", "homosexual", "bisexual"}

		userProfile := entity.UserProfile{
			UserID:           user.ID,
			FirstName:        sql.NullString{String: faker.FirstName(), Valid: true},
			LastName:         sql.NullString{String: faker.LastName(), Valid: true},
			Username:         sql.NullString{String: faker.Username(), Valid: true},
			Gender:           sql.NullString{String: genders[rand.Intn(len(genders))], Valid: true},
			SexualPreference: sql.NullString{String: sexualPreferences[rand.Intn(len(sexualPreferences))], Valid: true},
			Birthday:         sql.NullTime{Time: birthday, Valid: true},
			Occupation:       sql.NullString{String: faker.Word(), Valid: true},
			Biography:        sql.NullString{String: faker.Sentence(), Valid: true},
			FameRating:       sql.NullInt32{Int32: rand.Int31n(100) + 1, Valid: true}, // 1 to 100
			LocationName:     sql.NullString{String: faker.Word(), Valid: true},
		}
		tx.NamedExec(`INSERT INTO user_profiles (user_id, first_name, last_name, username, gender, sexual_preference, birthday, occupation, biography, fame_rating, location_name)
            VALUES (:user_id, :first_name, :last_name, :username, :gender, :sexual_preference, :birthday, :occupation, :biography, :fame_rating, :location_name)`, userProfile)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert user profiles: %v", err)
	}
	log.Printf("Created %d user profiles\n", len(users))
}

func createUserData(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	for _, user := range users {
		userData := entity.UserData{
			UserID:        user.ID,
			Latitude:      sql.NullFloat64{Float64: rand.Float64()*180 - 90, Valid: true},  // -90 to +90
			Longitude:     sql.NullFloat64{Float64: rand.Float64()*360 - 180, Valid: true}, // -180 to +180
			InternalScore: sql.NullInt32{Int32: rand.Int31n(1000) + 1, Valid: true},        // 1 to 1000
		}
		tx.NamedExec(`INSERT INTO user_data (user_id, latitude, longitude, internal_score)
            VALUES (:user_id, :latitude, :longitude, :internal_score)`, userData)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert user data: %v", err)
	}
	log.Printf("Created %d user data records\n", len(users))
}

func createTags(db *sqlx.DB) []entity.Tag {
	var tags []entity.Tag
	err := db.Select(&tags, "SELECT id, name FROM tags ORDER BY id")
	if err != nil {
		log.Fatalf("Failed to query existing tags: %v", err)
	}
	log.Printf("Fetched %d existing tags: %v\n", len(tags), tags)
	return tags
}

func main() {
	dbConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))

	db, err := sqlx.Connect("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Seeding database...")

	users := createUsers(db)            // Call createUsers
	tags := createTags(db)              // Call createTags
	createAuths(db, users)              // Call createAuths
	createUserData(db, users)           // Call createUserData
	createUserProfiles(db, users)       // Call createUserProfiles
	createUserTags(db, users, tags)     // Call createUserTags
	
	createBlocks(db, users)             // Call createBlocks
	createConnections(db, users)        // Call createConnections
	createLikes(db, users)              // Call createLikes
	createMessages(db, users)           // Call createMessages
	createNotifications(db, users)      // Call createNotifications
	createPasswordResets(db, users)     // Call createPasswordResets
	createRefreshTokens(db, users)      // Call createRefreshTokens
	createVerificationTokens(db, users) // Call createVerificationTokens
	createViews(db, users)              // Call createViews
	createPictures(db, users)           // Call createPictures

	log.Println("Seeding completed successfully!")
}

func createUsers(db *sqlx.DB) []entity.User {
	users := make([]entity.User, 0, numUsers)
	for i := 0; i < numUsers; i++ {
		user := entity.User{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
		}
		users = append(users, user)
	}

	tx := db.MustBegin()
	for _, user := range users {
		tx.MustExec("INSERT INTO users (id, created_at) VALUES ($1, $2)", user.ID, user.CreatedAt)
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert users: %v", err)
	}
	log.Printf("Created %d users\n", len(users))
	return users
}

func createAuths(db *sqlx.DB, users []entity.User) {
	tx := db.MustBegin()
	authCount := 0
	for _, user := range users {
		auth := entity.Auth{
			UserID:     user.ID,
			IsVerified: true,
		}

		// Randomly choose between email/password and OAuth provider
		if rand.Float32() < 0.7 { // 70% chance for email/password
			auth.Email = sql.NullString{String: faker.Email(), Valid: true}
			auth.PasswordHash = sql.NullString{String: hashPassword(faker.Password()), Valid: true} // Hash a fake password
			auth.Provider = entity.ProviderLocal
		} else { // 30% chance for OAuth provider
			provider := randomAuthProvider()
			auth.Provider = provider
			auth.ProviderUID = sql.NullString{String: faker.UUIDHyphenated(), Valid: true} // A fake provider UID
			// For OAuth, email might also be present
			if rand.Float32() < 0.5 {
				auth.Email = sql.NullString{String: faker.Email(), Valid: true}
			}
		}
		tx.NamedExec(`INSERT INTO auths (user_id, email, provider, provider_uid, is_verified, password_hash)
            VALUES (:user_id, :email, :provider, :provider_uid, :is_verified, :password_hash)`, auth)
		authCount++
	}
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to insert auths: %v", err)
	}
	log.Printf("Created %d auth records\n", authCount)
}

func randomAuthProvider() entity.AuthProvider {
	providers := []entity.AuthProvider{entity.ProviderGoogle, entity.ProviderGithub}
	return providers[rand.Intn(len(providers))]
}

// hashPassword is a dummy function for seeding, in a real app this would use a proper hashing algorithm
func hashPassword(password string) string {
	return fmt.Sprintf("hashed_%s", password)
}
