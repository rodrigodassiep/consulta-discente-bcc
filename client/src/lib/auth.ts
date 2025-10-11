import { browser } from '$app/environment';

/**
 * Logout user and clean localStorage
 */
export function logout() {
	if (browser) {
		localStorage.removeItem('user');
		localStorage.removeItem('userId');
		localStorage.removeItem('token');
		window.location.href = '/login';
	}
}

/**
 * Check if user is logged in
 */
export function isAuthenticated(): boolean {
	if (!browser) return false;
	const token = localStorage.getItem('token');
	const user = localStorage.getItem('user');
	return !!(token && user);
}

/**
 * Get current user data
 */
export function getCurrentUser() {
	if (!browser) return null;
	const userData = localStorage.getItem('user');
	if (!userData) return null;

	try {
		return JSON.parse(userData);
	} catch {
		// Invalid user data, logout
		logout();
		return null;
	}
}
