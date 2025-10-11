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

<div class="min-h-screen bg-gradient-to-br from-blue-50 via-indigo-50 to-white">
	<!-- Header -->
	<header class="border-b border-blue-100 bg-white/90 backdrop-blur-md shadow-sm">
		<nav class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-20 items-center justify-between">
				<div class="flex items-center gap-3">
					<!-- Logo/Icon -->
					<div>
						<h1 class="text-xl font-bold text-gray-900 sm:text-2xl" style="font-family: var(--font-display);">
							Consulta Discente <span class="text-blue-700">BCC</span>
						</h1>
						<p class="hidden text-xs text-gray-500 sm:block">IME-USP</p>
					</div>
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
	<main class="mx-auto max-w-7xl px-4 py-20 sm:px-6 sm:py-32 lg:px-8">
		<div class="text-center">
			<!-- Badge -->
			<div class="mb-8 inline-flex items-center gap-2 rounded-full border border-blue-200 bg-blue-50 px-4 py-2 text-sm font-medium text-blue-800">
				<svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
					<path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0zM6 18a1 1 0 001-1v-2.065a8.935 8.935 0 00-2-.712V17a1 1 0 001 1z" />
				</svg>
				Sistema Acadêmico IME-USP
			</div>

			<h2 class="text-5xl font-bold tracking-tight text-gray-900 sm:text-6xl md:text-7xl" style="font-family: var(--font-display); letter-spacing: -0.03em;">
				<span class="block">Avalie suas disciplinas</span>
				<span class="mt-2 block bg-gradient-to-r from-blue-700 via-blue-600 to-indigo-600 bg-clip-text text-transparent">
					Melhore o ensino
				</span>
			</h2>

			<p class="mx-auto mt-8 max-w-2xl text-xl leading-relaxed text-gray-600 sm:text-2xl">
				Compartilhe sua experiência e contribua para o aprimoramento do
				<span class="font-semibold text-gray-900">Bacharelado em Ciência da Computação</span> através de feedback estruturado e anônimo.
			</p>

			<div class="mt-12 flex flex-col justify-center gap-4 sm:flex-row">
				<a href="/login">
					<Button variant="primary" size="lg">
						<svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
						</svg>
						Começar Agora
					</Button>
				</a>
				<a href="/register">
					<Button variant="outline" size="lg">
						Criar Conta Gratuita
					</Button>
				</a>
			</div>

			<!-- Stats -->
			<div class="mt-16 grid grid-cols-3 gap-8 border-y border-gray-200 bg-white/50 py-8 backdrop-blur-sm">
				<div>
					<div class="text-3xl font-bold text-blue-700">100%</div>
					<div class="mt-1 text-sm text-gray-600">Anônimo</div>
				</div>
				<div>
					<div class="text-3xl font-bold text-blue-700">5 min</div>
					<div class="mt-1 text-sm text-gray-600">Por Avaliação</div>
				</div>
				<div>
					<div class="text-3xl font-bold text-blue-700">24/7</div>
					<div class="mt-1 text-sm text-gray-600">Disponível</div>
				</div>
			</div>
		</div>

		<!-- Features Section -->
		<div class="mt-32">
			<div class="mb-12 text-center">
				<h3 class="text-3xl font-bold text-gray-900" style="font-family: var(--font-display);">
					Como funciona
				</h3>
				<p class="mt-4 text-lg text-gray-600">Sistema simples, seguro e eficaz</p>
			</div>

			<div class="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
				<!-- Feature 1 -->
				<div class="group relative overflow-hidden rounded-2xl border border-blue-100 bg-white p-8 shadow-sm transition-all hover:shadow-xl hover:border-blue-200">
					<div class="absolute right-0 top-0 h-32 w-32 translate-x-8 -translate-y-8 rounded-full bg-blue-100 opacity-50 transition-transform group-hover:scale-150"></div>
					<div class="relative">
						<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-xl bg-gradient-to-br from-blue-600 to-blue-800 shadow-lg">
							<svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
							</svg>
						</div>
						<h3 class="mb-3 text-xl font-bold text-gray-900" style="font-family: var(--font-display);">
							Avaliações Estruturadas
						</h3>
						<p class="text-gray-600 leading-relaxed">
							Questionários organizados sobre disciplinas e professores com perguntas objetivas e personalizadas.
						</p>
					</div>
				</div>

				<!-- Feature 2 -->
				<div class="group relative overflow-hidden rounded-2xl border border-indigo-100 bg-white p-8 shadow-sm transition-all hover:shadow-xl hover:border-indigo-200">
					<div class="absolute right-0 top-0 h-32 w-32 translate-x-8 -translate-y-8 rounded-full bg-indigo-100 opacity-50 transition-transform group-hover:scale-150"></div>
					<div class="relative">
						<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-600 to-indigo-800 shadow-lg">
							<svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
							</svg>
						</div>
						<h3 class="mb-3 text-xl font-bold text-gray-900" style="font-family: var(--font-display);">
							100% Anônimo
						</h3>
						<p class="text-gray-600 leading-relaxed">
							Suas respostas são completamente confidenciais, garantindo feedback honesto e construtivo.
						</p>
					</div>
				</div>

				<!-- Feature 3 -->
				<div class="group relative overflow-hidden rounded-2xl border border-amber-100 bg-white p-8 shadow-sm transition-all hover:shadow-xl hover:border-amber-200">
					<div class="absolute right-0 top-0 h-32 w-32 translate-x-8 -translate-y-8 rounded-full bg-amber-100 opacity-50 transition-transform group-hover:scale-150"></div>
					<div class="relative">
						<div class="mb-6 flex h-16 w-16 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-amber-700 shadow-lg">
							<svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
							</svg>
						</div>
						<h3 class="mb-3 text-xl font-bold text-gray-900" style="font-family: var(--font-display);">
							Impacto Real
						</h3>
						<p class="text-gray-600 leading-relaxed">
							Contribua para melhorias no curso e ajude futuras turmas a terem uma experiência ainda melhor.
						</p>
					</div>
				</div>
			</div>
		</div>
	</main>

	<!-- Footer -->
	<footer class="mt-32 border-t border-gray-200 bg-gradient-to-br from-gray-50 to-white">
		<div class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
			<div class="flex flex-col items-center justify-between gap-4 sm:flex-row">
				<div class="flex items-center gap-2">
					<span class="font-semibold text-gray-900">Consulta Discente BCC</span>
				</div>
				<p class="text-sm text-gray-600">
					&copy; {new Date().getFullYear()} IME-USP. Sistema desenvolvido para melhorar a qualidade do ensino.
				</p>
			</div>
		</div>
	</footer>
</div>

