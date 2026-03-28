# Copilot Instructions

このファイルはリポジトリ内のコーディング規約・設計方針をまとめたものです。
GitHub Copilot はこれらの規約に従ってコード補完・生成を行ってください。

---

## プロジェクト概要

**Giter** は Go で実装された Web アプリケーションです。
GitHub のコミット履歴を分析し、コミットのストリーク（連続日数）を可視化します。

- **言語**: Go 1.23.4
- **HTTP フレームワーク**: [Gin](https://github.com/gin-gonic/gin) v1.10.0
- **ORM**: [GORM](https://gorm.io) v1.25.12（PostgreSQL）
- **認証**: JWT（`dgrijalva/jwt-go`）+ bcrypt + GitHub OAuth2（Goth）
- **ロギング**: [zerolog](https://github.com/rs/zerolog) v1.33.0
- **DI**: [Google Wire](https://github.com/google/wire) v0.6.0
- **GitHub API**: REST（`google/go-github`）+ GraphQL（`hasura/go-graphql-client`）

---

## アーキテクチャ

3 層アーキテクチャ（Repository パターン）を採用しています。

```
Controller（HTTP ハンドラ・リクエスト/レスポンス処理）
    ↓
Service（ビジネスロジック・処理の調整）
    ↓
Repository（データアクセス・DB クエリ・外部 API）
    ↓
PostgreSQL / GitHub API
```

各層はインターフェースで依存関係を定義し、Google Wire でコンパイル時に注入します。

---

## ディレクトリ構成

```
controllers/    # HTTP ハンドラ（Controller 層）
services/       # ビジネスロジック（Service 層）
repositories/   # データアクセス（Repository 層）
models/         # GORM モデル定義
dto/            # Data Transfer Objects（リクエスト入力など）
middlewares/    # Gin ミドルウェア
di/             # Wire による DI 設定（wire.go / wire_gen.go）
infra/          # DB 接続などインフラ設定
initializer/    # アプリ初期化・ロガー設定
utils/token/    # JWT ユーティリティ
query/          # GraphQL クエリ構造体
public/         # 汎用ユーティリティ関数
templates/      # Gin HTML テンプレート（*.tmpl）
migrations/     # GORM AutoMigrate スクリプト
```

---

## ファイル命名規約

| 種別 | パターン | 例 |
|---|---|---|
| Controller | `{リソース名}_controller.go` | `auth_controller.go` |
| Service | `{リソース名}_service.go` | `auth_service.go` |
| Repository | `{リソース名}_repository.go` | `auth_repository.go` |
| Middleware | `{機能名}_middleware.go` | `auth_middleware.go` |
| Model | `{モデル名}.go` | `user.go` |
| DTO | `{ドメイン名}.go` | `auth.go` |
| ユーティリティ | `{機能名}.go` | `token.go` |
| HTML テンプレート | `{画面名}.tmpl` | `mypage.tmpl` |

---

## 命名規約

### 構造体・インターフェース

- **インターフェース**: `I` プレフィックス + PascalCase（例: `IAuthService`, `IRequestController`）
- **具体型（struct）**: PascalCase（例: `AuthService`, `RequestController`）
- **非公開型**: camelCase の struct（例: `authService`）

```go
// インターフェース
type IAuthService interface {
    Register(user *models.User) (*models.User, error)
    Login(user *dto.LoginInput) (string, error)
}

// 公開 struct（DI コンストラクタが返す型）
type AuthService struct {
    repository repositories.IAuthRepository
    logger     zerolog.Logger
}
```

### 関数・メソッド

- **コンストラクタ**: `New{型名}()` 形式（例: `NewAuthController()`, `NewAuthService()`）
- **公開メソッド**: PascalCase（例: `GetCommits()`, `Register()`）
- **非公開メソッド/変数**: camelCase（例: `getLogFile()`, `dsn`）

### 変数

- Gin のコンテキスト引数: `ctx`
- エラー変数: `err`
- ロガー変数: `logger`
- DB 変数: `db`

---

## コンストラクタパターン

各層は `New{型名}` 関数でインスタンスを生成し、インターフェース型で返します。
ロガーは `initializer.DefaultLogger()` で初期化します。

```go
func NewAuthController(service services.IAuthService) IAuthControler {
    return &AuthController{
        service: service,
        logger:  initializer.DefaultLogger(),
    }
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
    return &AuthService{
        repository: repository,
        logger:     initializer.DefaultLogger(),
    }
}

func NewAuthRepository(db *gorm.DB) IAuthRepository {
    return &AuthRepository{
        db:     db,
        logger: initializer.DefaultLogger(),
    }
}
```

---

## エラー処理

- Go 標準の `if err != nil` パターンを使用します。
- 致命的なエラー（DB 接続失敗など）は `panic` を使用します。
- HTTP ハンドラでエラーが発生した場合は `ctx.JSON` でエラーレスポンスを返し `return` します。

```go
// 致命的なエラー
if err != nil {
    panic("Failed to connect database")
}

// HTTP ハンドラでのエラー処理
if err := ctx.ShouldBind(&input); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}

// Service / Repository でのエラー処理
if err != nil {
    return nil, err
}
```

---

## ミドルウェアパターン

`gin.HandlerFunc` を返すクロージャ形式で実装します。

```go
func JwtAuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        err := token.TokenValid(ctx)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
            ctx.Abort()
            return
        }
        ctx.Next()
    }
}
```

**ミドルウェア適用順序（main.go）:**

1. `requestid.New()` — リクエスト ID 付与
2. `middlewares.LogMiddleware()` — リクエスト開始/終了のログ記録
3. `middlewares.JwtAuthMiddleware()` — 保護ルートへの JWT 認証

---

## ロギング

zerolog を使用します。ロガーは各 struct フィールドに保持し、`initializer.DefaultLogger()` で初期化します。

```go
// ログ出力
logger.Info().Msg("REQUEST START")
logger.Error().Msg("error")

// ロガーの初期化（コンストラクタ内）
logger: initializer.DefaultLogger()
```

---

## データベース（GORM）

- `gorm.Model` を埋め込んで `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt` を自動付与します。
- フィールドタグで DB 制約と JSON キーを定義します。
- クエリは Repository 層で行い、Service / Controller には DB インスタンスを渡しません。

```go
type User struct {
    gorm.Model
    Email        string `gorm:"size:255;not null;unique" json:"email"`
    Password     string `gorm:"size:255;"               json:"password"`
    IsGithubUser bool   `gorm:"not null"                json:"is_github_user"`
    CanSendEmail bool   `gorm:"not null"                json:"can_send_email"`
}
```

---

## DTO（Data Transfer Object）

リクエストの入力値は DTO 構造体で受け取ります。
`json`, `form`, `binding` タグを付与します。

```go
type RegisterInput struct {
    Email    string `json:"email"    form:"email"    binding:"required"`
    Password string `json:"password" form:"password" binding:"required"`
}
```

---

## HTTP レスポンス

- JSON レスポンスは `ctx.JSON(http.StatusXXX, gin.H{...})` を使用します。
- HTML テンプレートのレンダリングは `ctx.HTML(http.StatusOK, "name.tmpl", gin.H{...})` を使用します。
- テンプレートは `templates/*.tmpl` に配置します。

```go
// JSON レスポンス
ctx.JSON(http.StatusOK, gin.H{"data": user.PrepareOutput()})

// HTML レンダリング
ctx.HTML(http.StatusOK, "mypage.tmpl", gin.H{"token": token})

// リダイレクト
ctx.Redirect(http.StatusMovedPermanently, "/")
```

---

## 依存性注入（Google Wire）

`di/wire.go` でプロバイダセットを定義し、`wire gen` で `di/wire_gen.go` を生成します。
各コンストラクタは `wire.Build` に渡し、Wire が自動的に依存関係を解決します。

```go
func InitAuthRouter(db *gorm.DB) controllers.IAuthControler {
    wire.Build(
        controllers.NewAuthController,
        services.NewAuthService,
        repositories.NewAuthRepository,
    )
    return &controllers.AuthController{}
}
```

---

## インポート順序

標準ライブラリ → 内部パッケージ → サードパーティの順でグループ化します。

```go
import (
    // 標準ライブラリ
    "fmt"
    "net/http"
    "os"

    // 内部パッケージ
    "giter/dto"
    "giter/initializer"
    "giter/services"

    // サードパーティ
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
)
```

---

## コメント規約

- コメントは**日本語**を基本とします。
- 関数・メソッドの説明はその上の行に記述します。
- 複雑なロジックにはインラインコメントを付けます。
- `TODO:` コメントで未実装事項を記録します。

```go
// パスワードをハッシュ化する
hashedPassword, err := bcrypt.GenerateFromPassword(...)

// TODO: 日付を跨いで実行している場合、前日のファイルに書き込みをしてしまう
func Log() { ... }
```

---

## 環境変数

`.env` ファイル（`godotenv` で読み込み）で管理します。
環境変数名はすべて大文字のスネークケースです。

| 変数名 | 用途 |
|---|---|
| `GITHUB_KEY` | GitHub OAuth2 クライアント ID |
| `GITHUB_SECRET` | GitHub OAuth2 クライアントシークレット |
| `API_SECRET` | JWT 署名キー |
| `HOST_NAME` | PostgreSQL ホスト名 |
| `POSTGRES_USER` | PostgreSQL ユーザー名 |
| `POSTGRES_PASSWORD` | PostgreSQL パスワード |
| `POSTGRES_DATABASE` | PostgreSQL データベース名 |
| `POSTGRES_PORT` | PostgreSQL ポート番号 |

`.env` ファイルは `.gitignore` に含まれています。コミットしないでください。

---

## 開発環境

- **Docker Compose**: アプリ（ポート 8888）と PostgreSQL 16（ポート 5432）を起動します。
- **Air**: ホットリロードで開発します（`Dockerfile` にインストール済み）。
- **タイムゾーン**: `Asia/Tokyo`（JST）を基準とします。
- **DB マイグレーション**: `migrations/migration.go` を個別に実行します。

```bash
docker-compose up   # 開発環境起動
```

---

## GraphQL クエリ

`graphql` タグを使用してクエリ変数を定義します。変数名は `$変数名` 形式です。

```go
type CommitsQuery struct {
    User User `json:"user" graphql:"user(login: $USER_NAME)"`
}
```

---

## JWT 認証フロー

1. ログイン成功時に `token.GenerateToken(userID)` でトークンを生成します。
2. レスポンスの Cookie（`jwt`）にトークンをセットします。
3. 保護ルートでは `JwtAuthMiddleware()` がトークンを検証します。
4. `token.ExtractTokenId(ctx)` でトークンからユーザー ID を取得します。

```go
ctx.SetCookie("jwt", token, 3600, "/", "localhost", true, true)
```
