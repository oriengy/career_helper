# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

恋爱翻译官 (Romance Translator) - An AI-powered social assistance application that helps users better understand communications with the opposite sex. The system consists of a Next.js frontend and a Go backend using gRPC/Connect protocol.

## Architecture

### Frontend: Next.js 14 (App Router)
- **Location**: `frontend/`
- **Framework**: Next.js 14 with TypeScript
- **UI Library**: TDesign React + Tailwind CSS
- **State Management**: Zustand with persist middleware
- **HTTP Client**: Axios with Connect protocol support

### Backend: Go + gRPC/Connect
- **Location**: `chat_server-main/app_server/`
- **Framework**: Gin Web Framework
- **RPC Protocol**: gRPC + Connect-Go
- **Database**: MySQL + GORM ORM
- **Authentication**: JWT tokens (365 day expiry)

### Key Architectural Patterns

**Frontend Structure**:
- `src/app/`: Next.js App Router pages (gender, login, sessions)
- `src/components/`: Reusable React components
- `src/services/api/`: API client layer (auth, message, profile, session, translate, upload)
- `src/stores/`: Zustand state stores (auth, user, ui)
- `src/types/`: TypeScript type definitions

**Backend Structure**:
- `cmd/api_server/main.go`: Application entry point
- `service/`: Business logic layer (auth, chat, message, profile, translate, user)
- `domain/`: Domain-driven design layer with transaction management
- `model/`: GORM database models
- `pkg/`: Reusable packages (jwt, db, openaic, ossc, aiapi, etc.)
- `proto/`: Generated Protocol Buffers code

**Service Layer Dependencies**: Service → Domain → Model (strict one-way flow)

## Common Development Commands

### Frontend (run from `frontend/`)

```bash
# Install dependencies
pnpm install  # or npm install

# Development server (http://localhost:3000)
pnpm dev

# Type checking
pnpm type-check  # or npm run type-check

# Linting
pnpm lint

# Build production
pnpm build
pnpm start
```

### Backend (run from `chat_server-main/app_server/`)

```bash
# Install dependencies
go mod tidy

# Run development server
go run cmd/api_server/main.go -c config.yaml

# Build binary
go build -o output/api_server cmd/api_server/main.go

# Run tests
go test ./...

# Integration testing script
go run scripts/integration_testing.go
```

### Protocol Buffers (run from `chat_server-main/`)

```bash
# Regenerate Connect/gRPC code from proto definitions
buf generate
```

### Docker Deployment

```bash
# Build and run both services
docker compose up -d --build

# Frontend runs on port 3000
# Backend runs on port 8082

# View logs
docker compose logs -f

# Stop services
docker compose down
```

**Common Issues**:
- Port conflicts: See `PORT-ISSUE-SOLUTIONS.md` for quick fixes
- Auto-fix script: `./fix-port-issues.sh`
- Alternative ports: Use `docker-compose.alternative-ports.yml`
- Production with Nginx: Use `docker-compose.yml`

## Key Technical Details

### Frontend-Backend Communication

The frontend communicates with the backend using Connect protocol (gRPC-web):
- **Base URL**: Configured via `NEXT_PUBLIC_API_BASE_URL` environment variable
- **Protocol**: Connect-RPC over HTTP
- **Headers**: All requests include `Authorization: Bearer {token}`, `X-App-Platform: web`, `Connect-Protocol-Version: 1`

### Authentication Flow

1. User enters phone number and requests verification code
2. Backend sends 6-digit code (development: magic code `1234` works for any number)
3. User submits code for verification
4. Backend returns JWT token (365 day validity)
5. Token stored in Zustand + localStorage for auto-login
6. If user hasn't set gender, redirect to gender selection page
7. Otherwise, redirect to sessions list

### State Management

Uses Zustand with persistence:
- `stores/auth.ts`: Authentication token and login state
- `stores/user.ts`: User profile information
- `stores/ui.ts`: UI state (loading, modals, etc.)

All stores automatically persist to localStorage.

### Message System Architecture

**Unified Message Model**: All messages (history, translations, consultations) are stored in the `consult_message` table with different `msg_type` values.

**Message Roles**:
- `SELF`: Messages from the user
- `FRIEND`: Messages from the conversation partner
- `AI`: AI-generated responses

**Translation Flow**:
1. User clicks "translate" on a friend's message
2. Backend fetches 24-hour conversation context
3. AI analyzes context and translates based on gender perspective
4. Translation result saved as new message with `ParentID` linking to original
5. Frontend displays translation inline

### File Upload & Storage

**OSS Integration**: Uses Aliyun OSS for file storage
- **Temporary uploads**: `uploads/` directory
- **User files**: `user/{user_id}/` directory after migration
- **Image compression**: Avatars (800x800, quality 0.8), Chat images (1200x1200, quality 0.85)

**Upload endpoints**:
- `/file/wx_upload`: WeChat-compatible multipart upload

### Database Models

**Core tables**:
- `user`: User accounts with WeChat OpenID/UnionID
- `profile`: Personal profiles (user's own + conversation partners)
- `chat_session`: Conversation sessions linking user + profile
- `consult_message`: Unified message storage

**ID Generation**: Uses Snowflake algorithm for distributed unique IDs

## Environment Configuration

### Frontend Environment Variables

Development (`.env.local`):
```
NEXT_PUBLIC_API_BASE_URL=https://local.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=development
```

Production (`.env.production`):
```
NEXT_PUBLIC_API_BASE_URL=https://api.chathandy.com
NEXT_PUBLIC_APP_VERSION=1.0.0
NEXT_PUBLIC_APP_ENV=production
```

### Backend Configuration

Main config file: `chat_server-main/app_server/config.yaml`

Required sections:
- `db.dsn`: MySQL connection string
- `jwt.secret`: JWT signing key
- `wechat.mp.translator`: WeChat Mini Program credentials
- `ai.volces`: Volcano Engine AI API configuration (primary AI provider)
- `aliyun.oss`: Aliyun OSS credentials and bucket names
- `server.address`: Server bind address (default `:8080`)

## API Services

All backend services follow Connect-RPC protocol:

- `user.UserService`: User authentication and profile retrieval
- `chat.ChatService`: Chat session CRUD operations
- `message.ConsultMessageService`: Consultation message queries
- `message.FriendMessageService`: Friend message management + OCR parsing
- `translate.TranslateService`: AI translation services
- `profile.ProfileService`: Personal profile management
- `/file/wx_upload`: HTTP file upload endpoint
- `/ping`: Health check endpoint

## Development Practices

### Frontend
- Use functional React components with hooks
- Maintain TypeScript strict mode
- Follow existing TDesign component patterns
- API calls go through `services/api/` layer, never directly from components
- Error handling: 401 → auto-logout, others → toast notification

### Backend
- Follow Go standard formatting (`gofmt`)
- Service layer must validate all inputs
- Complex operations use Domain layer transaction wrappers
- Use structured logging with `slog`
- Never commit real secrets to config files

### Commit Messages
- Keep commits concise and descriptive
- Messages often in Chinese (project convention)
- Include UI screenshots for frontend PRs

## Important Notes

- **NEVER** update the git config
- **NEVER** commit sensitive credentials (use example values)
- Frontend uses TypeScript path alias `@/*` mapped to `./src/*`
- Backend imports use module path `app_server/` prefix
- Protocol buffer definitions live in `chat_server-main/proto/`
- For detailed backend architecture, see `chat_server-main/CLAUDE.md`
- For detailed frontend structure, see `frontend/README.md`
