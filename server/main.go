package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"rsc.io/quote"
)

// User roles constants
const (
	RoleStudent   = "student"
	RoleProfessor = "professor"
	RoleAdmin     = "admin"
)

// Question types constants
const (
	QuestionTypeNPS      = "nps"
	QuestionTypeFreeText = "free_text"
	QuestionTypeRating   = "rating"
	QuestionTypeChoice   = "multiple_choice"
)

// JWT Claims structure
type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Password hashing utilities
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT token utilities
func GenerateJWT(userID uint, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback_secret_key_change_in_production"
	}
	return token.SignedString([]byte(jwtSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "fallback_secret_key_change_in_production"
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}

// User model with proper role handling
type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	FirstName     string    `json:"first_name" gorm:"not null"`
	LastName      string    `json:"last_name" gorm:"not null"`
	Email         string    `json:"email" gorm:"uniqueIndex;not null"`
	Password      string    `json:"password" gorm:"not null"`
	Role          string    `json:"role" gorm:"not null;check:role IN ('student','professor','admin')"`
	RequestedRole string    `json:"requested_role" gorm:"not null;default:'student'"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Subject (course information)
type Subject struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Code        string    `json:"code" gorm:"uniqueIndex;not null"`
	Description string    `json:"description"`
	ProfessorID uint      `json:"professor_id" gorm:"not null"`
	Professor   User      `json:"professor" gorm:"foreignKey:ProfessorID;references:ID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Semester (academic periods)
type Semester struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"` // e.g., "2024.1", "2024.2"
	Year      int       `json:"year" gorm:"not null"`
	Period    int       `json:"period" gorm:"not null"` // 1 or 2
	StartDate time.Time `json:"start_date" gorm:"not null"`
	EndDate   time.Time `json:"end_date" gorm:"not null"`
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// StudentEnrollment (student-subject-semester relationships)
type StudentEnrollment struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	StudentID  uint      `json:"student_id" gorm:"not null"`
	Student    User      `json:"student" gorm:"foreignKey:StudentID;references:ID"`
	SubjectID  uint      `json:"subject_id" gorm:"not null"`
	Subject    Subject   `json:"subject" gorm:"foreignKey:SubjectID;references:ID"`
	SemesterID uint      `json:"semester_id" gorm:"not null"`
	Semester   Semester  `json:"semester" gorm:"foreignKey:SemesterID;references:ID"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Survey (feedback forms created by professors)
type Survey struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	SubjectID   uint       `json:"subject_id" gorm:"not null"`
	Subject     Subject    `json:"subject" gorm:"foreignKey:SubjectID;references:ID"`
	SemesterID  uint       `json:"semester_id" gorm:"not null"`
	Semester    Semester   `json:"semester" gorm:"foreignKey:SemesterID;references:ID"`
	ProfessorID uint       `json:"professor_id" gorm:"not null"`
	Professor   User       `json:"professor" gorm:"foreignKey:ProfessorID;references:ID"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	OpenDate    time.Time  `json:"open_date"`
	CloseDate   time.Time  `json:"close_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Questions   []Question `json:"questions" gorm:"foreignKey:SurveyID"`
}

// Question (individual questions with types)
type Question struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SurveyID  uint      `json:"survey_id" gorm:"not null"`
	Survey    Survey    `json:"survey" gorm:"foreignKey:SurveyID;references:ID"`
	Type      string    `json:"type" gorm:"not null;check:type IN ('nps','free_text','rating','multiple_choice')"`
	Text      string    `json:"text" gorm:"not null"`
	Required  bool      `json:"required" gorm:"default:false"`
	Order     int       `json:"order" gorm:"not null"`
	Options   string    `json:"options"` // JSON string for multiple choice options
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Response (student answers)
type Response struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	SurveyID    uint      `json:"survey_id" gorm:"not null"`
	Survey      Survey    `json:"survey" gorm:"foreignKey:SurveyID;references:ID"`
	StudentID   uint      `json:"student_id" gorm:"not null"`
	Student     User      `json:"student" gorm:"foreignKey:StudentID;references:ID"`
	QuestionID  uint      `json:"question_id" gorm:"not null"`
	Question    Question  `json:"question" gorm:"foreignKey:QuestionID;references:ID"`
	Answer      string    `json:"answer" gorm:"not null"`
	SubmittedAt time.Time `json:"submitted_at" gorm:"autoCreateTime"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Global database variable
var db *gorm.DB

// CORSMiddleware handles OPTIONS requests and sets CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		// Handle OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			log.Println("OPTIONS request received")
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// JWT-based role middleware
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := tokenParts[1]
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if user role is allowed
		allowed := false
		for _, role := range allowedRoles {
			if claims.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user data in context for later use
		var user User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("currentUser", user)
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Next()
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	var host, user, password, dbname, port string
	host = os.Getenv("DBHOST")
	user = os.Getenv("DBUSER")
	password = os.Getenv("DBPASSWORD")
	dbname = os.Getenv("DBNAME")
	port = os.Getenv("DBPORT")
	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"

	var err_sql error
	db, err_sql = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err_sql != nil {
		panic("failed to connect database")
	}

	// Auto-migrate all the new models
	log.Println("🔧 Running database migrations...")
	err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
	if err != nil {
		log.Printf("⚠️  Migration error: %v", err)
		log.Println("🔄 Attempting to reset database...")

		// Drop all tables and recreate them
		db.Migrator().DropTable(&Response{}, &Question{}, &Survey{}, &StudentEnrollment{}, &Subject{}, &Semester{}, &User{})

		// Retry migration
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		if err != nil {
			log.Fatal("Failed to migrate database after reset: ", err)
		}
		log.Println("✅ Database reset and migrated successfully")
	}

	// Seed database with sample data (comment out after first run if you want to keep data)
	// Seed database if SEED_DB environment variable is set to "true"
	if os.Getenv("SEED_DB") == "true" {
		log.Println("🌱 SEED_DB=true detected, seeding database...")
		seedDatabase(db)
	}

	r := gin.Default()

	// Apply CORS middleware to all routes
	r.Use(CORSMiddleware())

	// Public endpoints
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Student Feedback System API")
	})

	r.GET("/quote", func(c *gin.Context) {
		c.JSON(200, quote.Go())
	})

	// Get current active semester (public endpoint)
	r.GET("/current-semester", func(c *gin.Context) {
		var semester Semester
		if err := db.Where("is_active = ?", true).First(&semester).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No active semester found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"semester": semester})
	})

	// Authentication endpoints
	r.POST("/register", func(c *gin.Context) {
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

		// Determine requested role (what the user is asking for)
		requestedRole := newUser.RequestedRole
		if requestedRole == "" {
			// Backwards compatibility: if only Role was sent, treat it as requested role
			if newUser.Role != "" {
				requestedRole = newUser.Role
			} else {
				requestedRole = RoleStudent
			}
		}

		// Validate requested role
		if requestedRole != RoleStudent && requestedRole != RoleProfessor && requestedRole != RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requested role"})
			return
		}

		// Store requested role separately
		newUser.RequestedRole = requestedRole

		// Effective role is always student at registration time
		newUser.Role = RoleStudent

		// Hash the password before storing
		hashedPassword, err := HashPassword(newUser.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}
		newUser.Password = hashedPassword

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
			"id":             newUser.ID,
			"first_name":     newUser.FirstName,
			"last_name":      newUser.LastName,
			"email":          newUser.Email,
			"role":           newUser.Role,
			"requested_role": newUser.RequestedRole,
			"created_at":     newUser.CreatedAt,
			"updated_at":     newUser.UpdatedAt,
		}

		c.JSON(http.StatusCreated, responseUser)
	})

	r.POST("/login", func(c *gin.Context) {
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

		// Verify password using bcrypt
		if !CheckPasswordHash(user.Password, foundUser.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Generate JWT token
		token, err := GenerateJWT(foundUser.ID, foundUser.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Create response with user data and token
		response := map[string]interface{}{
			"token": token,
			"user": map[string]interface{}{
				"id":             foundUser.ID,
				"first_name":     foundUser.FirstName,
				"last_name":      foundUser.LastName,
				"email":          foundUser.Email,
				"role":           foundUser.Role,
				"requested_role": foundUser.RequestedRole,
				"created_at":     foundUser.CreatedAt,
				"updated_at":     foundUser.UpdatedAt,
			},
		}

		c.JSON(http.StatusOK, response)
	})

	// =============================================================================
	// ADMIN ENDPOINTS
	// =============================================================================

	adminGroup := r.Group("/admin")
	adminGroup.Use(RequireRole(RoleAdmin))
	{
		// Semester Management
		adminGroup.POST("/semesters", func(c *gin.Context) {
			var semester Semester
			if err := c.BindJSON(&semester); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}
			if err := db.Create(&semester).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create semester"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"semester": semester})
		})

		adminGroup.GET("/semesters", func(c *gin.Context) {
			var semesters []Semester
			if err := db.Find(&semesters).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch semesters"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"semesters": semesters})
		})

		adminGroup.PUT("/semesters/:id/activate", func(c *gin.Context) {
			id := c.Param("id")
			// Deactivate all semesters first
			db.Model(&Semester{}).Update("is_active", false)
			// Activate the selected semester
			if err := db.Model(&Semester{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate semester"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Semester activated successfully"})
		})

		// Subject Management
		adminGroup.POST("/subjects", func(c *gin.Context) {
			var subject Subject
			if err := c.BindJSON(&subject); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}
			if err := db.Create(&subject).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subject"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"subject": subject})
		})

		adminGroup.GET("/subjects", func(c *gin.Context) {
			var subjects []Subject
			if err := db.Preload("Professor").Find(&subjects).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"subjects": subjects})
		})

		// Student Enrollment Management
		adminGroup.POST("/enrollments", func(c *gin.Context) {
			var enrollment StudentEnrollment
			if err := c.BindJSON(&enrollment); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}
			if err := db.Create(&enrollment).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enrollment"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"enrollment": enrollment})
		})

		adminGroup.GET("/enrollments", func(c *gin.Context) {
			var enrollments []StudentEnrollment
			if err := db.Preload("Student").Preload("Subject").Preload("Semester").Find(&enrollments).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
		})

		// View All Responses
		adminGroup.GET("/responses", func(c *gin.Context) {
			var responses []Response
			if err := db.Preload("Survey").Preload("Student").Preload("Question").Find(&responses).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"responses": responses})
		})

		// Get all users
		adminGroup.GET("/users", func(c *gin.Context) {
			var users []User
			if err := db.Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"users": users})
		})

		// List users whose requested role differs from their current role
		adminGroup.GET("/role-requests", func(c *gin.Context) {
			var users []User
			if err := db.Where("requested_role <> role").Find(&users).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch role requests"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"users": users})
		})

		// Update a user's effective role (approve/reject role requests)
		adminGroup.PUT("/users/:id/role", func(c *gin.Context) {
			userIDParam := c.Param("id")
			userID, err := strconv.ParseUint(userIDParam, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			var body struct {
				Role string `json:"role"`
			}
			if err := c.BindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			if body.Role != RoleStudent && body.Role != RoleProfessor && body.Role != RoleAdmin {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
				return
			}

			var user User
			if err := db.First(&user, userID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			user.Role = body.Role
			user.RequestedRole = body.Role

			if err := db.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"user": user})
		})

		// Seed database endpoint (admin only)
		adminGroup.POST("/seed", func(c *gin.Context) {
			seedDatabase(db)
			c.JSON(http.StatusOK, gin.H{"message": "Database seeded successfully"})
		})
	}

	// =============================================================================
	// PROFESSOR ENDPOINTS
	// =============================================================================

	professorGroup := r.Group("/professor")
	professorGroup.Use(RequireRole(RoleProfessor))
	{
		// Get professor's subjects
		professorGroup.GET("/subjects", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var subjects []Subject
			if err := db.Where("professor_id = ?", user.ID).Find(&subjects).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subjects"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"subjects": subjects})
		})

		// Create survey
		professorGroup.POST("/surveys", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var survey Survey
			if err := c.BindJSON(&survey); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			// Verify that the professor owns the subject
			var subject Subject
			if err := db.Where("id = ? AND professor_id = ?", survey.SubjectID, user.ID).First(&subject).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "You can only create surveys for your subjects"})
				return
			}

			survey.ProfessorID = user.ID
			if err := db.Create(&survey).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create survey"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"survey": survey})
		})

		// Get professor's surveys
		professorGroup.GET("/surveys", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var surveys []Survey
			if err := db.Preload("Subject").Preload("Semester").Preload("Questions").Where("professor_id = ?", user.ID).Find(&surveys).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch surveys"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"surveys": surveys})
		})

		// Add question to survey
		professorGroup.POST("/surveys/:id/questions", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")

			// Verify survey ownership
			var survey Survey
			if err := db.Where("id = ? AND professor_id = ?", surveyID, user.ID).First(&survey).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Survey not found or access denied"})
				return
			}

			var question Question
			if err := c.BindJSON(&question); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			surveyIDUint, _ := strconv.ParseUint(surveyID, 10, 32)
			question.SurveyID = uint(surveyIDUint)
			if err := db.Create(&question).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"question": question})
		})

		// Update question
		professorGroup.PUT("/surveys/:id/questions/:questionId", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")
			questionID := c.Param("questionId")

			// Verify survey ownership
			var survey Survey
			if err := db.Where("id = ? AND professor_id = ?", surveyID, user.ID).First(&survey).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Survey not found or access denied"})
				return
			}

			// Find the question
			var question Question
			if err := db.Where("id = ? AND survey_id = ?", questionID, surveyID).First(&question).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
				return
			}

			// Bind update data
			var updateData struct {
				Text     string `json:"text"`
				Type     string `json:"type"`
				Required bool   `json:"required"`
				Options  string `json:"options"`
				Order    int    `json:"order"`
			}
			if err := c.BindJSON(&updateData); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			// Update fields
			if updateData.Text != "" {
				question.Text = updateData.Text
			}
			if updateData.Type != "" {
				question.Type = updateData.Type
			}
			question.Required = updateData.Required
			if updateData.Options != "" {
				question.Options = updateData.Options
			}
			if updateData.Order > 0 {
				question.Order = updateData.Order
			}

			if err := db.Save(&question).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"question": question})
		})

		// Delete question
		professorGroup.DELETE("/surveys/:id/questions/:questionId", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")
			questionID := c.Param("questionId")

			// Verify survey ownership
			var survey Survey
			if err := db.Where("id = ? AND professor_id = ?", surveyID, user.ID).First(&survey).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Survey not found or access denied"})
				return
			}

			// Delete the question
			if err := db.Where("id = ? AND survey_id = ?", questionID, surveyID).Delete(&Question{}).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Question deleted successfully"})
		})

		// Get responses for professor's surveys
		professorGroup.GET("/responses", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var responses []Response
			if err := db.Preload("Survey").Preload("Student").Preload("Question").
				Joins("JOIN surveys ON responses.survey_id = surveys.id").
				Where("surveys.professor_id = ?", user.ID).
				Find(&responses).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"responses": responses})
		})

		// Get responses for specific survey
		professorGroup.GET("/surveys/:id/responses", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")

			// Verify survey ownership
			var survey Survey
			if err := db.Where("id = ? AND professor_id = ?", surveyID, user.ID).First(&survey).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "Survey not found or access denied"})
				return
			}

			var responses []Response
			if err := db.Preload("Student").Preload("Question").Where("survey_id = ?", surveyID).Find(&responses).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"responses": responses})
		})
	}

	// =============================================================================
	// STUDENT ENDPOINTS
	// =============================================================================

	studentGroup := r.Group("/student")
	studentGroup.Use(RequireRole(RoleStudent))
	{
		// Get student's enrolled subjects
		studentGroup.GET("/subjects", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var enrollments []StudentEnrollment
			if err := db.Preload("Subject").Preload("Semester").Where("student_id = ?", user.ID).Find(&enrollments).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch enrollments"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
		})

		// Get available surveys for student
		studentGroup.GET("/surveys", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var surveys []Survey
			if err := db.Preload("Subject").Preload("Semester").Preload("Questions").
				Joins("JOIN student_enrollments ON surveys.subject_id = student_enrollments.subject_id AND surveys.semester_id = student_enrollments.semester_id").
				Where("student_enrollments.student_id = ? AND surveys.is_active = ?", user.ID, true).
				Find(&surveys).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch surveys"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"surveys": surveys})
		})

		// Submit response to survey
		studentGroup.POST("/responses", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var response Response
			if err := c.BindJSON(&response); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
				return
			}

			// Verify student is enrolled in the survey's subject
			var enrollment StudentEnrollment
			if err := db.Joins("JOIN surveys ON student_enrollments.subject_id = surveys.subject_id AND student_enrollments.semester_id = surveys.semester_id").
				Where("student_enrollments.student_id = ? AND surveys.id = ?", user.ID, response.SurveyID).
				First(&enrollment).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "You are not enrolled in this survey's subject"})
				return
			}

			response.StudentID = user.ID
			if err := db.Create(&response).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit response"})
				return
			}
			c.JSON(http.StatusCreated, gin.H{"response": response})
		})

		// Get student's past responses
		studentGroup.GET("/responses", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)

			var responses []Response
			if err := db.Preload("Survey").Preload("Question").Where("student_id = ?", user.ID).Find(&responses).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"responses": responses})
		})

		// Get specific survey details with questions (for taking survey)
		studentGroup.GET("/surveys/:id", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")

			var survey Survey
			if err := db.Preload("Subject").Preload("Semester").Preload("Questions", func(db *gorm.DB) *gorm.DB {
				return db.Order("\"order\" ASC")
			}).
				Joins("JOIN student_enrollments ON surveys.subject_id = student_enrollments.subject_id AND surveys.semester_id = student_enrollments.semester_id").
				Where("student_enrollments.student_id = ? AND surveys.id = ? AND surveys.is_active = ?", user.ID, surveyID, true).
				First(&survey).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found or access denied"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"survey": survey})
		})

		// Get student's responses for a specific survey
		studentGroup.GET("/surveys/:id/responses", func(c *gin.Context) {
			currentUser, _ := c.Get("currentUser")
			user := currentUser.(User)
			surveyID := c.Param("id")

			// Verify student has access to this survey
			var survey Survey
			if err := db.Joins("JOIN student_enrollments ON surveys.subject_id = student_enrollments.subject_id AND surveys.semester_id = student_enrollments.semester_id").
				Where("student_enrollments.student_id = ? AND surveys.id = ?", user.ID, surveyID).
				First(&survey).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Survey not found or access denied"})
				return
			}

			var responses []Response
			if err := db.Preload("Question").Where("survey_id = ? AND student_id = ?", surveyID, user.ID).Find(&responses).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"responses": responses})
		})
	}

	// Legacy endpoint - can be removed later
	r.POST("/consulta", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This endpoint is deprecated. Use the new survey system."})
	})

	r.Run(":3030")
}
