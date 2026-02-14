# Runbook

## Deployment

### ローカル開発環境

```bash
# 1. 環境変数設定
#    .env, api/.env, wsgateway/.env, filesrv/.env を作成

# 2. バックエンドサービス起動
make up

# 3. フロントエンド起動
make web-install
make web-dev
# → http://localhost:5173 (Vite dev server, API は localhost:80 へプロキシ)
```

### 本番ビルド

```bash
# フロントエンドビルド
make web-build
# → web/dist/ に出力。nginx で配信
```

### Docker Compose サービス全体

```bash
# 起動
make up

# 停止
make down

# 停止 + データ削除 (PostgreSQL ボリューム含む)
make downv
```

## Monitoring

### サービスヘルスチェック

```bash
# 全サービスの状態確認
docker-compose ps

# 特定サービスのログ
docker-compose logs -f api
docker-compose logs -f wsgateway
docker-compose logs -f nginx

# DB 接続確認
docker-compose exec db pg_isready -U $POSTGRES_USER -d $POSTGRES_DB
```

### ツール系 UI

```bash
make tool
```

| Tool | URL | Purpose |
|------|-----|---------|
| Adminer | http://localhost:8080 | DB 管理・クエリ実行 |
| Redis Commander | http://localhost:8081 | Redis データ確認 |

## Common Issues and Fixes

### 1. `make up` 後に API が起動しない

**症状**: `docker-compose ps` で api が restart ループ

**原因**: DB マイグレーションが完了する前に API が接続を試みる

**対処**:
```bash
# DB のヘルスチェック待ち
docker-compose logs -f db
# "database system is ready to accept connections" を確認後
docker-compose restart api
```

### 2. フロントエンド開発サーバーで API 呼び出しが 502

**症状**: `npm run dev` で API リクエストが 502 Bad Gateway

**原因**: バックエンドサービスが起動していない、または nginx がリクエストを転送できない

**対処**:
```bash
# バックエンドが起動しているか確認
make up
docker-compose ps

# nginx のログ確認
docker-compose logs -f nginx
```

### 3. `make web-test` で全テスト失敗

**症状**: `Cannot find package '@/...'` エラー

**原因**: `web/` ディレクトリ外からテスト実行している

**対処**:
```bash
# 正しい実行方法
make web-test
# または
cd web && npm run test:run
```

### 4. PostgreSQL データの完全リセット

**症状**: DB スキーマ変更後にマイグレーションが当たらない

**対処**:
```bash
# ボリューム削除で DB を完全リセット
make downv
make up
```

### 5. WebSocket 接続ができない

**症状**: フロントエンドから WS 接続が確立しない

**原因**: WS Gateway は `Authorization: Bearer` ヘッダーで認証するが、ブラウザの `WebSocket()` API はカスタムヘッダーを送れない

**対処**: nginx で query param (`/ws?token=XXX`) を Authorization ヘッダーに変換する設定が必要（未実装。BE-03 #20 で対応予定）

### 6. シードデータ投入

```bash
make seed
```

DB が起動済みであることが前提。スキーマが変わった場合は `make downv` してから再実行。

## Rollback

### フロントエンドのロールバック

```bash
# 直前のコミットに戻す
git revert HEAD
make web-build

# 特定バージョンに戻す
git checkout <commit-hash> -- web/
make web-build
```

### Docker サービスのロールバック

```bash
# 特定バージョンのイメージでビルドし直し
git checkout <commit-hash>
docker-compose build
make up
```

### DB ロールバック

現時点では自動マイグレーションツールを使用していない。`db/schema/schema.sql` がスキーマの single source of truth。

```bash
# スキーマ変更をロールバックする場合
git checkout <commit-hash> -- db/schema/schema.sql
make downv  # ボリューム削除
make up     # 新スキーマで再作成
make seed   # シードデータ再投入
```

> **Warning**: `make downv` は DB データを完全に削除します。本番環境では pg_dump でバックアップを取ってから実行してください。
