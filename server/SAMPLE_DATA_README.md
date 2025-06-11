# Sample Data for Student Feedback System

## Overview

The system will automatically populate the database with comprehensive sample data when you first run the server. This includes realistic users, academic data, and surveys to test all functionality.

## 🔑 Test Credentials

### Admin User
- **Email**: `admin@univasf.edu.br`
- **Password**: `admin123`
- **Access**: Full system administration

### Professor Users
- **Maria Silva**: `maria.silva@univasf.edu.br` / `prof123`
- **João Santos**: `joao.santos@univasf.edu.br` / `prof123`
- **Ana Costa**: `ana.costa@univasf.edu.br` / `prof123`

### Student Users
- **Pedro Oliveira**: `pedro.oliveira@discente.univasf.edu.br` / `student123`
- **Julia Ferreira**: `julia.ferreira@discente.univasf.edu.br` / `student123`
- **Lucas Almeida**: `lucas.almeida@discente.univasf.edu.br` / `student123`
- **Carla Mendes**: `carla.mendes@discente.univasf.edu.br` / `student123`
- **Rafael Lima**: `rafael.lima@discente.univasf.edu.br` / `student123`

## 📊 Sample Data Structure

### Semesters
- **2023.2**: Past semester (inactive)
- **2024.1**: Current semester (active) ⭐
- **2024.2**: Future semester (inactive)

### Subjects (Computer Science)
1. **Estruturas de Dados** (COMP001) - Prof. Maria Silva
2. **Programação Orientada a Objetos** (COMP002) - Prof. Maria Silva
3. **Banco de Dados** (COMP003) - Prof. João Santos
4. **Engenharia de Software** (COMP004) - Prof. João Santos
5. **Algoritmos Avançados** (COMP005) - Prof. Ana Costa
6. **Inteligência Artificial** (COMP006) - Prof. Ana Costa

### Student Enrollments (2024.1)
- **Pedro**: 4 subjects (COMP001, COMP002, COMP003, COMP004)
- **Julia**: 3 subjects (COMP002, COMP004, COMP005)
- **Lucas**: 5 subjects (COMP001, COMP003, COMP004, COMP005, COMP006)
- **Carla**: 3 subjects (COMP002, COMP005, COMP006)
- **Rafael**: 2 subjects (COMP001, COMP003)

### Active Surveys
1. **Estruturas de Dados** - "Avaliação da Disciplina" ✅ Active
   - 4 questions (NPS, Rating, Multiple Choice, Free Text)
   
2. **POO** - "Feedback Semestral" ✅ Active
   - 4 questions (Rating, Multiple Choice, NPS, Free Text)
   
3. **Banco de Dados** - "Avaliação" ✅ Active
   - 3 questions (Multiple Choice, Rating, Free Text)
   
4. **Engenharia de Software** - "Feedback Final" ❌ Inactive
   
5. **Inteligência Artificial** - "Pesquisa de Satisfação" ⏰ Upcoming (opens in 2 days)
   - 2 questions (NPS, Multiple Choice)

### Question Types Examples
- **NPS**: "Em uma escala de 0 a 10, o quanto você recomendaria esta disciplina?"
- **Rating**: "Como você avalia a didática do professor?" (1-5 stars)
- **Multiple Choice**: "Qual aspecto da disciplina você mais gostou?"
- **Free Text**: "Deixe sugestões para melhorar a disciplina"

## 🧪 Testing Scenarios

### Student Login Flow
1. **Login as Pedro** (`pedro.oliveira@discente.univasf.edu.br`)
2. **View Dashboard**: See 4 enrolled subjects
3. **Available Surveys**: 3 active surveys from enrolled subjects
4. **Survey Status**: See different states (Active, Upcoming, etc.)

### Different Student Experiences
- **Julia**: 3 subjects, different survey combinations
- **Lucas**: 5 subjects, most comprehensive experience
- **Carla**: 3 subjects, including advanced courses
- **Rafael**: 2 subjects, minimal enrollment

### Survey Status Testing
- **Active Surveys**: Ready to be answered
- **Inactive Survey**: Not available for responses
- **Upcoming Survey**: Opens in 2 days (IA course)
- **Past Responses**: Some students have already answered

### Professor Testing
- **Maria Silva**: 2 subjects with active surveys
- **João Santos**: 1 active + 1 inactive survey
- **Ana Costa**: 1 upcoming survey

## 🔄 Re-seeding Data

The seeding function clears and recreates all data each time the server starts. To prevent this:

1. **Comment out** the seed function call in `main.go`:
```go
// seedDatabase(db)  // Comment this line after first run
```

2. **Or move** the seed function to a separate command/script

## 📈 Data Statistics

- **Total Users**: 9 (1 admin + 3 professors + 5 students)
- **Total Subjects**: 6 computer science courses
- **Total Enrollments**: 16 student-subject combinations
- **Total Surveys**: 5 (4 active, 1 inactive)
- **Total Questions**: 13 across all surveys
- **Question Types**: All 4 types represented (NPS, Rating, Multiple Choice, Free Text)
- **Sample Responses**: 3 responses to test the system

## 🎯 Perfect for Testing

This sample data provides:
- **Realistic academic structure** with proper relationships
- **Multiple user types** to test all roles
- **Various survey states** (active, inactive, upcoming)
- **Different question types** to test UI components
- **Diverse student enrollments** for comprehensive testing
- **Sample responses** to test the complete flow

Start the server and login with any of the test credentials to explore the system! 🚀 