-- 1. ユーザー基本データ (Core User Data - 内部/非公開情報を含む)

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- 1. ユーザーコアアカウント (users)
-- 外部プロバイダがメールを提供しない場合を考慮し、NOT NULLを削除
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_connection TIMESTAMP WITH TIME ZONE
);

---------------------------------------------------

-- 2. ユーザー機密データ (user_data) - 変更なし

CREATE TABLE user_data (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    latitude DECIMAL(10, 8) NULL, 
    longitude DECIMAL(11, 8) NULL,
    internal_score INT DEFAULT 0
);

CREATE INDEX idx_user_data_location ON user_data USING btree (latitude, longitude);

---------------------------------------------------

-- 3. 認証情報 (Auths) - 変更なし

CREATE TYPE auth_provider_enum AS ENUM ('local', 'google', 'facebook', 'apple', 'github');

CREATE TABLE auths (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider auth_provider_enum NOT NULL,
    provider_uid VARCHAR(255),
    email VARCHAR(255),
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    password_hash VARCHAR(255),
    UNIQUE (user_id, provider),
    UNIQUE (provider, provider_uid),
    UNIQUE (provider, email)
);

-- トークン管理テーブル (FKはusersを参照)
CREATE TABLE verification_tokens (
    token VARCHAR(64) PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE password_resets (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(64) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    PRIMARY KEY (user_id)
);

---------------------------------------------------

-- 4. 公開プロフィール情報 (user_profiles) - 変更なし

CREATE TYPE gender_enum AS ENUM ('male', 'female', 'other');
CREATE TYPE preference_enum AS ENUM ('heterosexual', 'homosexual', 'bisexual');

CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    first_name VARCHAR(50) NULL, 
    last_name VARCHAR(50) NULL, 
    username VARCHAR(50) NULL, 

    gender gender_enum NULL,
    sexual_preference preference_enum NULL,
    biography TEXT NULL,
    
    fame_rating INT DEFAULT 0 CHECK (fame_rating >= 0), 
    location_name VARCHAR(255)
);

-- (5. 関係性テーブル, 6. チャットテーブルも同様に users テーブルを参照)

---------------------------------------------------

-- 4. 興味タグと写真 (Tags & Pictures)

CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE user_tags (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tag_id)
);

CREATE TABLE pictures (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    url VARCHAR(255) NOT NULL,
    is_profile_pic BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---------------------------------------------------

-- 5. 関係性、履歴、通知 (Relationships, History, & Notifications)

CREATE TABLE likes (
    liker_id UUID REFERENCES users(id) ON DELETE CASCADE,
    liked_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (liker_id, liked_id),
    CHECK (liker_id <> liked_id)
);

CREATE TABLE connections (
    user1_id UUID REFERENCES users(id) ON DELETE CASCADE,
    user2_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user1_id, user2_id),
    CHECK (user1_id < user2_id)
);

CREATE TABLE views (
    viewer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    viewed_id UUID REFERENCES users(id) ON DELETE CASCADE,
    view_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (viewer_id, viewed_id, view_time)
);

CREATE TABLE blocks (
    blocker_id UUID REFERENCES users(id) ON DELETE CASCADE,
    blocked_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id <> blocked_id)
);

CREATE TYPE notification_type_enum AS ENUM ('like', 'view', 'match', 'unlike', 'message');

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    type notification_type_enum NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

---------------------------------------------------

-- 6. チャット機能 (Chat)

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    sent_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_read BOOLEAN DEFAULT FALSE
);

-- インデックスの最適化
CREATE INDEX idx_user_data_location ON user_data USING btree (latitude, longitude);
CREATE INDEX idx_messages_chat_history ON messages (sender_id, recipient_id, sent_at);

---------------------------------------------------

-- 7. リフレッシュトークン (Refresh Tokens)
CREATE TABLE refresh_tokens (
    token_hash VARCHAR(255) PRIMARY KEY,
    user_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

