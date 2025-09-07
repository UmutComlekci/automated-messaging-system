# Automated Messaging System

An automatic message sending system built with Go Fiber framework, PostgreSQL, Redis and Temporal.

## Features

- 🚀 **Automatic Message Processing**: Sends 2 messages every 2 minutes
- 📊 **REST API**: Full CRUD operations with Swagger documentation
- 💾 **Redis Caching**: Caches sent messages with metadata
- 🐘 **PostgreSQL**: Reliable message storage with proper indexing
- 🐳 **Docker Support**: Complete containerization with docker-compose
- 📈 **Scalable**: Horizontal scaling with multiple pod support

## API Endpoints

### Scheduler Control
- `POST /api/v1/scheduler/start` - Start automatic message sending
- `POST /api/v1/scheduler/stop` - Stop automatic message sending

### Message Management
- `POST /api/v1/messages` - Create a new message
- `GET /api/v1/messages/sent` - Retrieve sent messages (paginated)

### Health & Documentation
- `GET /health` - Health check endpoint
- `GET /swagger` - Swagger API documentation

## Debugging

### Prerequisites
- Docker and Docker Compose

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/UmutComlekci/automated-messaging-system
   cd automated-messaging-system
   ```

2. **Start the services**
   ```bash
   docker-compose up -d --build
   ```
> NOTE: When the application is initialized using Docker Compose, the system automatically seeds the database with 5 default messages.

3. **Access the application**
   - API: http://localhost:8081/swagger
   - Temporal: http://localhost:8080

## Usage Examples

### Start Message Scheduling
```bash
curl -X POST http://localhost:8081/api/v1/scheduler/start
```

### Stop Message Scheduling
```bash
curl -X POST http://localhost:8081/api/v1/scheduler/stop
```

### Add a New Message
```bash
curl -X POST http://localhost:8081/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello! This is a test message.",
    "phone_number": "+905466011474"
  }'
```

### Get Sent Messages
```bash
curl "http://localhost:8081/api/v1/messages/sent?page=1&limit=10"
```
