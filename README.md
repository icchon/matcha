# Matcha

42 School のマッチングアプリプロジェクト。

## Architecture

```
nginx (reverse proxy :80)
├── /api/v1/*   → api (Go)
├── /images/*   → filesrv (Go)
├── /ws         → wsgateway (Go)
└── /*          → web (React SPA)

db: PostgreSQL 13
cache: Redis 7
```

| Service | 言語/FW | 概要 |
|---------|---------|------|
| `web` | React 19 + Vite + TypeScript | SPA フロントエンド |
| `api` | Go | REST API |
| `wsgateway` | Go | WebSocket ゲートウェイ |
| `filesrv` | Go | ファイルアップロード/配信 |
| `nginx` | Nginx | リバースプロキシ + SPA 配信 |
| `db` | PostgreSQL 13 | メインDB |
| `redis` | Redis 7 | Pub/Sub + キャッシュ |

## Prerequisites

- Docker & Docker Compose
- Node.js 22+ (フロントエンド開発用)
- Go 1.22+ (バックエンド開発用)

## Setup

```bash
# .env ファイルを作成 (Makefile が include する)
cp .env.example .env  # または手動で作成

# 全サービス起動
make up

# 管理ツール (Adminer, Redis Commander) も起動
make tool
```

## Development

### Frontend (web/)

```bash
make web-install   # 依存インストール
make web-dev       # 開発サーバー起動
make web-build     # プロダクションビルド
make web-test      # テスト実行
make web-lint      # ESLint
make web-format    # Prettier
```

### Backend

```bash
make fmt           # Go コードフォーマット
```

### Infrastructure

```bash
make up            # Docker Compose 起動
make down          # 停止
make downv         # 停止 + ボリューム削除
make seed          # シードデータ投入
make diagram       # tbls で DB ドキュメント生成
```

## Docs

- [Frontend Plan](docs/frontend_plan.md)
- [Backend Plan](docs/backend_plan.md)
- [DB Schema](db/tbls/)
