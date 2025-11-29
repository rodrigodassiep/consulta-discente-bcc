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

	t.Run("CORS Headers Set On Preflight Request", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		// Preflight OPTIONS request
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		req.Header.Set("Access-Control-Request-Method", "GET")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// gin-contrib/cors returns 204 for preflight
		assert.Equal(t, 204, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("CORS Headers Set On Regular Request With Origin", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://localhost:5173")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Equal(t, "http://localhost:5173", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	t.Run("CORS Rejects Unauthorized Origin", func(t *testing.T) {
		r := gin.New()
		r.Use(CORSMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "test"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", "http://malicious-site.com")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// gin-contrib/cors returns 403 for unauthorized origins
		assert.Equal(t, 403, w.Code)
		assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate User model
	testDB.AutoMigrate(&User{})

	// Save and restore original db
	originalDB := db
	db = testDB
	defer func() { db = originalDB }()

	// Set JWT secret for testing
	os.Setenv("JWT_SECRET", "test_secret_key")
	defer os.Unsetenv("JWT_SECRET")

	// Create test users
	student := User{
		FirstName: "John",
		LastName:  "Student",
		Email:     "student@test.com",
		Password:  "password123",
		Role:      RoleStudent,
	}
	testDB.Create(&student)

	professor := User{
		FirstName: "Jane",
		LastName:  "Professor",
		Email:     "professor@test.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	testDB.Create(&professor)

	admin := User{
		FirstName: "Admin",
		LastName:  "User",
		Email:     "admin@test.com",
		Password:  "password123",
		Role:      RoleAdmin,
	}
	testDB.Create(&admin)

	// Generate JWT tokens for test users
	studentToken, _ := GenerateJWT(student.ID, student.Role)
	professorToken, _ := GenerateJWT(professor.ID, professor.Role)
	adminToken, _ := GenerateJWT(admin.ID, admin.Role)

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

	t.Run("Invalid Authorization Header Format", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
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
		req.Header.Set("Authorization", "Bearer invalid_token_here")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid or expired token")
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

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
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

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+professorToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Unauthorized Role Access", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleAdmin)) // Require admin role

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken) // Student trying to access admin endpoint
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Multiple Allowed Roles - Student Allowed", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor)) // Allow both student and professor

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+studentToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Multiple Allowed Roles - Professor Allowed", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+professorToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})

	t.Run("Multiple Allowed Roles - Admin Denied", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor)) // Only student and professor allowed

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Admin Access to Admin Endpoint", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleAdmin))

		r.GET("/test", func(c *gin.Context) {
			currentUser, exists := c.Get("currentUser")
			assert.True(t, exists)
			user := currentUser.(User)
			assert.Equal(t, admin.ID, user.ID)
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})
}
