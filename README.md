# 🚀 TaskFlow API

> A production-ready RESTful API for task management built with Go
> ## 🌐 Live Demo

**API Base URL:** https://taskflow-api.railway.app

Try it out:
- Health Check: https://taskflow-api.railway.app/health
- API Docs: See below

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Project Structure](#-project-structure)
- [Development](#-development)
- [Testing](#-testing)
- [Deployment](#-deployment)
- [Author](#-author)

## ✨ Features

- **User Authentication**
  - Secure registration with password hashing (bcrypt)
  - JWT-based authentication
  - Protected routes with middleware

- **Task Management**
  - Create, read, update, delete tasks (CRUD)
  - Task prioritization (high, medium, low)
  - Task categorization (personal, work, urgent)
  - Due date tracking
  - Completion status toggle
  - Advanced filtering and search

- **Analytics**
  - Task statistics (total, completed, pending)
  - Priority-based metrics
  - Overdue task tracking

- **Production Ready**
  - Clean architecture with separation of concerns
  - Comprehensive error handling
  - Database connection pooling
  - Docker containerization
  - Health check endpoint

## 🛠️ Tech Stack

| Technology | Purpose |
|------------|---------|
| **Go 1.24** | Backend programming language |
| **Gin** | HTTP web framework |
| **PostgreSQL 15** | Relational database |
| **JWT** | Authentication tokens |
| **Docker** | Containerization |
| **bcrypt** | Password hashing |

## 🏗️ Architecture
```
┌─────────────┐
│   Client    │
│ (Postman/   │
│  Browser)   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────┐
│      API Layer (Gin Router)     │
│  ┌────────────────────────────┐ │
│  │    Middleware (JWT Auth)   │ │
│  └────────────────────────────┘ │
│  ┌────────────────────────────┐ │
│  │   Handlers (Controllers)   │ │
│  │  - Auth Handler            │ │
│  │  - Task Handler            │ │
│  └────────────────────────────┘ │
│  ┌────────────────────────────┐ │
│  │   Business Logic Layer     │ │
│  │  - Validation              │ │
│  │  - Data Processing         │ │
│  └────────────────────────────┘ │
└────────────┬────────────────────┘
             │
             ▼
    ┌────────────────┐
    │   PostgreSQL   │
    │    Database    │
    └────────────────┘
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker & Docker Compose
- PostgreSQL 15 (if running locally)
- Git

### Installation

#### Option 1: Using Docker (Recommended)

1. **Clone the repository**
```bash
git clone https://github.com/YOUR_USERNAME/taskflow-api.git
cd taskflow-api
```

2. **Build and run with Docker Compose**
```bash
docker-compose build
docker-compose up -d
```

3. **Run database migrations**
```bash
docker exec -i taskflow-postgres psql -U postgres -d taskflow < migrations/001_init.sql
```

4. **Verify the API is running**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "message": "TaskFlow API is running",
  "status": "ok"
}
```

#### Option 2: Local Development

1. **Clone the repository**
```bash
git clone https://github.com/YOUR_USERNAME/taskflow-api.git
cd taskflow-api
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup PostgreSQL database**
```bash
createdb taskflow
psql taskflow < migrations/001_init.sql
```

4. **Create `.env` file**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. **Run the application**
```bash
go run cmd/api/main.go
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2026-01-30T10:00:00Z",
    "updated_at": "2026-01-30T10:00:00Z"
  }
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

#### Get Profile
```http
GET /auth/profile
Authorization: Bearer <token>
```

### Task Endpoints

All task endpoints require authentication (Bearer token).

#### Create Task
```http
POST /tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Complete project documentation",
  "description": "Write comprehensive README and API docs",
  "priority": "high",
  "category": "work",
  "due_date": "2026-02-01T23:59:59Z"
}
```

#### Get All Tasks
```http
GET /tasks
Authorization: Bearer <token>

# With filters
GET /tasks?priority=high&category=work&is_completed=false&search=documentation
```

#### Get Task by ID
```http
GET /tasks/:id
Authorization: Bearer <token>
```

#### Update Task
```http
PUT /tasks/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated title",
  "priority": "medium",
  "is_completed": true
}
```

#### Delete Task
```http
DELETE /tasks/:id
Authorization: Bearer <token>
```

#### Toggle Task Completion
```http
PATCH /tasks/:id/complete
Authorization: Bearer <token>
```

#### Get Task Statistics
```http
GET /tasks/stats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Statistics retrieved successfully",
  "data": {
    "total": 10,
    "completed": 5,
    "pending": 5,
    "high_priority": 3,
    "overdue": 1
  }
}
```

### Error Responses

All error responses follow this format:
```json
{
  "success": false,
  "error": "Error message here"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `409` - Conflict (e.g., email already exists)
- `500` - Internal Server Error

## 📁 Project Structure
```
taskflow-api/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── database/
│   │   └── database.go            # Database connection
│   ├── handlers/
│   │   ├── auth_handler.go        # Authentication endpoints
│   │   └── task_handler.go        # Task management endpoints
│   ├── middleware/
│   │   └── auth_middleware.go     # JWT authentication middleware
│   ├── models/
│   │   ├── user.go                # User data structures
│   │   └── task.go                # Task data structures
│   └── utils/
│       ├── password.go            # Password hashing utilities
│       └── response.go            # Standard response helpers
├── migrations/
│   └── 001_init.sql               # Database schema
├── .env.example                    # Environment variables template
├── .gitignore                      # Git ignore rules
├── docker-compose.yml              # Docker Compose configuration
├── Dockerfile                      # Docker build instructions
├── go.mod                          # Go module dependencies
├── go.sum                          # Go module checksums
├── Makefile                        # Development commands
└── README.md                       # This file
```

## 🔧 Development

### Useful Commands
```bash
# Build the application
go build -o bin/taskflow cmd/api/main.go

# Run tests
go test -v ./...

# Format code
go fmt ./...

# Check for issues
go vet ./...

# Docker commands
docker-compose up -d           # Start services
docker-compose down            # Stop services
docker-compose logs -f api     # View API logs
docker-compose restart         # Restart services
docker-compose down -v         # Remove all data
```

### Environment Variables
```env
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/taskflow?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
```

## 🧪 Testing

### Manual Testing with cURL

**Health Check:**
```bash
curl http://localhost:8080/health
```

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**Create Task:**
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"Test Task","priority":"high","category":"work"}'
```

### Testing with Postman

Import the API collection:
1. Open Postman
2. Create new collection: "TaskFlow API"
3. Add requests for all endpoints
4. Set up environment variables for the token

## 🚢 Deployment

### Docker Deployment

The application is containerized and ready for deployment to any Docker-compatible platform:

- **Docker Hub**: Build and push the image
- **AWS ECS/EKS**: Deploy with Fargate or Kubernetes
- **Google Cloud Run**: Serverless container deployment
- **Azure Container Instances**: Quick container deployment
- **Railway.app**: Easy one-click deployment
- **Render.com**: Free tier available

### Example: Deploy to Railway

1. Push code to GitHub
2. Connect Railway to your repository
3. Add PostgreSQL database
4. Set environment variables
5. Deploy!

## 🎯 Future Enhancements

- [ ] Add unit and integration tests
- [ ] Implement task sharing between users
- [ ] Add task comments and attachments
- [ ] Email notifications for due dates
- [ ] Task templates
- [ ] Export tasks to CSV/PDF
- [ ] RESTful API versioning
- [ ] Rate limiting
- [ ] Comprehensive logging with structured logs
- [ ] Metrics and monitoring (Prometheus)
- [ ] CI/CD pipeline with GitHub Actions

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

**Dani Abizar Ahmad**

- GitHub: [@daniabizar](https://github.com/daniabizar)
- Email: daniabizar.44@gmail.com
- LinkedIn: https://linkedin/daniabizar

---

## 🙏 Acknowledgments

This project was developed as part of my application to the **GoTo Engineering Bootcamp 2026**. It demonstrates:

- ✅ Clean code architecture and best practices
- ✅ RESTful API design principles
- ✅ Database design and optimization
- ✅ Authentication and security
- ✅ Docker containerization
- ✅ Production-ready development practices

**Built with ❤️ for GoTo Group**

---

<p align="center">
  <i>If you found this project helpful, please consider giving it a ⭐️</i>
</p>
