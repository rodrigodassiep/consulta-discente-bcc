# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **student feedback system** for the Bachelor in Computer Science program at IME-USP. The system enables students to evaluate professors and courses through structured surveys with role-based access control.

## Development Commands

### Backend (Go Server)
```bash
cd server
go run main.go                    # Start development server on port 3030
go test ./...                     # Run all tests
go test -coverprofile=coverage.out && go tool cover -html=coverage.out  # Generate and view coverage
go mod tidy                       # Update dependencies
```

### Frontend (SvelteKit)
```bash
cd client
npm run dev                       # Start development server on port 5173
npm run build                     # Build for production
npm run test                      # Run all tests
npm run test:unit                 # Run unit tests only
npm run format                    # Format with Prettier
npm run lint                      # Lint with ESLint
```

## Architecture Overview

### Tech Stack
- **Backend**: Go with Gin framework, GORM ORM, PostgreSQL/SQLite
- **Frontend**: SvelteKit with TypeScript, Tailwind CSS
- **Testing**: Testify (Go), Vitest + Testing Library (Svelte)

### Database Models & Relationships
The system uses 7 core entities with complex relationships:

1. **User** (student/professor/admin roles)
2. **Subject** (courses with assigned professors)
3. **Semester** (academic periods like "2024.1")
4. **StudentEnrollment** (student-subject-semester associations)
5. **Survey** (feedback forms with open/close scheduling)
6. **Question** (NPS, free_text, rating, multiple_choice types)
7. **Response** (student answers)

Key relationships: User→Subject (as professor), User→StudentEnrollment (as student), Survey→Question→Response chains.

### Project Structure
```
server/
├── main.go              # Single-file backend (650+ lines)
├── *_test.go           # Comprehensive test suite
└── seed_data.go        # Database seeding

client/
├── src/lib/components/ui/  # Reusable UI components
├── src/lib/api.ts         # Centralized API client
└── src/routes/            # SvelteKit file-based routing
    ├── dashboard/         # Role-based dashboards
    ├── login/register/    # Authentication
    └── surveys/           # Survey management
```

## API & Authentication

### Server Configuration
- **Port**: 3030 (hardcoded)
- **CORS**: Configured for localhost:5173
- **Auth**: Simple header-based (`X-User-ID`) - no JWT
- **Database**: Auto-migration with fallback table reset

### Route Structure
- Public: `/`, `/quote`, `/current-semester`
- Auth: `/register`, `/login`
- Student: `/student/*` (surveys, enrollments)
- Professor: `/professor/*` (survey creation, responses)
- Admin: `/admin/*` (system management)

## Development Patterns

### Backend (main.go)
- Single-file architecture with all logic in main.go
- GORM for database operations with automatic migrations
- Gin middleware for CORS and role-based authentication
- Comprehensive error handling and JSON responses
- Test-driven development with model/middleware/API/seed tests

### Frontend (SvelteKit)
- Component-based architecture with shared UI library
- Centralized API client with TypeScript types
- Role-based routing and dashboard segregation
- Local storage for authentication state management
- Tailwind CSS utility-first styling

## Testing Infrastructure

### Backend Testing
- **Model tests**: Database validation, relationships, constraints
- **Middleware tests**: CORS, authentication, authorization
- **API tests**: Endpoint functionality, validation, errors
- **Seed tests**: Data integrity and seeding behavior
- Uses in-memory SQLite for test isolation

### Frontend Testing
- Vitest with Testing Library for component testing
- JSDOM for browser environment simulation
- Separate client/server test workspace configuration

## Key Configuration

- **Environment**: `/server/.env` for database connection
- **Dependencies**: `/server/go.mod` (Go), `/client/package.json` (npm)
- **Build**: `/client/vite.config.ts` for frontend build configuration
- **Database**: PostgreSQL (production), SQLite (testing)

## Special Development Notes

1. **Database Reset**: System automatically drops/recreates tables on migration failures
2. **Seed Data**: Comprehensive sample data for all models and relationships available via `seed_data.go`
3. **Authentication Flow**: Uses simple header-based auth with role validation middleware
4. **CORS Setup**: Pre-configured for frontend development server communication
5. **Single File Backend**: All server logic consolidated in `main.go` for simplicity

## Current Status

The system has a complete backend with comprehensive testing, functional authentication, student dashboards, and survey status tracking. Frontend components and API integration are established. Key remaining work includes survey-taking interface, professor/admin dashboards, and response analytics.