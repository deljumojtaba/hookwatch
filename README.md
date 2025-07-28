# 🔗 HookWatch

A simple, lightweight webhook monitoring and testing service built with Go and Docker. HookWatch allows you to receive, store, monitor, and replay webhooks with a clean web interface.

## 📋 What is HookWatch?

HookWatch is a webhook testing and monitoring tool that provides:

- **Webhook Reception**: Accept webhooks on custom endpoints with any HTTP method
- **Real-time Monitoring**: View incoming webhooks in real-time with full request details
- **Data Storage**: Store webhook logs in MongoDB with complete request information
- **Webhook Replay**: Replay captured webhooks to external endpoints for testing
- **Web Dashboard**: Clean, responsive web interface for monitoring and management
- **Multi-format Support**: Handle JSON, form data, query parameters, and raw payloads

## 🏗️ Architecture

### Backend (Go/Gin)

- **Framework**: Gin web framework for high-performance HTTP handling
- **Database**: MongoDB for webhook log storage
- **Cache**: Redis ready for future caching needs
- **API**: RESTful endpoints for webhook operations

### Frontend (Vanilla JS/HTML/CSS)

- Simple, responsive web dashboard
- Real-time webhook log viewing
- Webhook testing interface
- Clean, modern UI design

### Infrastructure (Docker)

- Fully containerized with Docker Compose
- MongoDB and Redis services
- Nginx for static file serving
- Production-ready configuration

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose
- Git

### 1. Clone the Repository

```bash
git clone https://github.com/sajjadgozal/hookwatch.git
cd hookwatch
```

### 2. Create Environment File

Create a `.env` file in the root directory:

```bash
# Database Configuration
MONGO_INITDB_ROOT_USERNAME=hookwatch
MONGO_INITDB_ROOT_PASSWORD=your_secure_password
MONGO_INITDB_PORT=27017

# Redis Configuration
REDIS_PORT=6379

# Application Ports
HOOKWATCH_PORT=3000
WEB_UI_PORT=8080
```

### 3. Start the Services

```bash
docker-compose up -d
```

This will start:

- **HookWatch API**: http://localhost:3000
- **Web Dashboard**: http://localhost:8080
- **MongoDB**: localhost:27017
- **Redis**: localhost:6379

### 4. Verify Installation

Check the health endpoint:

```bash
curl http://localhost:3000/health
```

Expected response:

```json
{
  "status": "healthy",
  "service": "hookwatch",
  "message": "Service is running"
}
```

## 📡 API Endpoints

### Webhook Reception

Accept webhooks on any endpoint ID:

```
{METHOD} /webhooks/{endpointId}/receive
```

**Examples:**

```bash
# POST with JSON payload
curl -X POST http://localhost:3000/webhooks/my-app/receive \
  -H "Content-Type: application/json" \
  -d '{"event": "user.created", "user_id": 123}'

# GET with query parameters
curl -X GET "http://localhost:3000/webhooks/my-app/receive?event=test&id=456"

# PUT with custom headers
curl -X PUT http://localhost:3000/webhooks/payment-service/receive \
  -H "Content-Type: application/json" \
  -H "X-Signature: abc123" \
  -d '{"transaction_id": "tx_789", "amount": 100.00}'
```

### Webhook Management

```bash
# Get webhook logs for an endpoint
GET /webhooks/{endpointId}/logs?limit=50

# Clear webhook logs for an endpoint
DELETE /webhooks/{endpointId}/logs

# Replay a specific webhook
POST /webhooks/replay/{webhookLogId}
{
  "target_url": "https://example.com/webhook",
  "timeout": 30
}
```

### Health Check

```bash
GET /health
```

## 🔧 Configuration

### Environment Variables

| Variable                     | Description               | Default                     |
| ---------------------------- | ------------------------- | --------------------------- |
| `PORT`                       | HookWatch API port        | `3000`                      |
| `MONGO_URI`                  | MongoDB connection string | `mongodb://localhost:27017` |
| `REDIS_URI`                  | Redis connection string   | `redis://localhost:6379`    |
| `MONGO_INITDB_ROOT_USERNAME` | MongoDB root username     | -                           |
| `MONGO_INITDB_ROOT_PASSWORD` | MongoDB root password     | -                           |
| `REDIS_PORT`                 | Redis port mapping        | `6379`                      |
| `HOOKWATCH_PORT`             | HookWatch port mapping    | `3000`                      |
| `WEB_UI_PORT`                | Web UI port mapping       | `8080`                      |

### Docker Compose Services

- **hookwatch**: Main Go application
- **mongodb**: MongoDB database with persistent storage
- **redis**: Redis cache (ready for future features)
- **web-ui**: Nginx serving the web dashboard

## 🎯 Use Cases

### 1. Webhook Development & Testing

- Test webhook integrations during development
- Inspect webhook payloads and headers
- Verify webhook delivery and format

### 2. Webhook Debugging

- Capture failed webhooks for analysis
- Replay webhooks to test fixes
- Monitor webhook delivery patterns

### 3. Integration Testing

- Test webhook endpoints before production
- Validate webhook security headers
- Test different payload formats

### 4. Webhook Monitoring

- Monitor production webhook traffic
- Track webhook delivery rates
- Store webhook history for analysis

## 📊 Data Models

### WebhookLog

```go
type WebhookLog struct {
    ID          primitive.ObjectID `json:"id"`
    EndpointID  string            `json:"endpoint_id"`
    Method      string            `json:"method"`
    Headers     map[string]string `json:"headers"`
    Body        interface{}       `json:"body"`
    IPAddress   string            `json:"ip_address"`
    UserAgent   string            `json:"user_agent"`
    Status      string            `json:"status"`
    ProcessedAt *time.Time        `json:"processed_at"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
}
```

### Features Captured Per Webhook

- Complete HTTP method and headers
- Full request body (JSON, form data, raw text)
- Client IP address and User-Agent
- Processing timestamps and status
- Endpoint identification

## 🛠️ Development

### Local Development Setup

```bash
# Clone the repository
git clone https://github.com/sajjadgozal/hookwatch.git
cd hookwatch

# Start dependencies only
docker-compose up -d mongodb redis web-ui

# Run the Go application locally
cd hookwatch
go mod download
go run main.go
```

### Project Structure

```
hookwatch/
├── docker-compose.yaml       # Docker services configuration
├── .env                      # Environment variables
├── hookwatch/               # Go application
│   ├── main.go             # Application entry point
│   ├── Dockerfile          # Go app container
│   ├── go.mod              # Go modules
│   ├── config/             # Configuration management
│   ├── db/                 # Database connections
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   └── services/           # Business logic
├── web/                    # Frontend dashboard
│   ├── index.html         # Main dashboard
│   ├── styles.css         # Styling
│   └── script.js          # JavaScript functionality
└── init/                   # Database initialization
    ├── mongo.js           # MongoDB setup
    └── redis.conf         # Redis configuration
```

### Building and Deployment

```bash
# Build all services
docker-compose build

# Deploy to production
docker-compose -f docker-compose.yaml up -d

# View logs
docker-compose logs -f hookwatch
```

## 🔒 Security Considerations

- Run with non-root user in production
- Use strong MongoDB credentials
- Implement rate limiting for production use
- Consider adding webhook signature verification
- Use HTTPS in production environments
- Regularly backup MongoDB data

## 📈 Future Enhancements

- [ ] Webhook signature verification
- [ ] Rate limiting and throttling
- [ ] Webhook filtering and routing
- [ ] Email/Slack notifications
- [ ] Webhook transformation rules
- [ ] API authentication
- [ ] Metrics and analytics dashboard
- [ ] Export webhook logs
- [ ] Webhook retry logic

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Port Conflicts

If you encounter port conflicts, modify the ports in your `.env` file:

```bash
HOOKWATCH_PORT=3001
WEB_UI_PORT=8081
MONGO_INITDB_PORT=27018
REDIS_PORT=6380
```

### Database Connection Issues

1. Ensure MongoDB container is running: `docker-compose ps`
2. Check MongoDB logs: `docker-compose logs mongodb`
3. Verify connection string in environment variables

### Cannot Access Web Dashboard

1. Verify web-ui container is running
2. Check port mapping in docker-compose.yaml
3. Try accessing http://127.0.0.1:8080 instead of localhost

### Webhook Not Being Received

1. Check HookWatch API is running: `curl http://localhost:3000/health`
2. Verify endpoint URL format: `/webhooks/{endpointId}/receive`
3. Check Docker container logs: `docker-compose logs hookwatch`

---

**Built with ❤️ using Go, MongoDB, Redis, and Docker**
