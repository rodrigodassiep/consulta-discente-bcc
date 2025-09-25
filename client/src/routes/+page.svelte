<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';


	onMount(() => {
		if (browser) {
			// Check if user is logged in
			const userData = localStorage.getItem('user');
			if (userData) {
				try {
					const user = JSON.parse(userData);
					// Redirect to appropriate dashboard
					const roleRedirects = {
						student: '/dashboard/student',
						professor: '/dashboard/professor',
						admin: '/dashboard/admin'
					};
					
					const redirectPath = roleRedirects[user.role as keyof typeof roleRedirects];
					if (redirectPath) {
						window.location.href = redirectPath;
						return;
					}
				} catch (e) {
					// Invalid user data, clear it
					localStorage.removeItem('user');
					localStorage.removeItem('userId');
					localStorage.removeItem('token');
				}
			}
			
			// User not logged in, redirect to login
			window.location.href = '/login';
		}
	});
</script>

<svelte:head>
	<title>Sistema de Consulta Discente</title>
	<meta name="description" content="Sistema de feedback para estudantes" />
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gray-50">
	<div class="text-center">
		<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
		<p class="mt-4 text-gray-600">Redirecionando...</p>
	</div>
</div>


