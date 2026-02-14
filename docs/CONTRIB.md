# Contributing Guide

## Prerequisites

- Docker & Docker Compose
- Node.js (v20+)
- Go 1.21+ (バックエンド開発時)

## Environment Setup

### 1. 環境変数

各サービスの `.env` ファイルを設定する。`.env.example` は未作成のため、以下のキー一覧を参考に設定:

#### Root `.env` (docker-compose 用)

| Variable | Purpose |
|----------|---------|
| `POSTGRES_USER` | PostgreSQL ユーザー名 |
| `POSTGRES_PASSWORD` | PostgreSQL パスワード |
| `POSTGRES_DB` | PostgreSQL データベース名 |

#### `api/.env`

| Variable | Purpose |
|----------|---------|
| `SERVER_ADDR` | API サーバーアドレス (例: `:8080`) |
| `DATABASE_URL` | PostgreSQL 接続文字列 |
| `REDIS_ADDR` | Redis アドレス (例: `redis:6379`) |
| `JWT_SIGNING_KEY` | JWT 署名鍵 |
| `HMAC_SECRET_KEY` | HMAC 署名鍵 |
| `GOOGLE_CLIENT_ID` | Google OAuth クライアント ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth クライアントシークレット |
| `GITHUB_CLIENT_ID` | GitHub OAuth クライアント ID |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth クライアントシークレット |
| `REDIRECT_URI` | OAuth リダイレクト URI |
| `SMTP_HOST` | SMTP ホスト |
| `SMTP_PORT` | SMTP ポート |
| `SMTP_USERNAME` | SMTP ユーザー名 |
| `SMTP_PASSWORD` | SMTP パスワード |
| `SMTP_SENDER` | 送信元メールアドレス |
| `IMAGE_UPLOAD_ENDPOINT` | ファイルサーバーのアップロード URL |
| `BASE_URL` | アプリの公開 URL |

#### `wsgateway/.env`

| Variable | Purpose |
|----------|---------|
| `SERVER_ADDR` | WebSocket ゲートウェイアドレス |
| `REDIS_ADDR` | Redis アドレス |
| `JWT_SIGNING_KEY` | JWT 署名鍵 (API と同じ値) |

#### `filesrv/.env`

| Variable | Purpose |
|----------|---------|
| `UPLOAD_DIR` | アップロードファイル保存ディレクトリ |

### 2. サービス起動

```bash
# バックエンド全サービス起動
make up

# フロントエンド依存パッケージインストール
make web-install

# フロントエンド開発サーバー起動
make web-dev
```

## Available Scripts

### Make targets (ルート)

| Command | Description |
|---------|-------------|
| `make up` | Docker Compose で全サービス起動 |
| `make down` | 全サービス停止 |
| `make downv` | 全サービス停止 + ボリューム削除 |
| `make tool` | ツール系サービス起動 (Adminer, Redis Commander, tbls) |
| `make diagram` | tbls で DB スキーマドキュメント生成 |
| `make seed` | DB シードデータ投入 |
| `make fmt` | Go コードフォーマット |
| `make web-install` | フロントエンド依存パッケージインストール |
| `make web-dev` | フロントエンド開発サーバー起動 (Vite) |
| `make web-build` | フロントエンドビルド (TypeScript + Vite) |
| `make web-test` | フロントエンドテスト実行 (vitest run) |
| `make web-lint` | フロントエンド lint (ESLint) |
| `make web-format` | フロントエンドフォーマット (Prettier) |

### npm scripts (`web/`)

| Script | Command | Description |
|--------|---------|-------------|
| `dev` | `vite` | 開発サーバー起動 (HMR) |
| `build` | `tsc -b && vite build` | TypeScript チェック + 本番ビルド |
| `lint` | `eslint .` | ESLint 実行 |
| `format` | `prettier --write "src/**/*.{ts,tsx,css}"` | Prettier フォーマット |
| `preview` | `vite preview` | ビルド成果物プレビュー |
| `test` | `vitest` | テスト (watch モード) |
| `test:run` | `vitest run` | テスト (単発実行) |
| `test:coverage` | `vitest run --coverage` | テスト + カバレッジ |

## Development Workflow

### ブランチ戦略

```
main
  └── feat/fe-XX-feature-name   (フィーチャーブランチ)
```

### TDD フロー

1. **RED**: テストファイルを先に作成、期待する振る舞いを assertion で記述
2. **GREEN**: テストを通す最小限の実装
3. **IMPROVE**: リファクタリング（テストが緑のまま）

### コミット規約

Conventional Commits を使用:

- `feat:` — 新機能
- `fix:` — バグ修正
- `refactor:` — リファクタリング
- `docs:` — ドキュメント
- `test:` — テスト
- `chore:` — その他
- `perf:` — パフォーマンス改善
- `ci:` — CI/CD

### PR フロー

1. フィーチャーブランチで開発
2. `make web-test` + `make web-lint` パス確認
3. PR 作成（テストプラン記載）
4. code-reviewer + security-reviewer 通過
5. マージ

## Testing

### フロントエンドテスト

```bash
# 全テスト実行
make web-test

# watch モードで開発中テスト
cd web && npm test

# カバレッジ付き
cd web && npm run test:coverage
```

**テストスタック**: Vitest + Testing Library + jsdom

**テストファイル配置**: `web/src/tests/` にソースと対応するディレクトリ構造で配置

### モック戦略

- API 呼び出し: `vi.mock('@/api/client')`
- Zustand Store: `vi.mock('@/stores/authStore')`
- Router: `MemoryRouter` でナビゲーションテスト

## Architecture

[docs/architecture.md](./architecture.md) を参照。

### サービス構成

| Service | Port | Description |
|---------|------|-------------|
| nginx | 80 | リバースプロキシ |
| api | (内部) | REST API サーバー (Go) |
| wsgateway | (内部) | WebSocket ゲートウェイ (Go) |
| filesrv | (内部) | ファイルサーバー (Go) |
| db | (内部) | PostgreSQL 13 |
| redis | (内部) | Redis 7 |

### ツールサービス (`make tool`)

| Service | Port | Description |
|---------|------|-------------|
| adminer | 8080 | DB 管理 UI |
| redis-commander | 8081 | Redis 管理 UI |
| tbls | — | DB スキーマドキュメント生成 |

### フロントエンド技術スタック

| Technology | Version | Purpose |
|------------|---------|---------|
| React | 19 | UI ライブラリ |
| TypeScript | 5.9 | 型安全 |
| Vite | 7 | ビルドツール |
| Tailwind CSS | 4 | スタイリング |
| Zustand | 5 | 状態管理 |
| React Router | 7 | ルーティング |
| React Hook Form | 7 | フォーム管理 |
| Zod | 4 | バリデーション |
| Sonner | 2 | Toast 通知 |
| Vitest | 4 | テストフレームワーク |
| Testing Library | 16 | コンポーネントテスト |
