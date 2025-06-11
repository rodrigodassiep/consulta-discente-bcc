# Student Feedback System - Frontend

## Overview

This is the frontend for the Student Feedback System built with SvelteKit and Tailwind CSS. It provides role-based dashboards for students, professors, and administrators.

## Features Implemented

### ✅ **Shared Components Library**
- **Button**: Reusable button with variants (primary, secondary, outline, ghost)
- **Card**: Container component for content sections
- **Badge**: Status indicators with color variants
- **Layout**: Consistent page layout with navigation

### ✅ **Authentication Flow**
- **Login page**: Updated to redirect to role-based dashboards
- **Route protection**: Dashboard routes require authentication
- **Role-based routing**: Users redirected to appropriate dashboards

### ✅ **Student Dashboard** (`/dashboard/student`)
- **Enrolled subjects**: Display all subjects student is enrolled in
- **Available surveys**: Show surveys for enrolled subjects with status indicators
- **Survey status**: Active, Inactive, Upcoming, Expired states
- **Response tracking**: Shows which surveys have been answered
- **Recent activity**: List of recent survey responses

## File Structure

```
src/
├── lib/
│   ├── components/
│   │   ├── ui/
│   │   │   ├── Button.svelte
│   │   │   ├── Card.svelte
│   │   │   └── Badge.svelte
│   │   ├── Layout.svelte
│   │   └── index.ts
│   └── api.ts
├── routes/
│   ├── dashboard/
│   │   ├── student/
│   │   │   └── +page.svelte
│   │   └── +layout.svelte
│   ├── login/
│   │   └── +page.svelte
│   └── +page.svelte
```

## API Integration

The frontend uses a centralized API client (`/lib/api.ts`) that:
- Handles authentication headers (X-User-ID)
- Provides type-safe responses
- Manages error handling
- Supports all backend endpoints

## Authentication System

### How it works:
1. User logs in via `/login`
2. User data and ID stored in localStorage
3. Login redirects to role-based dashboard
4. Protected routes check for valid authentication
5. API requests include user ID header

### User Roles:
- **Student**: Access to survey participation and response history
- **Professor**: Survey creation and response viewing (not yet implemented)
- **Admin**: System-wide management (not yet implemented)

## Student Dashboard Features

### 📋 **Enrolled Subjects Section**
- Grid layout showing all enrolled subjects
- Subject name, code, professor, and semester info
- Responsive design (3 columns on large screens)

### 📊 **Available Surveys Section**
- List view with survey cards
- Color-coded status indicators:
  - 🟢 Green: Active surveys
  - 🟡 Yellow: Upcoming surveys
  - 🔴 Red: Expired surveys
  - ⚫ Gray: Inactive surveys
- Survey information includes:
  - Title and description
  - Subject and code
  - Open/close dates
  - Number of questions
  - Response status

### 📈 **Recent Activity Section**
- Shows last 5 survey responses
- Response date and survey information
- "View all responses" option for more history

## UI Design

- **Modern and clean**: Uses Tailwind CSS for styling
- **Responsive**: Mobile-first design approach
- **Accessible**: Proper color contrast and focus states
- **Loading states**: Spinner animations during data fetching
- **Error handling**: User-friendly error messages with retry options

## Next Steps

To complete the frontend, we still need to implement:

1. **Survey Taking Page** (`/surveys/[id]`)
   - Display survey questions
   - Handle different question types (NPS, free text, rating, multiple choice)
   - Submit responses

2. **Professor Dashboard** (`/dashboard/professor`)
   - Survey creation and management
   - Response viewing and analytics

3. **Admin Dashboard** (`/dashboard/admin`)
   - User management
   - Subject and semester management
   - System-wide analytics

4. **Additional Features**
   - Survey response editing
   - Advanced filtering and search
   - Data export functionality
   - Real-time notifications

## Development

To run the frontend:

```bash
cd client
npm run dev
```

The frontend connects to the backend API running on `http://localhost:3030`.

## Dependencies

- **SvelteKit**: Full-stack Svelte framework
- **Tailwind CSS**: Utility-first CSS framework
- **TypeScript**: Type safety and better development experience 