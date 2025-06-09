package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

// User model with proper role handling
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Role      string    `json:"role" gorm:"not null;check:role IN ('student','professor','admin')"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// CORSMiddleware handles OPTIONS requests and sets CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle OPTIONS requests
		if c.Request.Method == "OPTIONS" {
			log.Println("OPTIONS request received")
			c.AbortWithStatus(204)
			return
		}

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

	db, err_sql := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err_sql != nil {
		panic("failed to connect database")
	}

	// Auto-migrate all the new models
	db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})

	r := gin.Default()

	// Apply CORS middleware to all routes
	r.Use(CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hello world")
	})

	r.GET("/quote", func(c *gin.Context) {
		c.JSON(200, quote.Go())
	})

	r.POST("/register", func(c *gin.Context) {
		var newUser User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(500, "ERROR")
		}
		result := db.Create(&newUser)
		if result.Error != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(200, gin.H{"user": newUser})
	})

	r.POST("/login", func(c *gin.Context) {
		var user User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(500, "ERROR")
		}

		var foundUser User
		result := db.Where("email = ?", user.Email).First(&foundUser)
		if result.Error != nil {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}

		if foundUser.Password != user.Password {
			c.JSON(401, gin.H{"error": "Invalid credentials"})
			return
		}
		c.JSON(200, gin.H{"user": foundUser})

		// quando logar, criar um hashmap autorizando o usuario
	})

	// Legacy endpoint - can be removed later
	r.POST("/consulta", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This endpoint is deprecated. Use the new survey system."})
	})

	r.Run(":3030")
}
