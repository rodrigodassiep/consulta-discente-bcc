<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { logout, getCurrentUser } from '$lib/auth';
	import Button from '$lib/components/ui/Button.svelte';
	import Card from '$lib/components/ui/Card.svelte';

	let isLoggedIn = false;
	let user: any = null;

	onMount(() => {
		if (browser) {
			// Check if user is logged in
			const userData = localStorage.getItem('user');
			if (userData) {
				try {
					user = JSON.parse(userData);
					isLoggedIn = true;
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
					isLoggedIn = false;
					user = null;
					// Clear invalid data
					if (browser) {
						localStorage.removeItem('user');
						localStorage.removeItem('userId');
						localStorage.removeItem('token');
					}
				}
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

<svelte:head>
	<title>Sistema de Consulta Discente - BCC IME-USP</title>
	<meta name="description" content="Sistema de avaliação e feedback para estudantes do Bacharelado em Ciência da Computação - IME-USP" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
	<!-- Header -->
	<header class="border-b border-gray-200 bg-white/80 backdrop-blur-sm">
		<nav class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<h1 class="text-xl font-bold text-gray-900 sm:text-2xl">
						Consulta Discente BCC
					</h1>
				</div>
				<div class="flex items-center gap-4">
					{#if isLoggedIn && user}
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
						<a href="/login">
							<Button variant="primary" size="md">
								Fazer Login
							</Button>
						</a>
					{/if}
				</div>
			</div>
		</nav>
	</header>

	<!-- Hero Section -->
	<main class="mx-auto max-w-7xl px-4 py-16 sm:px-6 sm:py-24 lg:px-8">
		<div class="text-center">
			<h2 class="text-4xl font-extrabold tracking-tight text-gray-900 sm:text-5xl md:text-6xl">
				<span class="block">Avalie suas disciplinas</span>
				<span class="block text-blue-600">Melhore o ensino</span>
			</h2>
			<p class="mx-auto mt-6 max-w-2xl text-lg text-gray-600 sm:text-xl">
				Sistema de consulta discente do Bacharelado em Ciência da Computação do IME-USP.
				Compartilhe sua experiência e ajude a aprimorar o curso através de feedback estruturado.
			</p>
			<div class="mt-10 flex justify-center gap-4">
				<a href="/login">
					<Button variant="primary" size="lg">
						Começar Agora
					</Button>
				</a>
				<a href="/register">
					<Button variant="outline" size="lg">
						Criar Conta
					</Button>
				</a>
			</div>
		</div>

		<!-- Features Section -->
		<div class="mt-24 grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
			<Card>
				<div class="text-center">
					<div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-blue-100">
						<svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
						</svg>
					</div>
					<h3 class="mb-2 text-lg font-semibold text-gray-900">Avaliações Estruturadas</h3>
					<p class="text-sm text-gray-600">
						Responda questionários organizados sobre disciplinas e professores de forma rápida e objetiva.
					</p>
				</div>
			</Card>

			<Card>
				<div class="text-center">
					<div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-100">
						<svg class="h-6 w-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
						</svg>
					</div>
					<h3 class="mb-2 text-lg font-semibold text-gray-900">Feedback Anônimo</h3>
					<p class="text-sm text-gray-600">
						Suas respostas são confidenciais, permitindo avaliações honestas e construtivas.
					</p>
				</div>
			</Card>

			<Card>
				<div class="text-center">
					<div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-purple-100">
						<svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
						</svg>
					</div>
					<h3 class="mb-2 text-lg font-semibold text-gray-900">Melhoria Contínua</h3>
					<p class="text-sm text-gray-600">
						Contribua para o aprimoramento do curso e ajude futuras turmas a terem uma experiência melhor.
					</p>
				</div>
			</Card>
		</div>
	</main>

	<!-- Footer -->
	<footer class="border-t border-gray-200 bg-white">
		<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
			<p class="text-center text-sm text-gray-500">
				&copy; {new Date().getFullYear()} Sistema de Consulta Discente - BCC IME-USP.
			</p>
		</div>
	</footer>
</div>

