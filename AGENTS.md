# Repository Guidelines

## Project Structure & Module Organization
- `frontend/` is the Next.js 14 (App Router) web client. Source lives in `frontend/src/` with pages in `frontend/src/app/`, UI in `frontend/src/components/`, services in `frontend/src/services/`, and state in `frontend/src/stores/`.
- `chat_server-main/` is the Go backend (Gin + ConnectRPC). The main entry is `chat_server-main/app_server/cmd/api_server/main.go`. Core layers: `app_server/service/`, `app_server/domain/`, `app_server/model/`, `app_server/pkg/`. Protos are in `chat_server-main/proto/` and generated output in `chat_server-main/app_server/proto/`.
- Static assets for the web app are in `frontend/public/`.

## Build, Test, and Development Commands
Frontend (run from `frontend/`):
- `npm install` or `pnpm install`: install dependencies.
- `npm run dev`: start the dev server at `http://localhost:3000`.
- `npm run build` / `npm run start`: build and run production.
- `npm run lint`: run Next.js ESLint.
- `npm run type-check`: run TypeScript checks.

Backend (run from `chat_server-main/app_server/`):
- `go run cmd/api_server/main.go -c config.yaml`: start the API server.
- `go build -o output/api_server cmd/api_server/main.go`: build a binary.
- `buf generate` (from `chat_server-main/`): regenerate Connect/Swagger artifacts.

## Coding Style & Naming Conventions
- Frontend uses TypeScript, Next.js, Tailwind CSS, and TDesign. Prefer functional React components and hooks.
- Run `npm run lint` before PRs; keep formatting consistent with existing code.
- Backend follows standard Go formatting (`gofmt`) and idiomatic package naming.

## Testing Guidelines
- Frontend: no dedicated test runner configured; use `npm run lint` and `npm run type-check` as baseline checks.
- Backend: run `go test ./...` from `chat_server-main/app_server/` if you add or touch Go code.
- Integration script exists at `chat_server-main/app_server/scripts/integration_testing.go` (run as needed).

## Commit & Pull Request Guidelines
- Git history uses short, descriptive commit messages (often in Chinese). Keep messages brief and action-oriented.
- PRs should describe scope and impact, reference any related issue, and include UI screenshots for frontend changes.

## Configuration & Security Notes
- Frontend env files: `frontend/.env.local` for dev, `frontend/.env.production` for prod.
- Backend config: `chat_server-main/app_server/config.yaml`.
- Never commit real secrets; use example values and document required keys in the PR.
