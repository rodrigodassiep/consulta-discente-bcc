<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';

	let { children } = $props();

	onMount(() => {
		if (browser) {
			const userData = localStorage.getItem('user');
			const userId = localStorage.getItem('userId');
			
			if (!userData || !userId) {
				// Not logged in, redirect to login
				window.location.href = '/login';
				return;
			}

			try {
				const user = JSON.parse(userData);
				if (!user.id || !user.role) {
					throw new Error('Invalid user data');
				}
			} catch (e) {
				// Invalid user data, clear and redirect
				localStorage.removeItem('user');
				localStorage.removeItem('userId');
				localStorage.removeItem('token');
				window.location.href = '/login';
			}
		}
	});
</script>

{@render children()} 