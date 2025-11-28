# Feature Proposals for Large-Scale Use

This document outlines future features that would make the Student Feedback System suitable for a large Computer Science program at a Brazilian university.

## 1. Institution & Course Integration
- **University SSO integration:** Support OIDC/SAML login with the university identity provider (RA/USP number), replacing local credentials while keeping role-based access (student/professor/admin).
- **Enrollment sync:** Import subjects, semesters and enrollments from the institutional system (or via CSV upload) so that `StudentEnrollment` stays aligned with official records.
- **Multi-department support:** Add department/faculty fields to `Subject` and `User`, allowing different institutes to share the same instance with scoped admin permissions.

## 2. Student Experience & Pedagogy
- **Localization & accessibility:** First-class PT-BR UI (with optional EN), clear mobile-friendly layouts and accessibility improvements (keyboard navigation, ARIA labels).
- **Survey campaigns:** Allow admins to create university-wide campaigns (e.g., “Avaliação de Docentes 2025.1”) with progress tracking by course, semester and department.
- **Incentive mechanisms:** Non-intrusive features such as completion badges or summary dashboards showing how individual feedback contributes to course improvement.

## 3. Professor & TA Tools
- **Question banks & templates:** Let professors reuse question sets per course type (e.g., “Algoritmos 1”, “Laboratório”) and share templates inside departments.
- **TA role:** Introduce a teaching assistant role with limited access to survey creation and analytics for specific subjects.
- **Comparative dashboards:** Visualize trends across semesters (e.g., “Algoritmos 1 – 2023.2 vs 2024.1”) with breakdowns by question type and class section.

## 4. Analytics & Research Support
- **Export for data analysis:** Provide anonymized exports (CSV/Parquet) for use in Python/R, enabling students and faculty to run analyses in courses like Data Science and Statistics.
- **Built-in analytics:** Add dashboards with response distributions, heatmaps by question and subject, and filters by semester, course level and modality (presencial/EAD).
- **Research datasets:** Allow admins to generate long-term, LGPD-compliant datasets that can be reused in undergraduate projects and academic research.

## 5. Governance, Privacy & LGPD
- **LGPD-aware consent:** Explicit consent screens, clear privacy notices and configurable data retention periods per campaign.
- **Anonymity guarantees:** Enforce minimum response thresholds before showing results; aggregate outputs so that individual students cannot be re-identified.
- **Audit & access logs:** Track who accessed which reports and exports, supporting internal audits and institutional governance requirements.

