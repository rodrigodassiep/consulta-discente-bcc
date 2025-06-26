# Testing Documentation

This document explains the comprehensive testing setup for the Student Feedback System Go backend.

## Test Structure

The project includes several types of tests:

### 1. Unit Tests

#### Model Tests (`models_test.go`)
- Tests all database models and their validation
- Tests GORM relationships and constraints
- Tests constants and their values
- Uses in-memory SQLite database for isolation

**Coverage:**
- `User` model validation and uniqueness constraints
- `Subject` model and professor relationships
- `Semester` model and date validation
- `StudentEnrollment` model and foreign key relationships
- `Survey` model creation and validation
- `Question` model with different question types
- `Response` model with all required relationships

#### Middleware Tests (`middleware_test.go`)
- Tests CORS middleware functionality
- Tests role-based authorization middleware
- Tests user authentication and authorization flows
- Tests error handling for invalid requests

**Coverage:**
- CORS headers are set correctly
- OPTIONS requests are handled properly
- User authentication via X-User-ID header
- Role validation for different user types
- Authorization for multiple roles
- Error responses for missing/invalid credentials

#### API Tests (`api_test.go`)
- Tests core API endpoints
- Tests user registration and login
- Tests request validation and error handling
- Tests JSON response formats

**Coverage:**
- Root endpoint functionality
- User registration with validation
- User login authentication
- Error handling for invalid requests
- Response format validation

#### Database Seeding Tests (`seed_test.go`)
- Tests the database seeding functionality
- Verifies data consistency and relationships
- Tests seeding behavior with existing data
- Validates seeded data integrity

**Coverage:**
- Fresh database seeding
- Skip seeding when data exists
- Specific seeded data verification
- Data consistency across relationships
- Date constraints and validation

## Running Tests

### Prerequisites

Make sure you have the required dependencies installed:

```bash
go mod tidy
```

### Basic Test Commands

#### Run All Tests
```bash
go test ./...
```

#### Run Tests with Verbose Output
```bash
go test -v ./...
```

#### Run Specific Test File
```bash
go test -v models_test.go
go test -v middleware_test.go
go test -v api_test.go
go test -v seed_test.go
```

#### Run Specific Test Function
```bash
go test -v -run TestUserModel
go test -v -run TestCORSMiddleware
go test -v -run TestUserRegistration
```

### Coverage Analysis

#### Generate Coverage Report
```bash
go test ./... -coverprofile=coverage.out
```

#### View Coverage in Terminal
```bash
go tool cover -func=coverage.out
```

#### Generate HTML Coverage Report
```bash
go tool cover -html=coverage.out -o coverage.html
```

### Race Condition Detection
```bash
go test -race ./...
```

### Benchmark Tests
```bash
go test -bench=. -benchmem ./...
```

### Short Tests (Skip Long-Running Tests)
```bash
go test -short ./...
```

## Test Utilities

### `setupTestDB()` Function
- Creates an in-memory SQLite database for each test
- Auto-migrates all models
- Provides clean isolated test environment
- Used across all model tests

### `setupTestRouter()` Function
- Creates a test Gin router with middleware
- Sets up test database connection
- Configures CORS middleware
- Used for API endpoint testing

## Test Patterns

### Database Tests
```go
func TestModelName(t *testing.T) {
    db := setupTestDB()
    
    t.Run("Test Case Name", func(t *testing.T) {
        // Test implementation
        assert.NoError(t, err)
        assert.Equal(t, expected, actual)
    })
}
```

### API Tests
```go
func TestEndpointName(t *testing.T) {
    router, testDB := setupTestRouter()
    
    // Setup test data if needed
    
    t.Run("Test Case Name", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/endpoint", nil)
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, 200, w.Code)
        assert.Contains(t, w.Body.String(), "expected content")
    })
}
```

### Middleware Tests
```go
func TestMiddlewareName(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    t.Run("Test Case Name", func(t *testing.T) {
        r := gin.New()
        r.Use(MiddlewareFunction())
        
        // Test middleware behavior
    })
}
```

## Dependencies

The testing setup uses the following dependencies:

- `github.com/stretchr/testify/assert` - Assertion library
- `gorm.io/driver/sqlite` - In-memory SQLite for testing
- `github.com/gin-gonic/gin` - Web framework testing utilities
- `net/http/httptest` - HTTP testing utilities

## Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Clean State**: Use fresh database instances for each test
3. **Descriptive Names**: Test function names should clearly describe what they test
4. **Multiple Scenarios**: Test both success and failure cases
5. **Error Handling**: Verify proper error messages and status codes
6. **Data Validation**: Test all validation rules and constraints

## Coverage Goals

- **Model Tests**: 100% coverage of all model methods and validations
- **Middleware Tests**: 100% coverage of all middleware functions
- **API Tests**: 90%+ coverage of all endpoints
- **Integration Tests**: Cover complete user workflows

## Continuous Integration

For CI/CD pipelines, use:

```bash
# Run all tests with coverage
go test ./... -coverprofile=coverage.out -covermode=atomic

# Check minimum coverage (adjust percentage as needed)
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | awk '{if ($1 < 80) exit 1}'
```

## Troubleshooting

### Common Issues

1. **Database Connection Errors**: Ensure SQLite driver is properly imported
2. **Import Cycle**: Keep test files in the same package to avoid cycles
3. **Global Variables**: Reset global variables between tests if needed
4. **Concurrency Issues**: Use `-race` flag to detect race conditions

### Debug Tips

- Use `-v` flag for verbose output
- Add `t.Log()` statements for debugging
- Use `t.Skip()` to temporarily skip failing tests
- Use `t.Parallel()` for parallel test execution (when safe)
