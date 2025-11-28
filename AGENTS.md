# Repository Guidelines

This repository implements the Student Feedback System for IME-USP, with a Go backend and a SvelteKit + TypeScript frontend.

## Project Structure & Modules
- `server/`: Go API using Gin, GORM, PostgreSQL/SQLite; tests live in `*_test.go`.
- `client/`: SvelteKit app with TypeScript and Tailwind; routes under `client/src/routes`, shared components under `client/src/lib`.
- Additional docs: `README.md`, `client/FRONTEND_README.md`, `server/TESTING_README.md`, `server/MODEL_DOCUMENTATION.md`.

## Build, Test & Development
- Backend dev: `cd server && go run main.go` (API on port 3030).
- Backend tests: `cd server && go test ./...` (add `-coverprofile=coverage.out` for coverage).
- Frontend dev: `cd client && npm install` once, then `npm run dev` (SvelteKit on port 5173).
- Frontend build: `cd client && npm run build`.
- Frontend tests: `cd client && npm run test`.

## Coding Style & Naming
- Go: format with `go fmt ./...`; use `CamelCase` for exported types/functions and `camelCase` for locals; keep HTTP handlers small and reuse model helpers.
- Svelte/TS: run `npm run format` and `npm run lint` before committing; 2-space indentation; PascalCase for components, `camelCase` for variables/functions, `SCREAMING_SNAKE_CASE` for constants.
- Prefer clear, domain-focused names (e.g., `StudentEnrollment`, `SurveyResponse`) and avoid duplicating business rules across client and server.

## Testing Guidelines
- Backend: follow patterns in `server/models_test.go`, `server/api_test.go`, `server/middleware_test.go`, `server/seed_test.go`; use `testify/assert`; keep tests fast and deterministic (SQLite in-memory).
- Frontend: place Vitest specs near code (e.g., `*.spec.ts` in `src`); use Testing Library for user-centric tests; update tests and sample data when changing API contracts.

## Commit & Pull Request Guidelines
- Write clear, imperative commit messages; prefer prefixes like `feat`, `fix`, `chore` and optional scopes, e.g., `feat(auth): implement JWT login`.
- Group related changes into small PRs; include backend and frontend changes together when they depend on each other.
- In PR descriptions, include purpose, main changes, test commands, and screenshots/GIFs for UI updates; link related issues or roadmap items when applicable.

## Agent-Specific Notes
- Respect this `AGENTS.md` file when editing; align with existing architecture and testing patterns, and avoid large backend or routing refactors without coordinating with maintainers.

