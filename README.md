# Fansly API

A REST API wrapper around the `fansly-scraper` functionality, providing structured endpoints for interacting with Fansly content.

## 📋 Project Structure

```
fansly-api/
├── cmd/               # Main application entry points
│   └── api/           # API server entry point
├── internal/          # Private application code
│   ├── api/           # HTTP handlers and routing
│   ├── config/        # Configuration management
│   ├── service/       # Business logic and scraper integration
│   └── storage/       # Data storage interfaces and implementations
└── pkg/              # Reusable packages
    ├── auth/         # Authentication utilities
    ├── logger/       # Logging utilities
    └── utils/        # Shared utilities
```

## 🚀 Features

### Authentication & Authorization
- [x] JWT-based authentication
- [ ] User registration (if needed)
- [ ] Role-based access control

### Creator Management
- [x] List creators (mock data)
- [ ] Get creator details
- [ ] Search creators
- [ ] Follow/unfollow creators

### Content Management
- [ ] List creator content
- [ ] Filter content by type (images, videos, etc.)
- [ ] Search within creator content
- [ ] Download media

### Monitoring
- [ ] Set up monitoring for creators
- [ ] Get monitoring status
- [ ] Configure monitoring preferences
- [ ] Webhook/notification system

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Docker (for database, if needed)
- Fansly account credentials

### Setup
1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Copy `.env.example` to `.env` and configure
4. Run the server: `go run cmd/api/main.go`

### Testing
- Run unit tests: `go test ./...`
- Integration tests: (TBD)
- E2E tests: (TBD)

## 📚 API Documentation

### Authentication
- `POST /api/v1/auth/initiate` - Start authentication
- `POST /api/v1/auth/complete` - Complete authentication
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/logout` - Invalidate token

### Creators
- `GET /api/v1/creators` - List creators
- `GET /api/v1/creators/{id}` - Get creator details
- `GET /api/v1/creators/{id}/content` - Get creator content
- `POST /api/v1/creators/{id}/follow` - Follow a creator
- `DELETE /api/v1/creators/{id}/follow` - Unfollow a creator

## 📅 Roadmap

### Phase 1: Core Functionality
1. Scraper Integration
   - [ ] Create service layer for `fansly-scraper`
   - [ ] Implement authentication flow
   - [ ] Add error handling

2. API Endpoints
   - [ ] Implement creator listing with real data
   - [ ] Add content retrieval endpoints

3. Configuration
   - [ ] Set up configuration management
   - [ ] Add rate limiting
   - [ ] Configure logging

### Phase 2: Advanced Features
- Monitoring system
- User management
- API documentation

### Phase 3: Production Readiness
- Performance optimizations
- Security enhancements
- Deployment setup

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
