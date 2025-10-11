<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { logout } from '$lib/auth';
	import Button from '$lib/components/ui/Button.svelte';

	let { children } = $props();

	let user: any = $state(null);

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
				user = JSON.parse(userData);
				if (!user.id || !user.role) {
					throw new Error('Invalid user data');
				}
			} catch (e) {
				// Invalid user data, clear and redirect
				logout();
			}
		}
	});

	function getRoleDisplayName(role: string) {
		const roleMap: Record<string, string> = {
			student: 'Estudante',
			professor: 'Professor',
			admin: 'Administrador'
		};
		return roleMap[role] || role;
	}
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Navigation Header -->
	<header class="bg-white shadow">
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<h1 class="text-xl font-semibold text-gray-900">
						Sistema de Consulta Discente
					</h1>
				</div>
				<div class="flex items-center space-x-4">
					{#if user}
						<div class="flex items-center space-x-3">
							<div class="text-sm">
								<div class="font-medium text-gray-900">
									{user.first_name} {user.last_name}
								</div>
								<div class="text-gray-500">
									{getRoleDisplayName(user.role)}
								</div>
							</div>
							<Button variant="ghost" size="sm" onclick={logout}>
								Sair
							</Button>
						</div>
					{:else}
						<!-- Loading state -->
						<div class="h-8 w-32 animate-pulse rounded bg-gray-200"></div>
					{/if}
				</div>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		{@render children()}
	</main>
</div> 