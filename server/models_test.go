package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate all models
	db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})

	return db
}

func TestUserModel(t *testing.T) {
	db := setupTestDB()

	t.Run("Valid User Creation", func(t *testing.T) {
		user := User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
			Role:      RoleStudent,
		}

		result := db.Create(&user)
		assert.NoError(t, result.Error)
		assert.NotZero(t, user.ID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
		assert.Equal(t, RoleStudent, user.Role)
	})

	t.Run("User Email Uniqueness", func(t *testing.T) {
		user1 := User{
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane@example.com",
			Password:  "password123",
			Role:      RoleStudent,
		}

		user2 := User{
			FirstName: "John",
			LastName:  "Smith",
			Email:     "jane@example.com", // Same email
			Password:  "password456",
			Role:      RoleProfessor,
		}

		// First user should be created successfully
		result1 := db.Create(&user1)
		assert.NoError(t, result1.Error)

		// Second user with same email should fail
		result2 := db.Create(&user2)
		assert.Error(t, result2.Error)
	})

	t.Run("User Role Validation", func(t *testing.T) {
		validRoles := []string{RoleStudent, RoleProfessor, RoleAdmin}

		for _, role := range validRoles {
			user := User{
				FirstName: "Test",
				LastName:  "User",
				Email:     "test" + role + "@example.com",
				Password:  "password123",
				Role:      role,
			}

			result := db.Create(&user)
			assert.NoError(t, result.Error, "Role %s should be valid", role)
		}
	})
}

func TestSubjectModel(t *testing.T) {
	db := setupTestDB()

	// Create a professor first
	professor := User{
		FirstName: "Prof",
		LastName:  "Smith",
		Email:     "prof@example.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	db.Create(&professor)

	t.Run("Valid Subject Creation", func(t *testing.T) {
		subject := Subject{
			Name:        "Computer Science 101",
			Code:        "CS101",
			Description: "Introduction to Computer Science",
			ProfessorID: professor.ID,
		}

		result := db.Create(&subject)
		assert.NoError(t, result.Error)
		assert.NotZero(t, subject.ID)
		assert.Equal(t, "Computer Science 101", subject.Name)
		assert.Equal(t, "CS101", subject.Code)
		assert.Equal(t, professor.ID, subject.ProfessorID)
	})

	t.Run("Subject Code Uniqueness", func(t *testing.T) {
		subject1 := Subject{
			Name:        "Math 101",
			Code:        "MATH101",
			Description: "Basic Math",
			ProfessorID: professor.ID,
		}

		subject2 := Subject{
			Name:        "Advanced Math",
			Code:        "MATH101", // Same code
			Description: "Advanced Math Topics",
			ProfessorID: professor.ID,
		}

		result1 := db.Create(&subject1)
		assert.NoError(t, result1.Error)

		result2 := db.Create(&subject2)
		assert.Error(t, result2.Error)
	})
}

func TestSemesterModel(t *testing.T) {
	db := setupTestDB()

	t.Run("Valid Semester Creation", func(t *testing.T) {
		semester := Semester{
			Name:      "2024.1",
			Year:      2024,
			Period:    1,
			StartDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC),
			IsActive:  true,
		}

		result := db.Create(&semester)
		assert.NoError(t, result.Error)
		assert.NotZero(t, semester.ID)
		assert.Equal(t, "2024.1", semester.Name)
		assert.Equal(t, 2024, semester.Year)
		assert.Equal(t, 1, semester.Period)
		assert.True(t, semester.IsActive)
	})

	t.Run("Semester Date Validation", func(t *testing.T) {
		startDate := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC)

		semester := Semester{
			Name:      "2024.2",
			Year:      2024,
			Period:    2,
			StartDate: startDate,
			EndDate:   endDate,
			IsActive:  false,
		}

		result := db.Create(&semester)
		assert.NoError(t, result.Error)
		assert.True(t, semester.EndDate.After(semester.StartDate))
	})
}

func TestStudentEnrollmentModel(t *testing.T) {
	db := setupTestDB()

	// Setup test data
	student := User{
		FirstName: "Student",
		LastName:  "Test",
		Email:     "student@example.com",
		Password:  "password123",
		Role:      RoleStudent,
	}
	db.Create(&student)

	professor := User{
		FirstName: "Prof",
		LastName:  "Test",
		Email:     "prof@example.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	db.Create(&professor)

	subject := Subject{
		Name:        "Test Subject",
		Code:        "TEST101",
		Description: "Test Description",
		ProfessorID: professor.ID,
	}
	db.Create(&subject)

	semester := Semester{
		Name:      "2024.1",
		Year:      2024,
		Period:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}
	db.Create(&semester)

	t.Run("Valid Enrollment Creation", func(t *testing.T) {
		enrollment := StudentEnrollment{
			StudentID:  student.ID,
			SubjectID:  subject.ID,
			SemesterID: semester.ID,
		}

		result := db.Create(&enrollment)
		assert.NoError(t, result.Error)
		assert.NotZero(t, enrollment.ID)
		assert.Equal(t, student.ID, enrollment.StudentID)
		assert.Equal(t, subject.ID, enrollment.SubjectID)
		assert.Equal(t, semester.ID, enrollment.SemesterID)
	})
}

func TestSurveyModel(t *testing.T) {
	db := setupTestDB()

	// Setup test data
	professor := User{
		FirstName: "Prof",
		LastName:  "Test",
		Email:     "prof@example.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	db.Create(&professor)

	subject := Subject{
		Name:        "Test Subject",
		Code:        "TEST101",
		Description: "Test Description",
		ProfessorID: professor.ID,
	}
	db.Create(&subject)

	semester := Semester{
		Name:      "2024.1",
		Year:      2024,
		Period:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}
	db.Create(&semester)

	t.Run("Valid Survey Creation", func(t *testing.T) {
		survey := Survey{
			Title:       "Course Feedback",
			Description: "Please provide feedback about the course",
			SubjectID:   subject.ID,
			SemesterID:  semester.ID,
			ProfessorID: professor.ID,
			IsActive:    true,
			OpenDate:    time.Now(),
			CloseDate:   time.Now().AddDate(0, 1, 0),
		}

		result := db.Create(&survey)
		assert.NoError(t, result.Error)
		assert.NotZero(t, survey.ID)
		assert.Equal(t, "Course Feedback", survey.Title)
		assert.True(t, survey.IsActive)
		assert.Equal(t, professor.ID, survey.ProfessorID)
	})
}

func TestQuestionModel(t *testing.T) {
	db := setupTestDB()

	// Setup test data
	professor := User{
		FirstName: "Prof",
		LastName:  "Test",
		Email:     "prof@example.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	db.Create(&professor)

	subject := Subject{
		Name:        "Test Subject",
		Code:        "TEST101",
		Description: "Test Description",
		ProfessorID: professor.ID,
	}
	db.Create(&subject)

	semester := Semester{
		Name:      "2024.1",
		Year:      2024,
		Period:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}
	db.Create(&semester)

	survey := Survey{
		Title:       "Test Survey",
		Description: "Test Description",
		SubjectID:   subject.ID,
		SemesterID:  semester.ID,
		ProfessorID: professor.ID,
		IsActive:    true,
		OpenDate:    time.Now(),
		CloseDate:   time.Now().AddDate(0, 1, 0),
	}
	db.Create(&survey)

	t.Run("Valid Question Creation", func(t *testing.T) {
		validTypes := []string{QuestionTypeNPS, QuestionTypeFreeText, QuestionTypeRating, QuestionTypeChoice}

		for i, questionType := range validTypes {
			question := Question{
				SurveyID: survey.ID,
				Type:     questionType,
				Text:     "Test question " + questionType,
				Required: true,
				Order:    i + 1,
				Options:  `["Option1", "Option2", "Option3"]`,
			}

			result := db.Create(&question)
			assert.NoError(t, result.Error, "Question type %s should be valid", questionType)
			assert.NotZero(t, question.ID)
			assert.Equal(t, questionType, question.Type)
			assert.Equal(t, survey.ID, question.SurveyID)
		}
	})
}

func TestResponseModel(t *testing.T) {
	db := setupTestDB()

	// Setup complete test data
	student := User{
		FirstName: "Student",
		LastName:  "Test",
		Email:     "student@example.com",
		Password:  "password123",
		Role:      RoleStudent,
	}
	db.Create(&student)

	professor := User{
		FirstName: "Prof",
		LastName:  "Test",
		Email:     "prof@example.com",
		Password:  "password123",
		Role:      RoleProfessor,
	}
	db.Create(&professor)

	subject := Subject{
		Name:        "Test Subject",
		Code:        "TEST101",
		Description: "Test Description",
		ProfessorID: professor.ID,
	}
	db.Create(&subject)

	semester := Semester{
		Name:      "2024.1",
		Year:      2024,
		Period:    1,
		StartDate: time.Now(),
		EndDate:   time.Now().AddDate(0, 4, 0),
		IsActive:  true,
	}
	db.Create(&semester)

	survey := Survey{
		Title:       "Test Survey",
		Description: "Test Description",
		SubjectID:   subject.ID,
		SemesterID:  semester.ID,
		ProfessorID: professor.ID,
		IsActive:    true,
		OpenDate:    time.Now(),
		CloseDate:   time.Now().AddDate(0, 1, 0),
	}
	db.Create(&survey)

	question := Question{
		SurveyID: survey.ID,
		Type:     QuestionTypeFreeText,
		Text:     "How do you rate this course?",
		Required: true,
		Order:    1,
	}
	db.Create(&question)

	t.Run("Valid Response Creation", func(t *testing.T) {
		response := Response{
			SurveyID:   survey.ID,
			StudentID:  student.ID,
			QuestionID: question.ID,
			Answer:     "This course is excellent!",
		}

		result := db.Create(&response)
		assert.NoError(t, result.Error)
		assert.NotZero(t, response.ID)
		assert.Equal(t, "This course is excellent!", response.Answer)
		assert.Equal(t, student.ID, response.StudentID)
		assert.Equal(t, survey.ID, response.SurveyID)
		assert.Equal(t, question.ID, response.QuestionID)
	})
}

// TestConstants verifies that all constants are defined correctly
func TestConstants(t *testing.T) {
	t.Run("User Role Constants", func(t *testing.T) {
		assert.Equal(t, "student", RoleStudent)
		assert.Equal(t, "professor", RoleProfessor)
		assert.Equal(t, "admin", RoleAdmin)
	})

	t.Run("Question Type Constants", func(t *testing.T) {
		assert.Equal(t, "nps", QuestionTypeNPS)
		assert.Equal(t, "free_text", QuestionTypeFreeText)
		assert.Equal(t, "rating", QuestionTypeRating)
		assert.Equal(t, "multiple_choice", QuestionTypeChoice)
	})
}
