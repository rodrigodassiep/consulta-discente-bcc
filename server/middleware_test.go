package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
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
		assert.Equal(t, "Content-Type, Authorization, X-User-ID", w.Header().Get("Access-Control-Allow-Headers"))
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
}

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Auto-migrate User model
	testDB.AutoMigrate(&User{})

	// Set global db variable for middleware
	originalDB := db
	defer func() { db = originalDB }()

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
		LastName:  "Admin",
		Email:     "admin@test.com",
		Password:  "password123",
		Role:      RoleAdmin,
	}
	testDB.Create(&admin)

	// Temporarily set the global db variable
	originalGlobalDB := db
	db = testDB
	defer func() { db = originalGlobalDB }()

	t.Run("Missing User ID Header", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "User ID header required")
	})

	t.Run("Invalid User ID Header", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-User-ID", "invalid")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid user ID")
	})

	t.Run("User Not Found", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent))

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-User-ID", "999") // Non-existent user ID
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

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-User-ID", strconv.Itoa(int(student.ID)))
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
		req.Header.Set("X-User-ID", strconv.Itoa(int(professor.ID)))
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
		req.Header.Set("X-User-ID", strconv.Itoa(int(student.ID))) // Student trying to access admin endpoint
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 403, w.Code)
		assert.Contains(t, w.Body.String(), "Insufficient permissions")
	})

	t.Run("Multiple Allowed Roles", func(t *testing.T) {
		r := gin.New()
		r.Use(RequireRole(RoleStudent, RoleProfessor)) // Allow both student and professor

		r.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "authorized"})
		})

		// Test with student
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("X-User-ID", strconv.Itoa(int(student.ID)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")

		// Test with professor
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.Header.Set("X-User-ID", strconv.Itoa(int(professor.ID)))
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		assert.Equal(t, 200, w2.Code)
		assert.Contains(t, w2.Body.String(), "authorized")

		// Test with admin (should be denied)
		req3, _ := http.NewRequest("GET", "/test", nil)
		req3.Header.Set("X-User-ID", strconv.Itoa(int(admin.ID)))
		w3 := httptest.NewRecorder()
		r.ServeHTTP(w3, req3)

		assert.Equal(t, 403, w3.Code)
		assert.Contains(t, w3.Body.String(), "Insufficient permissions")
	})

	t.Run("Admin Access to All Roles", func(t *testing.T) {
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
		req.Header.Set("X-User-ID", strconv.Itoa(int(admin.ID)))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "authorized")
	})
}
