package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCORSMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CORS Headers Set Correctly", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("OPTIONS Request Handled", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("POST Request Gets CORS Headers", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.POST("/test", func(c *gin.Context) {
			c.JSON(201, gin.H{"message": "created"})
		})

		req, _ := http.NewRequest("POST", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRequireRoleWithJWT(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Save original env and restore after tests
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)
	os.Setenv("JWT_SECRET", "test-secret-for-middleware")

	// Setup test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate User model
	testDB.AutoMigrate(&User{})

	// Save original db and restore after tests
	originalDB := db
	db = testDB
	defer func() { db = originalDB }()

	// Create test users
	studentPassword, _ := HashPassword("password123")
	student := User{
		FirstName: "John",
		LastName:  "Student",
		Email:     "student@test.com",
		Password:  studentPassword,
		Role:      RoleStudent,
	}
	testDB.Create(&student)

	professorPassword, _ := HashPassword("password123")
	professor := User{
		FirstName: "Jane",
		LastName:  "Professor",
		Email:     "professor@test.com",
		Password:  professorPassword,
		Role:      RoleProfessor,
	}
	testDB.Create(&professor)

	adminPassword, _ := HashPassword("password123")
	admin := User{
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     "admin@test.com",
		Password:  adminPassword,
		Role:      RoleAdmin,
	}
	testDB.Create(&admin)

	t.Run("Missing Authorization Header", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Authorization header required")
	})

	t.Run("Invalid Authorization Header Format - No Bearer", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat token123")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("Invalid Authorization Header Format - Just Token", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "sometoken")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid authorization header format")
	})

	t.Run("Invalid JWT Token", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-jwt-token")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})

	t.Run("Valid Token But User Not Found", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		// Generate token for non-existent user
		token, _ := GenerateJWT(999, RoleStudent)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "User not found")
	})

	t.Run("Authorized Student Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			currentUser, exists := c.Get("currentUser")
			assert.True(t, exists)
			user := currentUser.(User)
			assert.Equal(t, student.ID, user.ID)
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(student.ID, student.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Authorized Professor Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			currentUser, exists := c.Get("currentUser")
			assert.True(t, exists)
			user := currentUser.(User)
			assert.Equal(t, professor.ID, user.ID)
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(professor.ID, professor.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Authorized Admin Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleAdmin))

		r.GET("/test", func(c *gin.Context) {
			currentUser, exists := c.Get("currentUser")
			assert.True(t, exists)
			user := currentUser.(User)
			assert.Equal(t, admin.ID, user.ID)
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(admin.ID, admin.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Unauthorized Role Access - Student to Admin Endpoint", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleAdmin))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(student.ID, student.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Unauthorized Role Access - Professor to Admin Endpoint", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleAdmin))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(professor.ID, professor.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Multiple Allowed Roles - Student Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(student.ID, student.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Multiple Allowed Roles - Professor Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(professor.ID, professor.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Multiple Allowed Roles - Admin Denied", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(admin.ID, admin.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Context Contains User ID", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("userID")
			assert.True(t, exists)
			assert.Equal(t, student.ID, userID.(uint))
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(student.ID, student.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Context Contains User Role", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			userRole, exists := c.Get("userRole")
			assert.True(t, exists)
			assert.Equal(t, RoleProfessor, userRole.(string))
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(professor.ID, professor.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Token With Wrong Secret Is Rejected", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		// Generate token with different secret
		os.Setenv("JWT_SECRET", "different-secret")
		token, _ := GenerateJWT(student.ID, student.Role)

		// Reset to original secret for validation
		os.Setenv("JWT_SECRET", "test-secret-for-middleware")

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
	})
}

func TestRequireRoleMiddlewareChain(t *testing.T) {
	gin.SetMode(gin.TestMode)

	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Setenv("JWT_SECRET", "")

	// Setup test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	testDB.AutoMigrate(&User{})

	originalDB := db
	db = testDB
	defer func() { db = originalDB }()

	// Create test user
	password, _ := HashPassword("password123")
	user := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@test.com",
		Password:  password,
		Role:      RoleStudent,
	}
	testDB.Create(&user)

	t.Run("Middleware Works With CORS", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		token, _ := GenerateJWT(user.ID, user.Role)

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("OPTIONS Request Bypasses Auth With CORS", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// OPTIONS should return 204 from CORS middleware, not 401 from auth
		assert.Equal(t, 204, w.Code)
	})
}
