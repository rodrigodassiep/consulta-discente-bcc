package main

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func seedDatabase(db *gorm.DB) {
	log.Println("🌱 Starting database seeding...")

	// Check if we already have data - if so, skip seeding
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		log.Println("📋 Database already has data, skipping seeding...")
		return
	}

	// Clear existing data (optional - comment out if you want to keep existing data)
	log.Println("🧹 Clearing existing data...")
	db.Exec("DELETE FROM responses")
	db.Exec("DELETE FROM questions")
	db.Exec("DELETE FROM surveys")
	db.Exec("DELETE FROM student_enrollments")
	db.Exec("DELETE FROM subjects")
	db.Exec("DELETE FROM semesters")
	db.Exec("DELETE FROM users")

	// Create users
	log.Println("👥 Creating users...")

	// Admin user
	admin := User{
		FirstName: "Carlos",
		LastName:  "Administrator",
		Email:     "admin@univasf.edu.br",
		Password:  "admin123",
		Role:      RoleAdmin,
	}
	db.Create(&admin)

	// Professor users
	professors := []User{
		{
			FirstName: "Maria",
			LastName:  "Silva",
			Email:     "maria.silva@univasf.edu.br",
			Password:  "prof123",
			Role:      RoleProfessor,
		},
		{
			FirstName: "João",
			LastName:  "Santos",
			Email:     "joao.santos@univasf.edu.br",
			Password:  "prof123",
			Role:      RoleProfessor,
		},
		{
			FirstName: "Ana",
			LastName:  "Costa",
			Email:     "ana.costa@univasf.edu.br",
			Password:  "prof123",
			Role:      RoleProfessor,
		},
	}

	for _, prof := range professors {
		db.Create(&prof)
	}

	// Student users
	students := []User{
		{
			FirstName: "Pedro",
			LastName:  "Oliveira",
			Email:     "pedro.oliveira@discente.univasf.edu.br",
			Password:  "student123",
			Role:      RoleStudent,
		},
		{
			FirstName: "Julia",
			LastName:  "Ferreira",
			Email:     "julia.ferreira@discente.univasf.edu.br",
			Password:  "student123",
			Role:      RoleStudent,
		},
		{
			FirstName: "Lucas",
			LastName:  "Almeida",
			Email:     "lucas.almeida@discente.univasf.edu.br",
			Password:  "student123",
			Role:      RoleStudent,
		},
		{
			FirstName: "Carla",
			LastName:  "Mendes",
			Email:     "carla.mendes@discente.univasf.edu.br",
			Password:  "student123",
			Role:      RoleStudent,
		},
		{
			FirstName: "Rafael",
			LastName:  "Lima",
			Email:     "rafael.lima@discente.univasf.edu.br",
			Password:  "student123",
			Role:      RoleStudent,
		},
	}

	for _, student := range students {
		db.Create(&student)
	}

	// Create semesters
	log.Println("📅 Creating semesters...")
	semesters := []Semester{
		{
			Name:      "2023.2",
			Year:      2023,
			Period:    2,
			StartDate: time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2023, 12, 20, 23, 59, 59, 0, time.UTC),
			IsActive:  false,
		},
		{
			Name:      "2024.1",
			Year:      2024,
			Period:    1,
			StartDate: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 7, 31, 23, 59, 59, 0, time.UTC),
			IsActive:  true, // Current semester
		},
		{
			Name:      "2024.2",
			Year:      2024,
			Period:    2,
			StartDate: time.Date(2024, 8, 1, 0, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2024, 12, 20, 23, 59, 59, 0, time.UTC),
			IsActive:  false,
		},
	}

	for _, semester := range semesters {
		db.Create(&semester)
	}

	// Get created users for foreign key references
	var createdProfessors []User
	var createdStudents []User
	var currentSemester Semester

	db.Where("role = ?", RoleProfessor).Find(&createdProfessors)
	db.Where("role = ?", RoleStudent).Find(&createdStudents)
	db.Where("is_active = ?", true).First(&currentSemester)

	// Check if we have the required data
	if len(createdProfessors) == 0 {
		log.Println("❌ No professors found, cannot create subjects")
		return
	}
	if len(createdStudents) == 0 {
		log.Println("❌ No students found, cannot create enrollments")
		return
	}
	if currentSemester.ID == 0 {
		log.Println("❌ No active semester found, cannot create subjects/enrollments")
		return
	}

	// Create subjects
	log.Println("📚 Creating subjects...")
	subjects := []Subject{
		{
			Name:        "Estruturas de Dados",
			Code:        "COMP001",
			Description: "Introdução às estruturas de dados fundamentais",
			ProfessorID: createdProfessors[0].ID,
		},
		{
			Name:        "Programação Orientada a Objetos",
			Code:        "COMP002",
			Description: "Conceitos e práticas de POO",
			ProfessorID: createdProfessors[0].ID,
		},
		{
			Name:        "Banco de Dados",
			Code:        "COMP003",
			Description: "Sistemas de gerenciamento de banco de dados",
			ProfessorID: createdProfessors[1].ID,
		},
		{
			Name:        "Engenharia de Software",
			Code:        "COMP004",
			Description: "Metodologias e processos de desenvolvimento",
			ProfessorID: createdProfessors[1].ID,
		},
		{
			Name:        "Algoritmos Avançados",
			Code:        "COMP005",
			Description: "Algoritmos de otimização e complexidade",
			ProfessorID: createdProfessors[2].ID,
		},
		{
			Name:        "Inteligência Artificial",
			Code:        "COMP006",
			Description: "Fundamentos de IA e machine learning",
			ProfessorID: createdProfessors[2].ID,
		},
	}

	for _, subject := range subjects {
		db.Create(&subject)
	}

	// Get created subjects
	var createdSubjects []Subject
	db.Find(&createdSubjects)

	if len(createdSubjects) == 0 {
		log.Println("❌ No subjects found, cannot create enrollments")
		return
	}

	// Create student enrollments
	log.Println("📝 Creating student enrollments...")
	enrollments := []StudentEnrollment{
		// Pedro enrolled in 4 subjects
		{StudentID: createdStudents[0].ID, SubjectID: createdSubjects[0].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[0].ID, SubjectID: createdSubjects[1].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[0].ID, SubjectID: createdSubjects[2].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[0].ID, SubjectID: createdSubjects[3].ID, SemesterID: currentSemester.ID},

		// Julia enrolled in 3 subjects
		{StudentID: createdStudents[1].ID, SubjectID: createdSubjects[1].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[1].ID, SubjectID: createdSubjects[3].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[1].ID, SubjectID: createdSubjects[4].ID, SemesterID: currentSemester.ID},

		// Lucas enrolled in 5 subjects
		{StudentID: createdStudents[2].ID, SubjectID: createdSubjects[0].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[2].ID, SubjectID: createdSubjects[2].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[2].ID, SubjectID: createdSubjects[3].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[2].ID, SubjectID: createdSubjects[4].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[2].ID, SubjectID: createdSubjects[5].ID, SemesterID: currentSemester.ID},

		// Carla enrolled in 3 subjects
		{StudentID: createdStudents[3].ID, SubjectID: createdSubjects[1].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[3].ID, SubjectID: createdSubjects[4].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[3].ID, SubjectID: createdSubjects[5].ID, SemesterID: currentSemester.ID},

		// Rafael enrolled in 2 subjects
		{StudentID: createdStudents[4].ID, SubjectID: createdSubjects[0].ID, SemesterID: currentSemester.ID},
		{StudentID: createdStudents[4].ID, SubjectID: createdSubjects[2].ID, SemesterID: currentSemester.ID},
	}

	for _, enrollment := range enrollments {
		db.Create(&enrollment)
	}

	// Create surveys
	log.Println("📊 Creating surveys...")
	now := time.Now()

	surveys := []Survey{
		{
			Title:       "Avaliação da Disciplina - Estruturas de Dados",
			Description: "Avalie a qualidade da disciplina e metodologia de ensino",
			SubjectID:   createdSubjects[0].ID,
			SemesterID:  currentSemester.ID,
			ProfessorID: createdProfessors[0].ID,
			IsActive:    true,
			OpenDate:    now.AddDate(0, 0, -7), // Opened 7 days ago
			CloseDate:   now.AddDate(0, 0, 14), // Closes in 14 days
		},
		{
			Title:       "Feedback Semestral - POO",
			Description: "Compartilhe sua experiência com a disciplina de Programação Orientada a Objetos",
			SubjectID:   createdSubjects[1].ID,
			SemesterID:  currentSemester.ID,
			ProfessorID: createdProfessors[0].ID,
			IsActive:    true,
			OpenDate:    now.AddDate(0, 0, -3), // Opened 3 days ago
			CloseDate:   now.AddDate(0, 0, 21), // Closes in 21 days
		},
		{
			Title:       "Avaliação - Banco de Dados",
			Description: "Avalie o conteúdo e dinâmica das aulas",
			SubjectID:   createdSubjects[2].ID,
			SemesterID:  currentSemester.ID,
			ProfessorID: createdProfessors[1].ID,
			IsActive:    true,
			OpenDate:    now.AddDate(0, 0, -1), // Opened yesterday
			CloseDate:   now.AddDate(0, 0, 30), // Closes in 30 days
		},
		{
			Title:       "Feedback Final - Engenharia de Software",
			Description: "Avaliação final da disciplina",
			SubjectID:   createdSubjects[3].ID,
			SemesterID:  currentSemester.ID,
			ProfessorID: createdProfessors[1].ID,
			IsActive:    false, // Inactive survey
			OpenDate:    now.AddDate(0, 0, -30),
			CloseDate:   now.AddDate(0, 0, -7),
		},
		{
			Title:       "Pesquisa de Satisfação - IA",
			Description: "Como você avalia a disciplina de Inteligência Artificial?",
			SubjectID:   createdSubjects[5].ID,
			SemesterID:  currentSemester.ID,
			ProfessorID: createdProfessors[2].ID,
			IsActive:    true,
			OpenDate:    now.AddDate(0, 0, 2), // Opens in 2 days (upcoming)
			CloseDate:   now.AddDate(0, 0, 45),
		},
	}

	for _, survey := range surveys {
		db.Create(&survey)
	}

	// Get created surveys
	var createdSurveys []Survey
	db.Find(&createdSurveys)

	// Create questions for surveys
	log.Println("❓ Creating questions...")

	// Questions for Survey 1 (Estruturas de Dados)
	survey1Questions := []Question{
		{
			SurveyID: createdSurveys[0].ID,
			Type:     QuestionTypeNPS,
			Text:     "Em uma escala de 0 a 10, o quanto você recomendaria esta disciplina para outros estudantes?",
			Required: true,
			Order:    1,
		},
		{
			SurveyID: createdSurveys[0].ID,
			Type:     QuestionTypeRating,
			Text:     "Como você avalia a didática do professor?",
			Required: true,
			Order:    2,
		},
		{
			SurveyID: createdSurveys[0].ID,
			Type:     QuestionTypeChoice,
			Text:     "Qual aspecto da disciplina você mais gostou?",
			Required: false,
			Order:    3,
			Options:  `["Conteúdo teórico", "Exercícios práticos", "Metodologia de ensino", "Material didático", "Avaliações"]`,
		},
		{
			SurveyID: createdSurveys[0].ID,
			Type:     QuestionTypeFreeText,
			Text:     "Deixe sugestões para melhorar a disciplina:",
			Required: false,
			Order:    4,
		},
	}

	// Questions for Survey 2 (POO)
	survey2Questions := []Question{
		{
			SurveyID: createdSurveys[1].ID,
			Type:     QuestionTypeRating,
			Text:     "Como você avalia a dificuldade da disciplina?",
			Required: true,
			Order:    1,
		},
		{
			SurveyID: createdSurveys[1].ID,
			Type:     QuestionTypeChoice,
			Text:     "Qual linguagem de programação você prefere para POO?",
			Required: false,
			Order:    2,
			Options:  `["Java", "Python", "C++", "C#", "JavaScript"]`,
		},
		{
			SurveyID: createdSurveys[1].ID,
			Type:     QuestionTypeNPS,
			Text:     "Você recomendaria esta disciplina? (0-10)",
			Required: true,
			Order:    3,
		},
		{
			SurveyID: createdSurveys[1].ID,
			Type:     QuestionTypeFreeText,
			Text:     "O que você achou mais desafiador na disciplina?",
			Required: false,
			Order:    4,
		},
	}

	// Questions for Survey 3 (Banco de Dados)
	survey3Questions := []Question{
		{
			SurveyID: createdSurveys[2].ID,
			Type:     QuestionTypeChoice,
			Text:     "Qual tópico você achou mais interessante?",
			Required: true,
			Order:    1,
			Options:  `["Modelagem ER", "SQL", "Normalização", "Transações", "NoSQL"]`,
		},
		{
			SurveyID: createdSurveys[2].ID,
			Type:     QuestionTypeRating,
			Text:     "Como você avalia os exercícios práticos?",
			Required: true,
			Order:    2,
		},
		{
			SurveyID: createdSurveys[2].ID,
			Type:     QuestionTypeFreeText,
			Text:     "Comentários gerais sobre a disciplina:",
			Required: false,
			Order:    3,
		},
	}

	// Questions for Survey 5 (IA)
	survey5Questions := []Question{
		{
			SurveyID: createdSurveys[4].ID,
			Type:     QuestionTypeNPS,
			Text:     "O quanto você está satisfeito com a disciplina? (0-10)",
			Required: true,
			Order:    1,
		},
		{
			SurveyID: createdSurveys[4].ID,
			Type:     QuestionTypeChoice,
			Text:     "Qual área de IA você tem mais interesse?",
			Required: false,
			Order:    2,
			Options:  `["Machine Learning", "Deep Learning", "Processamento de Linguagem Natural", "Visão Computacional", "Robótica"]`,
		},
	}

	// Create all questions
	allQuestions := append(survey1Questions, survey2Questions...)
	allQuestions = append(allQuestions, survey3Questions...)
	allQuestions = append(allQuestions, survey5Questions...)

	for _, question := range allQuestions {
		db.Create(&question)
	}

	// Create some sample responses
	log.Println("💬 Creating sample responses...")

	// Get some questions for responses
	var questions []Question
	db.Limit(6).Find(&questions)

	sampleResponses := []Response{
		{
			SurveyID:   createdSurveys[0].ID,
			StudentID:  createdStudents[1].ID, // Julia responding to Estruturas de Dados (she's not enrolled, but we'll allow for demo)
			QuestionID: questions[0].ID,
			Answer:     "8",
		},
		{
			SurveyID:   createdSurveys[0].ID,
			StudentID:  createdStudents[1].ID,
			QuestionID: questions[1].ID,
			Answer:     "4",
		},
		{
			SurveyID:   createdSurveys[1].ID,
			StudentID:  createdStudents[0].ID, // Pedro responding to POO
			QuestionID: questions[4].ID,
			Answer:     "3",
		},
	}

	for _, response := range sampleResponses {
		db.Create(&response)
	}

	log.Println("✅ Database seeding completed successfully!")
	log.Println("📊 Created:")
	log.Println("   - 1 Admin user")
	log.Println("   - 3 Professor users")
	log.Println("   - 5 Student users")
	log.Println("   - 3 Semesters (2024.1 is active)")
	log.Println("   - 6 Subjects")
	log.Println("   - 16 Student enrollments")
	log.Println("   - 5 Surveys (4 active, 1 inactive)")
	log.Println("   - 13 Questions")
	log.Println("   - 3 Sample responses")
	log.Println("")
	log.Println("🔑 Test credentials:")
	log.Println("   Admin: admin@univasf.edu.br / admin123")
	log.Println("   Professor: maria.silva@univasf.edu.br / prof123")
	log.Println("   Student: pedro.oliveira@discente.univasf.edu.br / student123")
}
