# Coupon System API

A robust and concurrent coupon management system built with Go, featuring caching, transaction support, and Swagger documentation.

## Architecture

### High-Level Design
```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│    API      │     │   Service   │     │ Repository  │
│  (Gin)      │────▶│   Layer     │────▶│   (GORM)    │
└─────────────┘     └─────────────┘     └─────────────┘
       │                   │                   │
       │                   │                   │
       ▼                   ▼                   ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Request    │     │    Cache    │     │  SQLite DB  │
│ Validation  │     │   (LRU)     │     │             │
└─────────────┘     └─────────────┘     └─────────────┘
```

### Components
1. **API Layer (Gin)**
   - RESTful endpoints
   - Request validation
   - Response formatting
   - Swagger documentation

2. **Service Layer**
   - Business logic
   - Cache management
   - Concurrency control
   - Error handling

3. **Repository Layer (GORM)**
   - Database operations
   - Transaction management
   - Data persistence
   - SQLite integration

4. **Cache Layer (LRU)**
   - In-memory caching
   - Thread-safe operations
   - Configurable capacity

## Concurrency & Caching

### Concurrency Control
- **Read-Write Locks**: Using `sync.RWMutex` for thread-safe operations
- **Database Transactions**: Atomic operations for data consistency
- **Context Support**: Request cancellation and timeout handling
- **Cache Synchronization**: Thread-safe cache operations

### Caching Strategy
- **LRU Cache**: Least Recently Used eviction policy
- **Cache Keys**: Based on request parameters
- **Cache Invalidation**: Automatic on data updates
- **Thread Safety**: Protected by mutex locks

### Locking Mechanism
```go
// Read operations (multiple readers)
mu.RLock()
defer mu.RUnlock()

// Write operations (exclusive access)
mu.Lock()
defer mu.Unlock()
```

## Setup Instructions

### Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Git

### Local Development
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/coupon-system.git
   cd coupon-system
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run cmd/main.go
   ```

### Docker Deployment
1. Build and run with Docker Compose:
   ```bash
   docker-compose up --build
   ```

2. Access the API:
   - Swagger UI: http://localhost:8080/swagger/index.html
   - API Base URL: http://localhost:8080

## API Documentation

### Swagger UI
Access the interactive API documentation at:
http://localhost:8080/swagger/index.html

### Available Endpoints
1. **Get Applicable Coupons**
   - Method: POST
   - Path: `/coupons/applicable`
   - Description: Retrieves coupons applicable to the given cart

2. **Validate Coupon**
   - Method: POST
   - Path: `/coupons/validate`
   - Description: Validates a coupon against cart items

3. **Create Coupon**
   - Method: POST
   - Path: `/coupons/create`
   - Description: Creates a new coupon

## Data Persistence

### SQLite Database
- Location: `./data/coupon.db`
- Persisted through Docker volume
- Automatic schema migration
- Transaction support

### Cache Configuration
- Type: In-memory LRU Cache
- Default capacity: 100 items
- Thread-safe operations
- Automatic eviction

## Error Handling

### Error Types
1. **Validation Errors**
   - Invalid request format
   - Missing required fields
   - Invalid data types

2. **Business Logic Errors**
   - Invalid coupon code
   - Expired coupons
   - Insufficient order value

3. **System Errors**
   - Database connection issues
   - Cache operation failures
   - Internal server errors

## Testing

### Running Tests
```bash
go test ./...
```

### Test Coverage
```bash
go test -cover ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.