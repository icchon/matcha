# Matcha API Specification

This document outlines the API endpoints for the Matcha application.

---

## Authentication

### Signup

-   **URL:** `/api/v1/auth/signup`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "email": "user@example.com",
        "password": "password123"
    }
    ```
-   **Response:**
    ```json
    {
        "message": "User created successfully. Please verify your email."
    }
    ```

### Login

-   **URL:** `/api/v1/auth/login`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "email": "user@example.com",
        "password": "password123"
    }
    ```
-   **Response:**
    ```json
    {
        "access_token": "...",
        "refresh_token": "..."
    }
    ```

### Logout

-   **URL:** `/api/v1/auth/logout`
-   **Method:** `POST`
-   **Request:** (No body, requires Authorization header)
-   **Response:**
    ```json
    {
        "message": "Logged out successfully"
    }
    ```

### Email Verification

-   **URL:** `/api/v1/auth/verify/{token}`
-   **Method:** `GET`
-   **Request:** URL parameter `token`.
-   **Response:**
    ```json
    {
        "message": "Email verified successfully"
    }
    ```
    
### Resend Verification Email

-   **URL:** `/api/v1/auth/verify/mail`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "email": "user@example.com"
    }
    ```
-   **Response:**
    ```json
    {
        "message": "Verification email sent"
    }
    ```

### Forgot Password

-   **URL:** `/api/v1/auth/password/forgot`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "email": "user@example.com"
    }
    ```
-   **Response:**
    ```json
    {
        "message": "Password reset email sent"
    }
    ```

### Reset Password

-   **URL:** `/api/v1/auth/password/reset`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "token": "...",
        "password": "new_password"
    }
    ```
-   **Response:**
    ```json
    {
        "message": "Password reset successfully"
    }
    ```

### OAuth - Google Login

-   **URL:** `/api/v1/auth/oauth/google/login`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "code": "oauth_code_from_google"
    }
    ```
-   **Response:**
    ```json
    {
        "access_token": "...",
        "refresh_token": "..."
    }
    ```

### OAuth - Github Login

-   **URL:** `/api/v1/auth/oauth/github/login`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "code": "oauth_code_from_github"
    }
    ```
-   **Response:**
    ```json
    {
        "access_token": "...",
        "refresh_token": "..."
    }
    ```
    
---

## User Actions

### Like a User

-   **URL:** `/api/v1/users/{userID}/like`
-   **Method:** `POST`
-   **Request:** URL parameter `userID`. Requires Authorization header.
-   **Response:**
    ```json
    {
        "connection": { /* connection object if a match is made */ },
        "message": "User liked successfully" or "It's a match!"
    }
    ```

### Unlike a User

-   **URL:** `/api/v1/users/{userID}/like`
-   **Method:** `DELETE`
-   **Request:** URL parameter `userID`. Requires Authorization header.
-   **Response:**
    ```json
    {
        "message": "User unliked successfully"
    }
    ```

### Block a User

-   **URL:** `/api/v1/users/{userID}/block`
-   **Method:** `POST`
-   **Request:** URL parameter `userID`. Requires Authorization header.
-   **Response:**
    ```json
    {
        "message": "User blocked successfully"
    }
    ```

---

## Current User (`/me`)

### Delete My Account

-   **URL:** `/api/v1/me/`
-   **Method:** `DELETE`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "message": "User account deleted successfully"
    }
    ```

### Get My Liked List

-   **URL:** `/api/v1/me/likes`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "likes": [ /* array of like objects */ ]
    }
    ```
    
### Get My Viewed List

-   **URL:** `/api/v1/me/views`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "views": [ /* array of view objects */ ]
    }
    ```

### Get My Blocked List

-   **URL:** `/api/v1/me/blocks`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "blocks": [ /* array of block objects */ ]
    }
    ```

### Get My Chats

-   **URL:** `/api/v1/me/chats`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of chat objects */ ]
    ```

### Get My Notifications

-   **URL:** `/api/v1/me/notifications`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of notification objects */ ]
    ```

### My User Data

-   **URL:** `/api/v1/me/data`
-   **Method:** `GET`, `POST`, `PUT`
-   **Request:** Requires Authorization header.
    -   `POST`/`PUT` Body:
        ```json
        {
            "latitude": 35.68,
            "longitude": 139.76,
            "internal_score": 100
        }
        ```
-   **Response:**
    ```json
    { /* user_data object */ }
    ```

### My User Tags

-   **URL:** `/api/v1/me/tags`
-   **Method:** `GET`, `POST`
-   **Request:**
    -   `GET`: Requires Authorization header.
    -   `POST` Body:
        ```json
        {
            "tag_id": 1
        }
        ```
-   **Response:**
    -   `GET`: `[ /* array of tag objects */ ]`
    -   `POST`: `{ "message": "Tag added successfully" }`

### Delete My User Tag

-   **URL:** `/api/v1/me/tags/{tagID}`
-   **Method:** `DELETE`
-   **Request:** URL parameter `tagID`. Requires Authorization header.
-   **Response:** `204 No Content`

### My Profile

-   **URL:** `/api/v1/me/profile`
-   **Method:** `POST`, `PUT`
-   **Request Body:**
    ```json
    {
        "first_name": "Test",
        "last_name": "User",
        "username": "testuser",
        "gender": "male",
        "sexual_preference": "bisexual",
        "birthday": "1990-01-01T00:00:00Z",
        "occupation": "Developer",
        "biography": "...",
        "location_name": "Tokyo"
    }
    ```
-   **Response:**
    ```json
    { /* user_profile object */ }
    ```
    
### Get Who Liked Me

-   **URL:** `/api/v1/me/profile/likes`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "likes": [ /* array of like objects */ ]
    }
    ```

### Get Who Viewed Me

-   **URL:** `/api/v1/me/profile/views`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "views": [ /* array of view objects */ ]
    }
    ```

---

## Profiles

### Get User Profiles (Filtered)

-   **URL:** `/api/v1/profiles`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
    -   Query Params: `age_min`, `age_max`, `gender`
-   **Response:**
    ```json
    [ /* array of user_profile objects */ ]
    ```
    
### Get Recommended Profiles

-   **URL:** `/api/v1/profiles/recommends`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of user_profile objects, sorted by recommendation score */ ]
    ```

### Get a Specific User's Profile

-   **URL:** `/api/v1/users/{userID}/profile`
-   **Method:** `GET`
-   **Request:** URL parameter `userID`. Requires Authorization header.
-   **Response:**
    ```json
    { /* user_profile object */ }
    ```

---

## Tags

### Get All Tags

-   **URL:** `/api/v1/tags`
-   **Method:** `GET`
-   **Request:** (No parameters)
-   **Response:**
    ```json
    [ /* array of tag objects */ ]
    ```

---

## Chat

### Get Chat Messages

-   **URL:** `/api/v1/chats/{userID}/messages`
-   **Method:** `GET`
-   **Request:**
    -   URL parameter `userID`.
    -   Query Params: `limit`, `offset`.
    -   Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of message objects */ ]
    ```

---

## WebSockets

### Real-time Communication

-   **URL:** `/ws`
-   **Method:** `GET` (WebSocket upgrade request)
-   **Authentication:** Requires JWT token in query parameter (e.g., `/ws?token=...`) or `Authorization` header during handshake.
-   **Description:** This endpoint is used for real-time communication for features like chat, presence (online/offline status), and notifications.
-   **Messages (examples):**
    -   **Client -> Server (Send Chat Message):**
        ```json
        {
            "type": "chat_message",
            "recipient_id": "uuid_of_recipient",
            "content": "Hello there!"
        }
        ```
    -   **Server -> Client (New Chat Message):**
        ```json
        {
            "type": "chat_message",
            "sender_id": "uuid_of_sender",
            "content": "Hello there!",
            "sent_at": "timestamp",
            "id": "message_id"
        }
        ```
    -   **Server -> Client (Notification):**
        ```json
        {
            "type": "notification",
            "sender_id": "uuid_of_sender",
            "notification_type": "like",
            "message": "User X liked your profile"
        }
        ```
    -   **Server -> Client (Presence Update):**
        ```json
        {
            "type": "presence",
            "user_id": "uuid_of_user",
            "status": "online"
        }
        ```