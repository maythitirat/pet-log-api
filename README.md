# Pet Log API 🐾

A RESTful API for managing pet information, built with Go using clean architecture principles.

## 📁 Project Structure

```
pet-log-api/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/                  # Private application code
│   ├── config/               # Configuration management
│   ├── handler/              # HTTP handlers (controllers)
│   ├── middleware/           # HTTP middleware
│   ├── model/                # Data models and DTOs
│   ├── repository/           # Data access layer
│   ├── router/               # Route definitions
│   └── service/              # Business logic layer
├── pkg/                      # Public packages (can be imported by other projects)
│   ├── response/             # HTTP response helpers
│   └── validator/            # Request validation
├── migrations/               # Database migration files
├── .air.toml                 # Hot reload configuration
├── .env.example              # Environment variables template
├── .gitignore
├── docker-compose.yml        # Docker compose for local development
├── Dockerfile                # Production Docker image
├── go.mod                    # Go modules
├── Makefile                  # Development commands
└── README.md
```

## 🏗️ Architecture

This project follows **Clean Architecture** with clear separation of concerns:

```
Handler (HTTP) → Service (Business Logic) → Repository (Data Access) → Database
```

### Layers Explained:

1. **Handler Layer** (`internal/handler/`)
   - Handles HTTP requests and responses
   - Input validation
   - Maps HTTP requests to service calls

2. **Service Layer** (`internal/service/`)
   - Contains business logic
   - Orchestrates data flow between handlers and repositories
   - Independent of HTTP concerns

3. **Repository Layer** (`internal/repository/`)
   - Data access abstraction
   - Database queries
   - Can be easily swapped (e.g., PostgreSQL → MySQL)

4. **Model Layer** (`internal/model/`)
   - Domain entities
   - Request/Response DTOs
   - Data transformations

## 🚀 Getting Started

### Prerequisites

- Go 1.26+
- Docker & Docker Compose (for local development)
- PostgreSQL (if running without Docker)

### Quick Start

1. **Clone and setup:**
   ```bash
   cd pet-log-api
   cp .env.example .env
   ```

2. **Start with Docker:**
   ```bash
   make docker-up
   ```

3. **Or run locally:**
   ```bash
   # Start PostgreSQL
   docker-compose up -d postgres

   # Download dependencies
   make deps

   # Run the application
   make run
   ```

4. **Access the API:**
   - Health check: http://localhost:8080/health
   - API base URL: http://localhost:8080/api/v1

### Development with Hot Reload

```bash
# Install air (hot reload tool)
go install github.com/air-verse/air@latest

# Run with hot reload
make dev
```

## 📚 API Endpoints

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/ready` | Readiness check |

### Users
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/users` | Create a new user |
| GET | `/api/v1/users/{id}` | Get user by ID |
| PUT | `/api/v1/users/{id}` | Update user |
| DELETE | `/api/v1/users/{id}` | Delete user |

### Pets
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/pets` | Get all pets (with pagination) |
| POST | `/api/v1/pets` | Create a new pet |
| GET | `/api/v1/pets/{id}` | Get pet by ID |
| PUT | `/api/v1/pets/{id}` | Update pet |
| DELETE | `/api/v1/pets/{id}` | Delete pet |
| GET | `/api/v1/users/{userId}/pets` | Get pets by owner |

## 📝 API Examples

### Create a User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "password": "securepassword123"
  }'
```

### Create a Pet
```bash
curl -X POST http://localhost:8080/api/v1/pets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Buddy",
    "species": "dog",
    "breed": "Golden Retriever",
    "weight": 30.5,
    "owner_id": 1
  }'
```

### Get All Pets (with pagination)
```bash
curl "http://localhost:8080/api/v1/pets?page=1&page_size=10"
```

## 🛠️ Development

### Available Make Commands

```bash
make build          # Build the application
make run            # Run the application
make dev            # Run with hot reload
make test           # Run tests
make test-coverage  # Run tests with coverage
make lint           # Run linter
make fmt            # Format code
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations
make tools          # Install development tools
make help           # Show all commands
```

### Adding New Features

#### 1. Adding a New Entity (e.g., `log`)

**Step 1: Create the model** (`internal/model/log.go`)
```go
package model

type Log struct {
    ID        int64     `json:"id" db:"id"`
    PetID     int64     `json:"pet_id" db:"pet_id"`
    Type      string    `json:"type" db:"type"`
    Notes     string    `json:"notes" db:"notes"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

**Step 2: Create the repository** (`internal/repository/log_repository.go`)
```go
package repository

type LogRepository interface {
    Create(ctx context.Context, log *model.Log) error
    GetByPetID(ctx context.Context, petID int64) ([]*model.Log, error)
    // ... other methods
}
```

**Step 3: Create the service** (`internal/service/log_service.go`)
```go
package service

type LogService interface {
    Create(ctx context.Context, req *model.CreateLogRequest) (*model.LogResponse, error)
    GetByPetID(ctx context.Context, petID int64) ([]*model.LogResponse, error)
    // ... other methods
}
```

**Step 4: Create the handler** (`internal/handler/log_handler.go`)
```go
package handler

type LogHandler struct {
    service service.LogService
}

func (h *LogHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Handle HTTP request
}
```

**Step 5: Register routes** (`internal/router/router.go`)
```go
r.Route("/logs", func(r chi.Router) {
    r.Post("/", h.Log.Create)
    r.Get("/{id}", h.Log.GetByID)
})
```

**Step 6: Create migration** (`migrations/000002_add_logs.up.sql`)
```sql
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    pet_id BIGINT NOT NULL REFERENCES pets(id),
    type VARCHAR(50) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Code Style Guidelines

1. **Use meaningful variable names**
2. **Keep functions small and focused**
3. **Always handle errors explicitly**
4. **Use interfaces for dependencies (for testability)**
5. **Write tests for business logic**
6. **Follow Go naming conventions**

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./internal/service/... -run TestPetService
```

## 🔧 Configuration

Environment variables (see `.env.example`):

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | pet-log-api |
| `APP_ENV` | Environment (development/production) | development |
| `APP_PORT` | Server port | 8080 |
| `APP_DEBUG` | Debug mode | true |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database user | postgres |
| `DB_PASSWORD` | Database password | postgres |
| `DB_NAME` | Database name | pet_log |
| `DB_SSL_MODE` | SSL mode | disable |

## 📦 Dependencies

| Package | Purpose |
|---------|---------|
| [chi](https://github.com/go-chi/chi) | HTTP router |
| [sqlx](https://github.com/jmoiron/sqlx) | Database extensions |
| [zerolog](https://github.com/rs/zerolog) | Structured logging |
| [godotenv](https://github.com/joho/godotenv) | Environment variables |
| [pq](https://github.com/lib/pq) | PostgreSQL driver |

## 🚀 Deployment

### Building for Production

```bash
# Build binary
make build

# Or build Docker image
make docker-build
```

### Running in Production

```bash
docker run -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_HOST=your-db-host \
  -e DB_PASSWORD=your-secure-password \
  pet-log-api
```

## 📞 Support

If you have questions or need help, please:
1. Check the existing documentation
2. Look at the code examples in handlers
3. Ask your team lead or senior developer

---

Happy coding! 🎉
