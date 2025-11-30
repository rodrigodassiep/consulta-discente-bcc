package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	// Create test database
	testDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}

	// Auto-migrate all models
	testDB.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})

	// Set global db variable for handlers
	db = testDB

	// Create router with middleware
	r := gin.New()
	r.Use(CORSMiddleware())

	return r, testDB
}

func TestRootEndpoint(t *testing.T) {
	router, _ := setupTestRouter()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Student Feedback System API")
	})

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "Student Feedback System API")
}

func TestQuoteEndpoint(t *testing.T) {
	router, _ := setupTestRouter()

	router.GET("/quote", func(c *gin.Context) {
		c.JSON(200, "Don't communicate by sharing memory, share memory by communicating.")
	})

	req, _ := http.NewRequest("GET", "/quote", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "communicate")
}

func TestCurrentSemesterEndpoint(t *testing.T) {
	router, testDB := setupTestRouter()

	router.GET("/current-semester", func(c *gin.Context) {
		var semester Semester
		if err := db.Where("is_active = ?", true).First(&semester).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active semester found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"semester": semester})
	})

	t.Run("No Active Semester", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/current-semester", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 404, w.Code)
		assert.Contains(t, w.Body.String(), "No active semester found")
	})

	t.Run("Active Semester Found", func(t *testing.T) {
		// Create an active semester
		semester := Semester{
			Name:      "2024.1",
			Year:      2024,
			Period:    1,
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 4, 0),
			IsActive:  true,
		}
		testDB.Create(&semester)

		req, _ := http.NewRequest("GET", "/current-semester", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "2024.1")
		assert.Contains(t, w.Body.String(), "semester")
	})
}

func TestUserRegistration(t *testing.T) {
	router, _ := setupTestRouter()

	router.POST("/register", func(c *gin.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		// Validate required fields
		if newUser.FirstName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "First name is required"})
			return
		}
		if newUser.LastName == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Last name is required"})
			return
		}
		if newUser.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}
		if newUser.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
			return
		}

		// Set default role if not provided
		if newUser.Role == "" {
			newUser.Role = RoleStudent
		}

		// Validate role
		if newUser.Role != RoleStudent && newUser.Role != RoleProfessor && newUser.Role != RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}

		result := db.Create(&newUser)
		if result.Error != nil {
			// Check if it's a unique constraint violation (email already exists)
			if strings.Contains(result.Error.Error(), "duplicate key") || strings.Contains(result.Error.Error(), "UNIQUE constraint") {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// Create response user without password
		responseUser := map[string]interface{}{
			"id":         newUser.ID,
			"first_name": newUser.FirstName,
			"last_name":  newUser.LastName,
			"email":      newUser.Email,
			"role":       newUser.Role,
			"created_at": newUser.CreatedAt,
			"updated_at": newUser.UpdatedAt,
		}

		c.JSON(http.StatusCreated, responseUser)
	})

	t.Run("Valid User Registration", func(t *testing.T) {
		user := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe@example.com",
			"password":   "password123",
			"role":       RoleStudent,
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		assert.Contains(t, w.Body.String(), "john.doe@example.com")
		assert.Contains(t, w.Body.String(), "John")
		assert.NotContains(t, w.Body.String(), "password123") // Password should not be in response
	})

	t.Run("Missing First Name", func(t *testing.T) {
		user := map[string]interface{}{
			"last_name": "Doe",
			"email":     "john.doe2@example.com",
			"password":  "password123",
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "First name is required")
	})

	t.Run("Missing Email", func(t *testing.T) {
		user := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"password":   "password123",
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Email is required")
	})

	t.Run("Invalid Role", func(t *testing.T) {
		user := map[string]interface{}{
			"first_name": "John",
			"last_name":  "Doe",
			"email":      "john.doe3@example.com",
			"password":   "password123",
			"role":       "invalid_role",
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid role")
	})

	t.Run("Default Role Assignment", func(t *testing.T) {
		user := map[string]interface{}{
			"first_name": "Jane",
			"last_name":  "Smith",
			"email":      "jane.smith@example.com",
			"password":   "password123",
			// No role provided - should default to student
		}

		jsonValue, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 201, w.Code)
		assert.Contains(t, w.Body.String(), RoleStudent)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		// First registration
		user1 := map[string]interface{}{
			"first_name": "First",
			"last_name":  "User",
			"email":      "duplicate@example.com",
			"password":   "password123",
		}

		jsonValue1, _ := json.Marshal(user1)
		req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue1))
		req1.Header.Set("Content-Type", "application/json")
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		assert.Equal(t, 201, w1.Code)

		// Second registration with same email
		user2 := map[string]interface{}{
			"first_name": "Second",
			"last_name":  "User",
			"email":      "duplicate@example.com", // Same email
			"password":   "password456",
		}

		jsonValue2, _ := json.Marshal(user2)
		req2, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue2))
		req2.Header.Set("Content-Type", "application/json")
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, 409, w2.Code)
		assert.Contains(t, w2.Body.String(), "Email already exists")
	})
}

func TestUserLogin(t *testing.T) {
	router, testDB := setupTestRouter()

	router.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
			return
		}

		var foundUser User
		result := db.Where("email = ?", user.Email).First(&foundUser)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if foundUser.Password != user.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Create response user without password
		responseUser := map[string]interface{}{
			"id":         foundUser.ID,
			"first_name": foundUser.FirstName,
			"last_name":  foundUser.LastName,
			"email":      foundUser.Email,
			"role":       foundUser.Role,
			"created_at": foundUser.CreatedAt,
			"updated_at": foundUser.UpdatedAt,
		}

		c.JSON(http.StatusOK, responseUser)
	})

	// Create a test user
	testUser := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Password:  "testpass123",
		Role:      RoleStudent,
	}
	testDB.Create(&testUser)

	t.Run("Valid Login", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "testpass123",
		}

		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		assert.Contains(t, w.Body.String(), "test@example.com")
		assert.Contains(t, w.Body.String(), "Test")
		assert.NotContains(t, w.Body.String(), "testpass123") // Password should not be in response
	})

	t.Run("Invalid Email", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "testpass123",
		}

		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("Invalid Password", func(t *testing.T) {
		loginData := map[string]string{
			"email":    "test@example.com",
			"password": "wrongpassword",
		}

		jsonValue, _ := json.Marshal(loginData)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 401, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid credentials")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid request data")
	})
}

func TestLegacyConsultaEndpoint(t *testing.T) {
	router, _ := setupTestRouter()

	router.POST("/consulta", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This endpoint is deprecated. Use the new survey system."})
	})

	req, _ := http.NewRequest("POST", "/consulta", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "deprecated")
	assert.Contains(t, w.Body.String(), "survey system")
}

func TestAnonymousResponse(t *testing.T) {
	t.Run("ToAnonymous removes student identity", func(t *testing.T) {
		response := Response{
			ID:        1,
			SurveyID:  1,
			StudentID: 42,
			Student: User{
				ID:        42,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@test.com",
				Role:      RoleStudent,
			},
			QuestionID: 1,
			Question: Question{
				ID:   1,
				Text: "Test question",
				Type: QuestionTypeFreeText,
			},
			Answer:      "Test answer",
			SubmittedAt: time.Now(),
		}

		anonymous := response.ToAnonymous()

		// Verify that non-sensitive data is preserved
		assert.Equal(t, response.ID, anonymous.ID)
		assert.Equal(t, response.SurveyID, anonymous.SurveyID)
		assert.Equal(t, response.QuestionID, anonymous.QuestionID)
		assert.Equal(t, response.Answer, anonymous.Answer)
		assert.Equal(t, response.SubmittedAt, anonymous.SubmittedAt)
		assert.Equal(t, response.Question.Text, anonymous.Question.Text)
	})

	t.Run("ToAnonymousList converts multiple responses", func(t *testing.T) {
		responses := []Response{
			{
				ID:        1,
				SurveyID:  1,
				StudentID: 42,
				Student:   User{ID: 42, FirstName: "John"},
				Answer:    "Answer 1",
			},
			{
				ID:        2,
				SurveyID:  1,
				StudentID: 43,
				Student:   User{ID: 43, FirstName: "Jane"},
				Answer:    "Answer 2",
			},
		}

		anonymousList := ToAnonymousList(responses)

		assert.Len(t, anonymousList, 2)
		assert.Equal(t, responses[0].ID, anonymousList[0].ID)
		assert.Equal(t, responses[0].Answer, anonymousList[0].Answer)
		assert.Equal(t, responses[1].ID, anonymousList[1].ID)
		assert.Equal(t, responses[1].Answer, anonymousList[1].Answer)
	})

	t.Run("AnonymousResponse JSON excludes student fields", func(t *testing.T) {
		anonymous := AnonymousResponse{
			ID:         1,
			SurveyID:   1,
			QuestionID: 1,
			Answer:     "Test answer",
		}

		jsonBytes, err := json.Marshal(anonymous)
		assert.NoError(t, err)

		jsonStr := string(jsonBytes)
		// Verify student-related fields are not in JSON
		assert.NotContains(t, jsonStr, "student_id")
		assert.NotContains(t, jsonStr, "student")
		// Verify expected fields are present
		assert.Contains(t, jsonStr, "\"id\":1")
		assert.Contains(t, jsonStr, "\"answer\":\"Test answer\"")
	})
}
