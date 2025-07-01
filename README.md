# Custom HTTP Server with Authentication & Database

A production-ready HTTP server built with Go featuring JWT authentication, PostgreSQL database integration, and comprehensive middleware system.

## 🚀 Features

- **Custom HTTP Server**: Built with gorilla/mux for robust routing
- **JWT Authentication**: Secure token-based authentication with bcrypt password hashing
- **Database Integration**: PostgreSQL with connection pooling and transaction support
- **Middleware System**: CORS, Logging, Security headers, Rate limiting, Authentication
- **ORM-like Repository Pattern**: Clean data access layer with user management
- **Unit Tests**: Comprehensive test coverage for all packages
- **Docker Support**: Ready for containerization with Docker Compose
- **Graceful Shutdown**: Clean server termination with context timeout
- **Health Checks**: Database and server health monitoring
- **Security Features**: Rate limiting, input validation, security headers

## 📁 Project Structure

```
.
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── config/
│   └── config.go          # Configuration management
├── server/
│   └── server.go          # HTTP server implementation
├── middleware/
│   └── middleware.go      # Middleware functions
├── auth/
│   └── auth.go           # JWT authentication logic
├── database/
│   └── database.go       # Database connection and ORM
├── models/
│   └── user.go           # User model and repository
├── handlers/
│   ├── auth.go           # Authentication handlers
│   ├── user.go           # User management handlers
│   └── health.go         # Health check handlers
├── tests/                 # Unit tests
│   ├── auth_test.go      # Authentication tests
│   └── handlers_test.go  # Handler tests
├── env.example           # Environment variables template
├── docker-compose.yml    # Docker services configuration
├── Dockerfile            # Docker build configuration
└── README.md             # This file
```

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker & Docker Compose (optional)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd ServerGo
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Set up PostgreSQL database**
   ```bash
   # Create database and user
   createdb myapp
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

### Docker Deployment

1. **Build and run with Docker Compose**
   ```bash
   docker-compose up --build
   ```

2. **Or build manually**
   ```bash
   docker build -t custom-http-server .
   docker run -p 8080:8080 custom-http-server
   ```

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Login User
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Protected Endpoints

#### Get User Profile
```http
GET /auth/me
Authorization: Bearer <jwt_token>
```

#### Update User Profile
```http
PUT /auth/me
Authorization: Bearer <jwt_token>
Content-Type: application/json

{
  "name": "Updated Name"
}
```

#### List Users (with pagination)
```http
GET /users?limit=10&offset=0
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "email": "user@example.com",
      "name": "John Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "limit": 10,
  "offset": 0,
  "total": 1
}
```

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "database": "connected",
  "version": "2.0.0",
  "uptime": "1h2m3s"
}
```

## 🧪 Testing

Run all tests:
```bash
go test ./...
```

Run specific test packages:
```bash
go test ./tests/...
go test ./auth/...
```

Run tests with coverage:
```bash
go test -cover ./...
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://user:password@localhost/dbname?sslmode=disable` |
| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key-change-this-in-production` |
| `LOG_LEVEL` | Logging level | `info` |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limit per IP | `100` |

### Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
```

## 🔒 Security Features

- **Password Hashing**: bcrypt with cost 14
- **JWT Tokens**: HMAC-SHA256 signed tokens with 24-hour expiration
- **Rate Limiting**: 100 requests per minute per IP address
- **Security Headers**: XSS protection, content type options, frame options
- **Input Validation**: Email format, password strength requirements
- **SQL Injection Protection**: Parameterized queries
- **CORS Support**: Configurable cross-origin resource sharing

## 🚀 Performance Optimizations

- **Connection Pooling**: Database connection pool with optimized settings
- **Goroutines**: Concurrent request handling
- **Middleware Chain**: Efficient middleware execution
- **Indexed Queries**: Database indexes for fast lookups
- **JSON Logging**: Structured logging for better performance
- **Graceful Shutdown**: Clean resource cleanup

## 📊 Monitoring & Health Checks

- **Health Endpoint**: `/health` for server and database status
- **Uptime Tracking**: Server uptime monitoring
- **Database Connectivity**: Real-time database connection status
- **Request Logging**: Detailed request/response logging with timing

## 🐳 Docker Support

The application includes full Docker support:

- **Multi-stage Build**: Optimized Docker image size
- **Docker Compose**: Complete development environment
- **Health Checks**: Database and service health monitoring
- **Volume Mounting**: Development code hot-reloading
- **Environment Variables**: Configurable via Docker environment

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions:
- Create an issue in the repository
- Check the API documentation
- Review the test examples

## 🔄 Version History

- **v2.0.0**: Production-ready release with comprehensive features
- **v1.0.0**: Initial release with basic functionality 