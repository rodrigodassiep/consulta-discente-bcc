# Professor Interface Development Plan

## 🎯 **Overview**

Building a comprehensive professor interface for survey creation, management, and response viewing. All backend APIs are already implemented and working - we need to build the frontend interfaces.

### **Core User Stories**
As a professor, I want to:
1. **See my current teaching load** (subjects I'm teaching this semester)
2. **Create surveys** for my subjects with proper scheduling
3. **Manage my surveys** (view all, edit, activate/deactivate)
4. **View survey responses** and basic analytics
5. **Navigate between past and current semesters**

---

## 📋 **Implementation Steps**

### **Step 1: Professor Dashboard Foundation** 
*Status: Not Started*
*Estimated Time: 1-2 days*
*Priority: HIGH - Foundation*

**🎯 Goal:** Create a working professor dashboard that shows current subjects and provides navigation structure.

**What we build:**
- Replace placeholder professor dashboard with functional interface
- Display professor's subjects for current active semester
- Basic navigation structure between sections
- Empty states for when no subjects/surveys exist

**Requirements:**
- Fetch professor's subjects using existing API (`GET /professor/subjects`)
- Show subject cards with: Name, Code, Semester info
- "Create Survey" button for each subject
- Navigation sections: "My Subjects", "All Surveys", "Responses"
- Handle loading and error states

**UI Components Needed:**
- Subject cards with course information
- Navigation menu/tabs
- Loading skeletons
- Empty states

**API Integration:**
- `GET /professor/subjects` - fetch professor's assigned subjects

---

### **Step 2: Survey Creation Form**
*Status: Not Started*
*Estimated Time: 2-3 days*
*Priority: HIGH - Core Functionality*

**🎯 Goal:** Allow professors to create new surveys for their subjects.

**What we build:**
- Survey creation form (modal or page)
- Form validation and submission
- Success/error handling and user feedback
- Redirect to question management after creation

**Requirements:**
- Form fields: Title, Description, Subject (pre-selected), Open Date, Close Date
- Date validation (close date > open date, both in future or present)
- Subject selection (filtered to professor's subjects)
- Active/inactive toggle
- Submit to existing API (`POST /professor/surveys`)

**Form Fields:**
```
Survey Title: [Required text field]
Description:  [Optional textarea]
Subject:      [Dropdown of professor's subjects]
Open Date:    [Date/time picker]
Close Date:   [Date/time picker]
Active:       [Toggle - Active immediately or save as draft]
```

**Validation Rules:**
- Title: Required, max 200 characters
- Description: Optional, max 1000 characters
- Open Date: Cannot be in the past (unless admin override)
- Close Date: Must be after open date
- Subject: Must be one of professor's assigned subjects

**API Integration:**
- `POST /professor/surveys` - create new survey

---

### **Step 3: Question Management**
*Status: Not Started*
*Estimated Time: 2-3 days*
*Priority: HIGH - Core Functionality*

**🎯 Goal:** Allow professors to add and manage questions for their surveys.

**What we build:**
- Add questions interface with all supported question types
- Question form with dynamic fields based on type
- Question preview showing how it appears to students
- Question ordering and basic editing capabilities

**Question Types Support:**
1. **Free Text Questions**
   - Question text + required flag
   - Character limit option

2. **NPS Questions** 
   - Question text + required flag
   - 0-10 scale (automatic)

3. **Rating Questions**
   - Question text + required flag  
   - 1-5 star scale (automatic)

4. **Multiple Choice Questions**
   - Question text + required flag
   - List of options (add/remove options)
   - Single selection

**Requirements:**
- Question type selector (dropdown/radio buttons)
- Dynamic form fields based on question type
- Question ordering (drag-and-drop or up/down arrows)
- Question preview section
- Save individual questions using existing API
- Ability to edit existing questions (before responses exist)

**UI Components Needed:**
- Question type selector
- Dynamic form components for each question type
- Question preview component
- Question list with ordering controls
- Add/Edit question modals

**API Integration:**
- `POST /professor/surveys/:id/questions` - add new question
- Questions are automatically ordered by creation order initially

---

### **Step 4: Survey Management Dashboard**
*Status: Not Started*
*Estimated Time: 2 days*
*Priority: MEDIUM - Management*

**🎯 Goal:** Provide overview and management of all professor's surveys across all semesters.

**What we build:**
- List all professor's surveys (current + historical)
- Survey status indicators and metadata
- Basic survey actions based on current state
- Filtering and search capabilities

**Requirements:**
- Fetch all surveys using existing API (`GET /professor/surveys`)
- Display survey cards/table with key information
- Status indicators with clear visual design
- Action buttons contextual to survey state
- Filter by semester, status, subject
- Search by survey title

**Survey Information Display:**
- Title and description
- Associated subject and semester
- Creation date and last modified
- Open/close date schedule
- Response count (if any)
- Current status

**Survey Status States:**
- **Draft**: Has no questions yet, or not activated
- **Ready**: Has questions, scheduled but not yet open
- **Active**: Currently accepting responses (within date range and active)
- **Closed**: Past close date or manually deactivated
- **Completed**: Closed with responses collected

**Action Buttons (contextual):**
- **Draft**: Edit, Add Questions, Delete, Activate
- **Ready**: Edit Dates, Add Questions, Deactivate, Delete (if no responses)
- **Active**: View Responses, Deactivate, Edit Dates (extend only)
- **Closed**: View Responses, Archive, Duplicate

**API Integration:**
- `GET /professor/surveys` - fetch all professor's surveys

---

### **Step 5: Response Viewing Interface**
*Status: Not Started*
*Estimated Time: 2-3 days*
*Priority: MEDIUM - Analytics*

**🎯 Goal:** Allow professors to view and analyze responses to their surveys.

**What we build:**
- Response overview for individual surveys
- Individual response detail viewing
- Basic filtering and search capabilities
- Response statistics (basic counts and averages)

**Requirements:**
- Fetch responses using existing API (`GET /professor/surveys/:id/responses`)
- Display response list with key metadata
- Click-through to view individual complete responses
- Filter by date range, completion status
- Basic statistics: total responses, completion rate, average ratings

**Response List View:**
- Student identifier (name or anonymous ID)
- Submission date and time
- Completion status (partial vs complete)
- Quick preview of key responses
- Action to view full response

**Individual Response View:**
- All questions with student's answers
- Submission timestamp
- Navigation between responses
- Export individual response option

**Basic Statistics:**
- Total responses vs enrolled students
- Response completion rate
- Average ratings for NPS/rating questions
- Most common multiple choice answers

**API Integration:**
- `GET /professor/surveys/:id/responses` - fetch responses for specific survey
- `GET /professor/responses` - fetch all responses for professor's surveys

---

## 🔧 **Technical Considerations**

### **Shared Components Needed:**
- Survey card component (reusable across different views)
- Question form components (for each question type)
- Date/time picker components
- Loading states and error handling
- Modal/dialog components

### **State Management:**
- Professor's subjects (cache vs refetch)
- Survey list state
- Form state management for complex forms
- Navigation state between sections

### **UX/UI Decisions (To Be Determined):**
- Modal vs separate pages for survey creation
- Inline editing vs separate edit pages
- Question preview style and detail level
- Response viewing layout (table vs cards)
- Mobile responsiveness approach

### **Data Flow:**
- How to handle semester switching
- Survey draft auto-save behavior
- Real-time response count updates
- Navigation between related items (survey → questions → responses)

---

## 🚀 **Implementation Order**

**Recommended sequence:**
1. **Start with Step 1** (Dashboard Foundation) - gives immediate visual progress
2. **Move to Step 2** (Survey Creation) - enables core workflow
3. **Follow with Step 3** (Question Management) - completes survey creation flow
4. **Then Step 4** (Survey Management) - provides oversight capabilities  
5. **Finally Step 5** (Response Viewing) - closes the feedback loop

**Rationale:** This order allows us to build and test the complete survey creation workflow before adding management and analytics features.

---

## ✅ **Definition of Done**

Each step is complete when:
- [ ] Frontend interface is fully functional
- [ ] All API integrations work correctly
- [ ] Error handling and loading states implemented
- [ ] Basic responsive design works
- [ ] User can complete the intended workflow
- [ ] No critical bugs or broken functionality

**Overall professor interface completion:**
- [ ] Professor can log in and see their dashboard
- [ ] Professor can create surveys with questions
- [ ] Professor can manage their surveys
- [ ] Professor can view responses to their surveys
- [ ] All workflows are intuitive and bug-free

---

*Last Updated: Current Date*
*Version: 1.0*
*Next: Begin Step 1 - Professor Dashboard Foundation* 