<script lang="ts">
	import { onMount } from 'svelte';
	import { logout } from '$lib/auth';
	import Button from './ui/Button.svelte';

	interface Props {
		title: string;
		children: any;
	}

	let { title, children }: Props = $props();

	let user: any = null;
	let loading = true;

	onMount(() => {
		// Get user from localStorage or session
		const userData = localStorage.getItem('user');
		if (userData) {
			user = JSON.parse(userData);
		}
		loading = false;
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

<svelte:head>
	<title>{title} - Sistema de Consulta Discente</title>
</svelte:head>

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
					{#if user && !loading}
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
					{/if}
				</div>
			</div>
		</div>
	</header>

	<!-- Main Content -->
	<main class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
		<div class="mb-8">
			<h2 class="text-3xl font-bold text-gray-900">{title}</h2>
		</div>
		{@render children()}
	</main>
</div> 