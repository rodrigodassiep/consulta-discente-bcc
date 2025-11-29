package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSeedDatabase(t *testing.T) {
	t.Run("Fresh Database Seeding", func(t *testing.T) {
		// Create a fresh in-memory database
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate all models
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		assert.NoError(t, err)

		// Verify database is empty
		var userCount int64
		db.Model(&User{}).Count(&userCount)
		assert.Equal(t, int64(0), userCount)

		// Run seeding
		seedDatabase(db)

		// Verify users were created
		db.Model(&User{}).Count(&userCount)
		assert.Greater(t, userCount, int64(0))

		// Verify admin user exists
		var admin User
		err = db.Where("role = ? AND email = ?", RoleAdmin, "admin@usp.br").First(&admin).Error
		assert.NoError(t, err)
		assert.Equal(t, "Carlos", admin.FirstName)
		assert.Equal(t, "Administrator", admin.LastName)
		assert.Equal(t, RoleAdmin, admin.Role)

		// Verify professors were created
		var professorCount int64
		db.Model(&User{}).Where("role = ?", RoleProfessor).Count(&professorCount)
		assert.Equal(t, int64(3), professorCount) // Should have 3 professors

		// Verify students were created
		var studentCount int64
		db.Model(&User{}).Where("role = ?", RoleStudent).Count(&studentCount)
		assert.Equal(t, int64(5), studentCount) // Should have 5 students

		// Verify semesters were created
		var semesterCount int64
		db.Model(&Semester{}).Count(&semesterCount)
		assert.Greater(t, semesterCount, int64(0))

		// Verify there's an active semester
		var activeSemester Semester
		err = db.Where("is_active = ?", true).First(&activeSemester).Error
		assert.NoError(t, err)
		assert.Equal(t, "2024.1", activeSemester.Name)
		assert.Equal(t, 2024, activeSemester.Year)
		assert.Equal(t, 1, activeSemester.Period)

		// Verify subjects were created
		var subjectCount int64
		db.Model(&Subject{}).Count(&subjectCount)
		assert.Greater(t, subjectCount, int64(0))

		// Verify enrollments were created
		var enrollmentCount int64
		db.Model(&StudentEnrollment{}).Count(&enrollmentCount)
		assert.Greater(t, enrollmentCount, int64(0))

		// Verify surveys were created
		var surveyCount int64
		db.Model(&Survey{}).Count(&surveyCount)
		assert.Greater(t, surveyCount, int64(0))

		// Verify questions were created
		var questionCount int64
		db.Model(&Question{}).Count(&questionCount)
		assert.Greater(t, questionCount, int64(0))

		// Verify responses were created
		var responseCount int64
		db.Model(&Response{}).Count(&responseCount)
		assert.Greater(t, responseCount, int64(0))
	})

	t.Run("Skip Seeding When Data Already Exists", func(t *testing.T) {
		// Create a database with existing data
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate all models
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		assert.NoError(t, err)

		// Create a test user to simulate existing data
		existingUser := User{
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@test.com",
			Password:  "password123",
			Role:      RoleStudent,
		}
		db.Create(&existingUser)

		// Verify user exists
		var userCount int64
		db.Model(&User{}).Count(&userCount)
		assert.Equal(t, int64(1), userCount)

		// Run seeding - should skip
		seedDatabase(db)

		// Verify no additional users were created
		db.Model(&User{}).Count(&userCount)
		assert.Equal(t, int64(1), userCount) // Still only 1 user

		// Verify the existing user is still there
		var foundUser User
		err = db.Where("email = ?", "existing@test.com").First(&foundUser).Error
		assert.NoError(t, err)
		assert.Equal(t, "Existing", foundUser.FirstName)
	})

	t.Run("Verify Specific Seeded Data", func(t *testing.T) {
		// Create a fresh database
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate all models
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		assert.NoError(t, err)

		// Run seeding
		seedDatabase(db)

		// Test specific professors
		var maria User
		err = db.Where("email = ?", "maria.silva@usp.br").First(&maria).Error
		assert.NoError(t, err)
		assert.Equal(t, "Maria", maria.FirstName)
		assert.Equal(t, "Silva", maria.LastName)
		assert.Equal(t, RoleProfessor, maria.Role)

		var joao User
		err = db.Where("email = ?", "joao.santos@usp.br").First(&joao).Error
		assert.NoError(t, err)
		assert.Equal(t, "João", joao.FirstName)
		assert.Equal(t, "Santos", joao.LastName)
		assert.Equal(t, RoleProfessor, joao.Role)

		// Test specific students
		var pedro User
		err = db.Where("email = ?", "pedro.oliveira@usp.br").First(&pedro).Error
		assert.NoError(t, err)
		assert.Equal(t, "Pedro", pedro.FirstName)
		assert.Equal(t, "Oliveira", pedro.LastName)
		assert.Equal(t, RoleStudent, pedro.Role)

		// Test semesters
		var semester2023 Semester
		err = db.Where("name = ?", "2023.2").First(&semester2023).Error
		assert.NoError(t, err)
		assert.Equal(t, 2023, semester2023.Year)
		assert.Equal(t, 2, semester2023.Period)
		assert.False(t, semester2023.IsActive)

		var semester2024 Semester
		err = db.Where("name = ?", "2024.1").First(&semester2024).Error
		assert.NoError(t, err)
		assert.Equal(t, 2024, semester2024.Year)
		assert.Equal(t, 1, semester2024.Period)
		assert.True(t, semester2024.IsActive)

		// Test subjects
		var estruturas Subject
		err = db.Where("code = ?", "COMP001").First(&estruturas).Error
		assert.NoError(t, err)
		assert.Equal(t, "Estruturas de Dados", estruturas.Name)
		assert.Equal(t, "COMP001", estruturas.Code)
		assert.NotZero(t, estruturas.ProfessorID)

		// Verify relationships
		var subjectWithProfessor Subject
		err = db.Preload("Professor").Where("code = ?", "COMP001").First(&subjectWithProfessor).Error
		assert.NoError(t, err)
		assert.Equal(t, RoleProfessor, subjectWithProfessor.Professor.Role)
	})

	t.Run("Verify Data Consistency", func(t *testing.T) {
		// Create a fresh database
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate all models
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		assert.NoError(t, err)

		// Run seeding
		seedDatabase(db)

		// Verify that all subjects have valid professors
		var subjects []Subject
		err = db.Preload("Professor").Find(&subjects).Error
		assert.NoError(t, err)

		for _, subject := range subjects {
			assert.NotZero(t, subject.ProfessorID)
			assert.Equal(t, RoleProfessor, subject.Professor.Role)
			assert.NotEmpty(t, subject.Professor.Email)
		}

		// Verify that all enrollments have valid references
		var enrollments []StudentEnrollment
		err = db.Preload("Student").Preload("Subject").Preload("Semester").Find(&enrollments).Error
		assert.NoError(t, err)

		for _, enrollment := range enrollments {
			assert.NotZero(t, enrollment.StudentID)
			assert.NotZero(t, enrollment.SubjectID)
			assert.NotZero(t, enrollment.SemesterID)
			assert.Equal(t, RoleStudent, enrollment.Student.Role)
			assert.NotEmpty(t, enrollment.Subject.Name)
			assert.NotEmpty(t, enrollment.Semester.Name)
		}

		// Verify that all surveys have valid references
		var surveys []Survey
		err = db.Preload("Professor").Preload("Subject").Preload("Semester").Find(&surveys).Error
		assert.NoError(t, err)

		for _, survey := range surveys {
			assert.NotZero(t, survey.ProfessorID)
			assert.NotZero(t, survey.SubjectID)
			assert.NotZero(t, survey.SemesterID)
			assert.Equal(t, RoleProfessor, survey.Professor.Role)
			assert.NotEmpty(t, survey.Subject.Name)
			assert.NotEmpty(t, survey.Semester.Name)
			assert.NotEmpty(t, survey.Title)
		}

		// Verify that all questions have valid survey references
		var questions []Question
		err = db.Preload("Survey").Find(&questions).Error
		assert.NoError(t, err)

		for _, question := range questions {
			assert.NotZero(t, question.SurveyID)
			assert.NotEmpty(t, question.Survey.Title)
			assert.NotEmpty(t, question.Text)
			assert.Contains(t, []string{QuestionTypeNPS, QuestionTypeFreeText, QuestionTypeRating, QuestionTypeChoice}, question.Type)
		}

		// Verify that all responses have valid references
		var responses []Response
		err = db.Preload("Student").Preload("Survey").Preload("Question").Find(&responses).Error
		assert.NoError(t, err)

		for _, response := range responses {
			assert.NotZero(t, response.StudentID)
			assert.NotZero(t, response.SurveyID)
			assert.NotZero(t, response.QuestionID)
			assert.Equal(t, RoleStudent, response.Student.Role)
			assert.NotEmpty(t, response.Survey.Title)
			assert.NotEmpty(t, response.Question.Text)
			assert.NotEmpty(t, response.Answer)
		}
	})

	t.Run("Verify Date Constraints", func(t *testing.T) {
		// Create a fresh database
		db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		assert.NoError(t, err)

		// Auto-migrate all models
		err = db.AutoMigrate(&User{}, &Subject{}, &Semester{}, &StudentEnrollment{}, &Survey{}, &Question{}, &Response{})
		assert.NoError(t, err)

		// Run seeding
		seedDatabase(db)

		// Verify semester dates make sense
		var semesters []Semester
		err = db.Find(&semesters).Error
		assert.NoError(t, err)

		for _, semester := range semesters {
			assert.True(t, semester.EndDate.After(semester.StartDate), "End date should be after start date for semester %s", semester.Name)

			// Verify year matches the semester name
			if semester.Name == "2023.2" {
				assert.Equal(t, 2023, semester.Year)
				assert.Equal(t, 2, semester.Period)
			} else if semester.Name == "2024.1" {
				assert.Equal(t, 2024, semester.Year)
				assert.Equal(t, 1, semester.Period)
			} else if semester.Name == "2024.2" {
				assert.Equal(t, 2024, semester.Year)
				assert.Equal(t, 2, semester.Period)
			}
		}

		// Verify survey dates make sense
		var surveys []Survey
		err = db.Find(&surveys).Error
		assert.NoError(t, err)

		for _, survey := range surveys {
			if !survey.CloseDate.IsZero() && !survey.OpenDate.IsZero() {
				assert.True(t, survey.CloseDate.After(survey.OpenDate), "Close date should be after open date for survey %s", survey.Title)
			}
		}
	})
}
