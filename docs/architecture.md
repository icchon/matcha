# Matcha プロジェクトアーキテクチャ

```mermaid
graph TB
    Client[クライアント<br/>Web Browser]
    
    subgraph "Nginx (リバースプロキシ :80)"
        Nginx[Nginx]
    end
    
    subgraph "API サーバー"
        API[API Server<br/>REST API]
        Auth[認証サービス<br/>OAuth/JWT]
        User[ユーザーサービス]
        Profile[プロフィールサービス]
        Chat[チャットサービス]
        Notice[通知サービス]
        Mail[メールサービス]
    end
    
    subgraph "WebSocket ゲートウェイ"
        WSGateway[WebSocket Gateway<br/>/ws]
        Presence[プレゼンス管理]
        ChatSub[チャット購読]
        ReadSub[既読購読]
    end
    
    subgraph "ファイルサーバー"
        FileSrv[File Server<br/>/images, /dog]
    end
    
    subgraph "データストレージ"
        PostgreSQL[(PostgreSQL<br/>メインデータベース)]
        Redis[(Redis<br/>キャッシュ/メッセージング)]
        Storage[ローカルストレージ<br/>uploads/]
    end
    
    subgraph "外部サービス"
        SMTP[SMTP<br/>メール送信]
        Google[Google OAuth]
        GitHub[GitHub OAuth]
    end
    
    Client -->|HTTP/HTTPS| Nginx
    Client -->|WebSocket| Nginx
    
    Nginx -->|/api/v1/*| API
    Nginx -->|/ws| WSGateway
    Nginx -->|/images/, /dog| FileSrv
    Nginx -->|/| Static[静的ファイル<br/>web/index.html]
    
    API --> Auth
    API --> User
    API --> Profile
    API --> Chat
    API --> Notice
    API --> Mail
    
    API --> PostgreSQL
    API --> Redis
    API --> FileSrv
    API --> SMTP
    API --> Google
    API --> GitHub
    
    WSGateway --> Redis
    WSGateway -->|購読| ChatSub
    WSGateway -->|購読| Presence
    WSGateway -->|購読| ReadSub
    
    ChatSub -.->|イベント処理| API
    Presence -.->|イベント処理| API
    ReadSub -.->|イベント処理| API
    
    FileSrv --> Storage
    
    style Client fill:#e1f5ff
    style Nginx fill:#fff4e1
    style API fill:#e8f5e9
    style WSGateway fill:#f3e5f5
    style FileSrv fill:#fce4ec
    style PostgreSQL fill:#e3f2fd
    style Redis fill:#ffebee
    style Storage fill:#f1f8e9
```

## 主要コンポーネント

### 1. Nginx (リバースプロキシ)
- ポート80でリクエストを受信
- ルーティング:
  - `/api/v1/*` → APIサーバー
  - `/ws` → WebSocketゲートウェイ
  - `/images/`, `/dog` → ファイルサーバー
  - `/` → 静的ファイル配信

### 2. APIサーバー
- REST API提供
- 主要機能:
  - 認証 (JWT, OAuth)
  - ユーザー管理
  - プロフィール管理
  - チャット機能
  - 通知機能
  - メール送信

### 3. WebSocketゲートウェイ
- WebSocket接続管理
- Redis経由でメッセージング
- プレゼンス管理
- チャット・既読イベント処理

### 4. ファイルサーバー
- 画像アップロード・配信
- ローカルストレージに保存

### 5. データストレージ
- **PostgreSQL**: メインデータベース
- **Redis**: キャッシュ・メッセージング・プレゼンス管理

### 6. 外部サービス
- **SMTP**: メール送信
- **Google OAuth**: Google認証
- **GitHub OAuth**: GitHub認証
