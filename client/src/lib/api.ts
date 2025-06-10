const API_BASE_URL = 'http://localhost:3030';

interface ApiResponse<T> {
	success: boolean;
	data?: T;
	error?: string;
}

class ApiClient {
	private getHeaders(): HeadersInit {
		const headers: HeadersInit = {
			'Content-Type': 'application/json',
		};

		// Add user ID header if available
		const userId = localStorage.getItem('userId');
		if (userId) {
			headers['X-User-ID'] = userId;
		}

		return headers;
	}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {}
	): Promise<ApiResponse<T>> {
		try {
			console.log(`${API_BASE_URL}${endpoint}`);
			const response = await fetch(`${API_BASE_URL}${endpoint}`, {
				...options,
				headers: {
					...this.getHeaders(),
					...options.headers,
				},
			});

			const data = await response.json();

			if (!response.ok) {
				return {
					success: false,
					error: data.error || `HTTP error! status: ${response.status}`,
				};
			}

			return {
				success: true,
				data,
			};
		} catch (error) {
			return {
				success: false,
				error: error instanceof Error ? error.message : 'Network error',
			};
		}
	}

	// Auth endpoints
	async login(email: string, password: string) {
		return this.request('/login', {
			method: 'POST',
			body: JSON.stringify({ email, password }),
		});
	}

	async register(userData: any) {
		return this.request('/register', {
			method: 'POST',
			body: JSON.stringify(userData),
		});
	}

	// Student endpoints
	async getStudentSubjects() {
		return this.request('/student/subjects');
	}

	async getStudentSurveys() {
		return this.request('/student/surveys');
	}

	async submitResponse(response: any) {
		return this.request('/student/responses', {
			method: 'POST',
			body: JSON.stringify(response),
		});
	}

	async getStudentResponses() {
		return this.request('/student/responses');
	}

	async getSurveyById(surveyId: string) {
		return this.request(`/student/surveys/${surveyId}`);
	}

	// Professor endpoints
	async getProfessorSubjects() {
		return this.request('/professor/subjects');
	}

	async getProfessorSurveys() {
		return this.request('/professor/surveys');
	}

	async createSurvey(survey: any) {
		return this.request('/professor/surveys', {
			method: 'POST',
			body: JSON.stringify(survey),
		});
	}

	async addQuestionToSurvey(surveyId: string, question: any) {
		return this.request(`/professor/surveys/${surveyId}/questions`, {
			method: 'POST',
			body: JSON.stringify(question),
		});
	}

	async getProfessorResponses() {
		return this.request('/professor/responses');
	}

	async getSurveyResponses(surveyId: string) {
		return this.request(`/professor/surveys/${surveyId}/responses`);
	}

	// Admin endpoints
	async createSemester(semester: any) {
		return this.request('/admin/semesters', {
			method: 'POST',
			body: JSON.stringify(semester),
		});
	}

	async getSemesters() {
		return this.request('/admin/semesters');
	}

	async activateSemester(semesterId: string) {
		return this.request(`/admin/semesters/${semesterId}/activate`, {
			method: 'PUT',
		});
	}

	async createSubject(subject: any) {
		return this.request('/admin/subjects', {
			method: 'POST',
			body: JSON.stringify(subject),
		});
	}

	async getSubjects() {
		return this.request('/admin/subjects');
	}

	async createEnrollment(enrollment: any) {
		return this.request('/admin/enrollments', {
			method: 'POST',
			body: JSON.stringify(enrollment),
		});
	}

	async getEnrollments() {
		return this.request('/admin/enrollments');
	}

	async getAllResponses() {
		return this.request('/admin/responses');
	}

	async getAllUsers() {
		return this.request('/admin/users');
	}
}

export const api = new ApiClient();
export type { ApiResponse }; 