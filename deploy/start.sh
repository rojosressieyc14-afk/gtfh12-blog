#!/bin/bash
set -euo pipefail

echo "=== PulseBlog Production Start ==="

# Check dependencies
command -v docker >/dev/null 2>&1 || { echo "需要安装 Docker"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "需要安装 docker-compose"; exit 1; }

# Load production env if exists
if [ -f .env.prod ]; then
    export $(grep -v '^\s*#' .env.prod | xargs)
fi

# Build frontends
echo ">>> Building web frontend..."
cd web && npm ci && npm run build && cd ..

echo ">>> Building admin frontend..."
cd admin && npm ci && npm run build && cd ..

# Start services
echo ">>> Starting Docker services..."
docker-compose up --build -d

echo "=== Deployment complete ==="
echo "Web:  http://localhost"
echo "Admin: http://localhost/admin"
echo "API:  http://localhost/api"
