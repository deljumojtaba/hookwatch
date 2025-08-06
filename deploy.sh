#!/bin/bash
# Hookwatch Deployment Script

echo "🚀 Deploying Hookwatch Project..."

# Update docker-compose.yaml
cat > docker-compose.yaml << 'EOF'
services:
  redis:
    image: redis:latest
    container_name: redis
    restart: unless-stopped
    tty: true
    ports:
      - "${REDIS_PORT:-6379}:6379"
    volumes:
      - ./init/redis.conf:/usr/local/etc/redis/redis.conf
      - redis_data:/data
    command: redis-server /usr/local/etc/redis/redis.conf
    environment:
      - REDIS_PASSWORD=${REDIS_PASSWORD:-}
    networks:
      - main-network

  mongodb:
    image: mongo:latest
    container_name: mongodb
    ports:
      - "${MONGO_INITDB__PORT:-27017}:27017"
    volumes:
      - ./data/mongodb_data:/data/db
      - ./init/mongo.js:/docker-entrypoint-initdb.d/mongo.js
      - mongodb_data:/data/db
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${MONGO_INITDB_ROOT_USERNAME:-admin}
      MONGO_INITDB_ROOT_PASSWORD: ${MONGO_INITDB_ROOT_PASSWORD:-SecurePassword123!}
      MONGO_INITDB_DATABASE: ${MONGO_INITDB_DATABASE:-hookwatch}
    networks:
      - main-network
    restart: unless-stopped

  hookwatch:
    build:
      context: ./hookwatch
      dockerfile: Dockerfile
    container_name: hookwatch
    restart: unless-stopped
    expose:
      - "3000"
    environment:
      - PORT=3000
      - MONGO_URI=mongodb://${MONGO_INITDB_ROOT_USERNAME}:${MONGO_INITDB_ROOT_PASSWORD}@mongodb:27017
      - REDIS_URI=redis://redis:6379
      - VIRTUAL_HOST=api.hookwatch.antcoders.dev
      - LETSENCRYPT_HOST=api.hookwatch.antcoders.dev
      - LETSENCRYPT_EMAIL=antcoderstoken@gmail.com
    depends_on:
      - mongodb
      - redis
    networks:
      - main-network
      - nginx-proxy

  web-ui:
    image: nginx:alpine
    container_name: hookwatch-web-ui
    restart: unless-stopped
    expose:
      - "80"
    environment:
      - VIRTUAL_HOST=hookwatch.antcoders.dev
      - LETSENCRYPT_HOST=hookwatch.antcoders.dev
      - LETSENCRYPT_EMAIL=antcoderstoken@gmail.com
    volumes:
      - ./web:/usr/share/nginx/html
    networks:
      - main-network
      - nginx-proxy

networks:
  main-network:
    driver: bridge
  nginx-proxy:
    name: nginx-proxy
    external: true

volumes:
  mongodb_data:
    driver: local
  redis_data:
    driver: local
EOF

# Update web/script.js with correct API URL (HTTP first, then HTTPS)
echo "📝 Updating API URL in script.js..."
sed -i 's|const API_BASE_URL = "http://localhost:3000";|const API_BASE_URL = "http://api.hookwatch.antcoders.dev";|g' web/script.js

# Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📋 Creating .env file..."
    cat > .env << 'EOF'
# Database Ports (accessible from host)
REDIS_PORT=6379
MONGO_INITDB__PORT=27017

# MongoDB Configuration
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=SecurePassword123!
MONGO_INITDB_DATABASE=hookwatch

# Redis Configuration
REDIS_PASSWORD=

# Hookwatch Configuration
HOOKWATCH_PORT=3000

# Web UI Configuration
WEB_UI_PORT=8080

# Development/Debug Settings
NODE_ENV=production
GIN_MODE=release
EOF
fi

echo "✅ Files updated successfully!"
echo ""

# Check which docker compose command is available
if command -v docker-compose &> /dev/null; then
    COMPOSE_CMD="docker-compose"
elif docker compose version &> /dev/null 2>&1; then
    COMPOSE_CMD="docker compose"
else
    echo "❌ Docker Compose not found! Please install it first:"
    echo "sudo apt update && sudo apt install docker-compose -y"
    exit 1
fi

echo "🔄 Starting deployment..."

# Stop existing containers
echo "⏹️ Stopping existing containers..."
$COMPOSE_CMD down

# Build and start services
echo "🚀 Building and starting services..."
$COMPOSE_CMD up -d --build

# Wait for services to start
echo "⏱️ Waiting for services to start..."
sleep 10

# Check status
echo "📊 Checking service status..."
$COMPOSE_CMD ps

echo ""
echo "🌐 Your services should be available at:"
echo "• Web UI: http://hookwatch.antcoders.dev (HTTP first, HTTPS will be ready in 2-3 minutes)"
echo "• API: http://api.hookwatch.antcoders.dev (HTTP first, HTTPS will be ready in 2-3 minutes)"
echo ""
echo "🔍 Test commands:"
echo "curl http://api.hookwatch.antcoders.dev/health"
echo "curl http://hookwatch.antcoders.dev"
echo ""
echo "🗄️ Database Access:"
echo "• MongoDB: mongodb://admin:SecurePassword123!@localhost:27017"
echo "• Redis: redis://localhost:6379"
echo ""
echo "⏰ SSL certificates will be generated automatically by Let's Encrypt in 2-3 minutes."
echo "📋 Check certificate generation: docker logs nginx-letsencrypt -f"
