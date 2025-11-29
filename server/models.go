package main

import (
	"time"
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
