# Student Feedback System - Data Models Documentation

## Overview

This system is designed to collect feedback from students about their courses every semester. It allows professors to create surveys with various question types, students to respond to surveys for their enrolled subjects, and administrators to view all responses across the system.

## System Roles

The system supports three main user roles:

- **Student**: Can respond to surveys for subjects they are enrolled in
- **Professor**: Can create surveys for subjects they teach and view responses
- **Admin**: Can view all responses across all subjects and semesters

## Data Models

### 1. User Model

**Purpose**: Represents all users in the system (students, professors, and administrators)

```go
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
```

**Key Features**:
- Email must be unique across the system
- Role is constrained to three values: `student`, `professor`, `admin`
- Database-level validation ensures data integrity

**Relationships**:
- One-to-many with `Subject` (as professor)
- One-to-many with `StudentEnrollment` (as student)
- One-to-many with `Survey` (as professor)
- One-to-many with `Response` (as student)

### 2. Subject Model

**Purpose**: Represents courses/subjects taught in the institution

```go
type Subject struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    Name         string    `json:"name" gorm:"not null"`
    Code         string    `json:"code" gorm:"uniqueIndex;not null"`
    Description  string    `json:"description"`
    ProfessorID  uint      `json:"professor_id" gorm:"not null"`
    Professor    User      `json:"professor" gorm:"foreignKey:ProfessorID;references:ID"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}
```

**Key Features**:
- Each subject has a unique code (e.g., "CS101", "MATH201")
- Each subject is assigned to one professor
- Supports optional descriptions for detailed course information

**Business Logic**:
- Only users with `professor` role can be assigned as subject professors
- Subject codes must be unique to prevent duplicates

### 3. Semester Model

**Purpose**: Represents academic periods when courses are offered

```go
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
```

**Key Features**:
- Naming convention: "YYYY.P" (e.g., "2024.1" for first semester of 2024)
- Only one semester can be active at a time
- Start and end dates define the academic period

**Business Logic**:
- `IsActive` flag helps identify the current semester
- Period typically represents 1st or 2nd semester of the year

### 4. StudentEnrollment Model

**Purpose**: Links students to subjects in specific semesters (handles many-to-many relationships)

```go
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
```

**Key Features**:
- Represents "Student X is enrolled in Subject Y during Semester Z"
- Enables students to see only surveys for their enrolled subjects
- Historical tracking of enrollments across semesters

**Business Logic**:
- Students can only respond to surveys for subjects they are enrolled in
- Enrollments are semester-specific (same student can take same subject in different semesters)

### 5. Survey Model

**Purpose**: Represents feedback forms created by professors for their subjects

```go
type Survey struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Title       string    `json:"title" gorm:"not null"`
    Description string    `json:"description"`
    SubjectID   uint      `json:"subject_id" gorm:"not null"`
    Subject     Subject   `json:"subject" gorm:"foreignKey:SubjectID;references:ID"`
    SemesterID  uint      `json:"semester_id" gorm:"not null"`
    Semester    Semester  `json:"semester" gorm:"foreignKey:SemesterID;references:ID"`
    ProfessorID uint      `json:"professor_id" gorm:"not null"`
    Professor   User      `json:"professor" gorm:"foreignKey:ProfessorID;references:ID"`
    IsActive    bool      `json:"is_active" gorm:"default:true"`
    OpenDate    time.Time `json:"open_date"`
    CloseDate   time.Time `json:"close_date"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    Questions   []Question `json:"questions" gorm:"foreignKey:SurveyID"`
}
```

**Key Features**:
- Each survey belongs to one subject in one semester
- Time-based availability (open/close dates)
- Active/inactive status for manual control
- Contains multiple questions

**Business Logic**:
- Only the professor teaching the subject can create surveys for it
- Students can only see surveys for subjects they're enrolled in
- Survey availability is controlled by both `IsActive` flag and date range

### 6. Question Model

**Purpose**: Individual questions within surveys, supporting multiple question types

```go
type Question struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    SurveyID   uint      `json:"survey_id" gorm:"not null"`
    Survey     Survey    `json:"survey" gorm:"foreignKey:SurveyID;references:ID"`
    Type       string    `json:"type" gorm:"not null;check:type IN ('nps','free_text','rating','multiple_choice')"`
    Text       string    `json:"text" gorm:"not null"`
    Required   bool      `json:"required" gorm:"default:false"`
    Order      int       `json:"order" gorm:"not null"`
    Options    string    `json:"options"` // JSON string for multiple choice options
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

**Supported Question Types**:
- **NPS** (`nps`): Net Promoter Score (0-10 scale)
- **Free Text** (`free_text`): Open-ended text responses
- **Rating** (`rating`): Numeric rating scale (e.g., 1-5 stars)
- **Multiple Choice** (`multiple_choice`): Predefined options

**Key Features**:
- Questions are ordered within surveys
- Required/optional question support
- Flexible options storage as JSON string
- Database-level validation for question types

### 7. Response Model

**Purpose**: Stores student answers to survey questions

```go
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
```

**Key Features**:
- Links student answers to specific questions and surveys
- Automatic timestamp tracking
- Stores all answer types as strings (numbers, text, selections)
- Enables tracking of when responses were submitted

**Business Logic**:
- One response per student per question
- Students can view their historical responses
- Professors can view all responses to their surveys
- Admins can view all responses system-wide

## System Workflow

### 1. Setup Phase
1. Admin creates semesters and marks one as active
2. Admin creates subjects and assigns professors
3. Admin enrolls students in subjects for each semester

### 2. Survey Creation Phase
1. Professor creates survey for their subject in current semester
2. Professor adds questions of various types to the survey
3. Professor sets survey availability dates

### 3. Response Collection Phase
1. Students see available surveys for their enrolled subjects
2. Students submit responses to survey questions
3. System stores responses with timestamps

### 4. Analysis Phase
1. Professors view responses for their surveys
2. Admins view all responses across the system
3. Historical data is maintained for trend analysis

## Database Relationships Summary

- **User** → **Subject** (1:many, as professor)
- **User** → **StudentEnrollment** (1:many, as student)
- **User** → **Survey** (1:many, as professor)
- **User** → **Response** (1:many, as student)
- **Subject** → **StudentEnrollment** (1:many)
- **Subject** → **Survey** (1:many)
- **Semester** → **StudentEnrollment** (1:many)
- **Semester** → **Survey** (1:many)
- **Survey** → **Question** (1:many)
- **Survey** → **Response** (1:many)
- **Question** → **Response** (1:many)

## Constants Reference

### User Roles
```go
const (
    RoleStudent   = "student"
    RoleProfessor = "professor"
    RoleAdmin     = "admin"
)
```

### Question Types
```go
const (
    QuestionTypeNPS      = "nps"
    QuestionTypeFreeText = "free_text"
    QuestionTypeRating   = "rating"
    QuestionTypeChoice   = "multiple_choice"
)
```

This documentation provides a complete reference for understanding the student feedback system's data structure and business logic. 
