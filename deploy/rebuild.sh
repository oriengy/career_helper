#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
COMPOSE_FILE="$ROOT_DIR/deploy/docker-compose.nginx.yml"

cd "$ROOT_DIR"

echo "[1/4] Pull latest code"
git pull 

echo "[2/4] Stop existing containers"
docker compose -f "$COMPOSE_FILE" down --remove-orphans

echo "[3/4] Build images"
docker compose -f "$COMPOSE_FILE" build

echo "[4/4] Start services"
docker compose -f "$COMPOSE_FILE" up -d
echo "[5/5] Done"
