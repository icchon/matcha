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
-   **Response:**
    ```json
    {
        "user_id": "string (UUID)",
        "is_verified": "boolean",
        "auth_method": "string (e.g., local, google, github)",
        "access_token": "string",
        "refresh_token": "string"
    }
    ```

### Logout

-   **URL:** `/api/v1/auth/logout`
-   **Method:** `POST`
-   **Request:** (No body, requires Authorization header)
-   **Response:**
    ```json
    {}
    ```

### Email Verification

-   **URL:** `/api/v1/auth/verify/{token}`
-   **Method:** `GET`
-   **Request:** URL parameter `token`.
-   **Response:**
    ```json
    {}
    ```
    
### Resend Verification Email

-   **URL:** `/api/v1/auth/verify/mail`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "user_id": "string (UUID)",
        "email": "user@example.com"
    }
    ```
-   **Response:**
    ```json
    {
        "message": "Please check your email to verify your account"
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
        "message": "Please check your email to reset your password"
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
    {}
    ```

### OAuth - Google Login

-   **URL:** `/api/v1/auth/oauth/google/login`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "code": "oauth_code_from_google",
        "code_verifier": "string"
    }
    ```
-   **Response:**
    ```json
    {
        "user_id": "string (UUID)",
        "is_verified": "boolean",
        "auth_method": "string (e.g., local, google, github)",
        "access_token": "string",
        "refresh_token": "string"
    }
    ```

### OAuth - Github Login

-   **URL:** `/api/v1/auth/oauth/github/login`
-   **Method:** `POST`
-   **Request Body:**
    ```json
    {
        "code": "oauth_code_from_github",
        "code_verifier": "string"
    }
    ```
-   **Response:**
    ```json
    {
        "user_id": "string (UUID)",
        "is_verified": "boolean",
        "auth_method": "string (e.g., local, google, github)",
        "access_token": "string",
        "refresh_token": "string"
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
        "connection": "Connection Object or null (if no match is made)",
        "message": "string (e.g., User liked successfully or It's a match!)"
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
        "likes": [ /* array of Like Objects */ ]
    }
    ```
    
### Get My Viewed List

-   **URL:** `/api/v1/me/views`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "views": [ /* array of View Objects */ ]
    }
    ```

### Get My Blocked List

-   **URL:** `/api/v1/me/blocks`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    {
        "blocks": [ /* array of Block Objects */ ]
    }
    ```

### Get My Chats

-   **URL:** `/api/v1/me/chats`
-   **Method:** `GET`
-   **Response:**
    ```json
    [
        {
            "other_user": "UserProfile Object",
            "last_message": "Message Object or null"
        }
    ]
    ```

### Get My Notifications

-   **URL:** `/api/v1/me/notifications`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of Notification Objects */ ]
    ```

### My User Data

-   **URL:** `/api/v1/me/data`
-   **Method:** `GET`, `POST`, `PUT`
-   **Request:** Requires Authorization header.
    -   `POST`/`PUT` Body: `UserData Object` (excluding `user_id` which is taken from the authenticated user)
        ```json
        {
            "latitude": "number or null (float)",
            "longitude": "number or null (float)",
            "internal_score": "integer or null"
        }
        ```
-   **Response:** `UserData Object`
    ```json
    {
        "user_id": "string (UUID)",
        "latitude": "number or null (float)",
        "longitude": "number or null (float)",
        "internal_score": "integer or null"
    }
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
    -   `GET`: `[ /* array of Tag Objects */ ]`
    -   `POST`: `{ "message": "Tag added successfully" }`

### Delete My User Tag

-   **URL:** `/api/v1/me/tags/{tagID}`
-   **Method:** `DELETE`
-   **Request:** URL parameter `tagID`. Requires Authorization header.
-   **Response:** `{}`

### My Profile

-   **URL:** `/api/v1/me/profile`
-   **Method:** `POST`, `PUT`
-   **Request Body:**
    ```json
    {
        "first_name": "string or null",
        "last_name": "string or null",
        "username": "string or null",
        "gender": "string or null (male, female, other)",
        "sexual_preference": "string or null (heterosexual, homosexual, bisexual)",
        "birthday": "string or null (timestamp)",
        "occupation": "string or null",
        "biography": "string or null",
        "location_name": "string or null"
    }
    ```
-   **Response:** `UserProfile Object`
    ```json
    {
        "user_id": "string (UUID)",
        "first_name": "string or null",
        "last_name": "string or null",
        "username": "string or null",
        "gender": "string or null (e.g., male, female, other)",
        "sexual_preference": "string or null (e.g., heterosexual, homosexual, bisexual)",
        "birthday": "string or null (timestamp)",
        "occupation": "string or null",
        "biography": "string or null",
        "fame_rating": "integer or null",
        "location_name": "string or null",
        "distance": "number or null (float)"
    }
    ```

### Upload Profile Picture

-   **URL:** `/api/v1/me/profile/pictures`
-   **Method:** `POST`
-   **Request:** `multipart/form-data` with a file field named `image`. Requires Authorization header.
-   **Response:** `Picture Object`
    ```json
    {
        "id": "integer",
        "user_id": "string (UUID)",
        "url": "string (URL)",
        "is_profile_pic": "boolean or null",
        "created_at": "string (timestamp)"
    }
    ```

-   **Response:**
    ```json
    {
        "message": "Picture deleted successfully"
    }
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
    [ /* array of UserProfile Objects */ ]
    ```
    
### Get Recommended Profiles

-   **URL:** `/api/v1/profiles/recommends`
-   **Method:** `GET`
-   **Request:** Requires Authorization header.
-   **Response:**
    ```json
    [ /* array of UserProfile Objects, sorted by recommendation score */ ]
    ```

### Get a Specific User's Profile

-   **URL:** `/api/v1/users/{userID}/profile`
-   **Method:** `GET`
-   **Request:** URL parameter `userID`. Requires Authorization header.
-   **Response:** `UserProfile Object`
    ```json
    {
        "user_id": "string (UUID)",
        "first_name": "string or null",
        "last_name": "string or null",
        "username": "string or null",
        "gender": "string or null (e.g., male, female, other)",
        "sexual_preference": "string or null (e.g., heterosexual, homosexual, bisexual)",
        "birthday": "string or null (timestamp)",
        "occupation": "string or null",
        "biography": "string or null",
        "fame_rating": "integer or null",
        "location_name": "string or null",
        "distance": "number or null (float)"
    }
    ```

---

## Tags

### Get All Tags

-   **URL:** `/api/v1/tags`
-   **Method:** `GET`
-   **Request:** (No parameters)
-   **Response:**
    ```json
    [ /* array of Tag Objects */ ]
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
    [ /* array of Message Objects */ ]
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
    -   **Server -> Client (Notification):
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

---

## Data Models

### Like Object

Represents a user's like action.

```json
{
    "liker_id": "string (UUID)",
    "liked_id": "string (UUID)",
    "created_at": "string (timestamp)"
}

### Tag Object

Represents a user-defined tag.

```json
{
    "id": "integer",
    "name": "string"
}

### View Object

Represents a user's view action.

```json
{
    "viewer_id": "string (UUID)",
    "viewed_id": "string (UUID)",
    "view_time": "string (timestamp)"
}

### Block Object

Represents a user's block action.

```json
{
    "blocker_id": "string (UUID)",
    "blocked_id": "string (UUID)"
}

### UserData Object

Represents additional user data like location and internal scoring.

```json
{
    "user_id": "string (UUID)",
    "latitude": "number or null (float)",
    "longitude": "number or null (float)",
    "internal_score": "integer or null"
}

### Picture Object

Represents a user's uploaded picture.

```json
{
    "id": "integer",
    "user_id": "string (UUID)",
    "url": "string (URL)",
    "is_profile_pic": "boolean or null",
    "created_at": "string (timestamp)"
}



### Connection Object

Represents a mutual like between two users.

```json
{
    "user1_id": "string (UUID)",
    "user2_id": "string (UUID)",
    "created_at": "string (timestamp)"
}



```
```

### UserProfile Object

Represents a user's profile information.

```json
{
    "user_id": "string (UUID)",
    "first_name": "string or null",
    "last_name": "string or null",
    "username": "string or null",
    "gender": "string or null (e.g., male, female, other)",
    "sexual_preference": "string or null (e.g., heterosexual, homosexual, bisexual)",
    "birthday": "string or null (timestamp)",
    "occupation": "string or null",
    "biography": "string or null",
    "fame_rating": "integer or null",
    "location_name": "string or null",
        "distance": "number or null (float)"
    }
    
    ### Message Object
    
    Represents a chat message.
    
    ```json
    {
        "id": "integer",
        "sender_id": "string (UUID)",
        "recipient_id": "string (UUID)",
        "content": "string",
        "sent_at": "string (timestamp)",
            "is_read": "boolean or null"
        }
        
        ### Notification Object
        
        Represents a user notification.
        
        ```json
        {
            "id": "integer",
            "recipient_id": "string (UUID)",
            "sender_id": "string (UUID) or null",
            "type": "string (e.g., like, view, message)",
            "is_read": "boolean or null",
            "created_at": "string (timestamp)"
        }
        ```    