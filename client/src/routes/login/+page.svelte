<script lang="ts">
	import { onMount } from 'svelte';

	let email = '';
	let password = '';
	let loading = false;
	let emailError = '';
	let passwordError = '';
	let loginError = '';
	let successMessage = '';
	let showPassword = false;

	onMount(() => {
		// Check if user was redirected from registration
		const urlParams = new URLSearchParams(window.location.search);
		if (urlParams.get('registered') === 'true') {
			successMessage = 'Conta criada com sucesso! Faça login para continuar.';
			// Clean up URL
			window.history.replaceState({}, document.title, window.location.pathname);
		}
	});

	async function login() {
		try {
			const response = await fetch('http://localhost:3030/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password })
			});

			const data = await response.json();

			if (response.ok) {
				let token, user;

				if (data.token && data.user) {
					token = data.token;
					user = data.user;
				} else {
					user = data;
				}

				if (token) {
					localStorage.setItem('token', token);
				}
				localStorage.setItem('user', JSON.stringify(user));
				localStorage.setItem('userId', user.id.toString());

				const roleRedirects = {
					student: '/dashboard/student',
					professor: '/dashboard/professor',
					admin: '/dashboard/admin'
				};

				const redirectPath = roleRedirects[user.role as keyof typeof roleRedirects] || '/';
				window.location.href = redirectPath;
			} else {
				if (response.status === 401) {
					loginError = 'Email ou senha incorretos. Verifique suas credenciais e tente novamente.';
				} else {
					loginError = data.error || 'Erro no servidor. Tente novamente mais tarde.';
				}
			}
		} catch (error) {
			console.error('Network error during login:', error);
			loginError = 'Erro de conexão. Verifique sua internet e tente novamente.';
		} finally {
			loading = false;
		}
	}

	function validateEmail(email: string) {
		const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return re.test(email);
	}

	function handleSubmit() {
		emailError = '';
		passwordError = '';
		loginError = '';

		let isValid = true;

		if (!email) {
			emailError = 'Email é obrigatório';
			isValid = false;
		} else if (!validateEmail(email)) {
			emailError = 'Por favor, insira um email válido';
			isValid = false;
		}

		if (!password) {
			passwordError = 'Senha é obrigatória';
			isValid = false;
		} else if (password.length < 6) {
			passwordError = 'Senha deve ter pelo menos 6 caracteres';
			isValid = false;
		}

		if (isValid) {
			loading = true;
			login();
		}
	}
</script>

<svelte:head>
	<title>Login - Sistema de Consulta Discente</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-50 to-blue-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md">
		<!-- Logo/Brand Section -->
		<div class="mb-8 text-center">
			<div class="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-xl bg-gradient-to-br from-blue-700 to-blue-800 shadow-lg">
				<svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
				</svg>
			</div>
			<h1 class="text-2xl font-bold text-gray-900">Sistema de Consulta Discente</h1>
			<p class="mt-1 text-sm text-gray-500">Bacharelado em Ciência da Computação - IME-USP</p>
		</div>

		<!-- Login Card -->
		<div class="rounded-xl border border-gray-200 bg-white p-8 shadow-lg">
			<div class="mb-6">
				<h2 class="text-xl font-semibold text-gray-900">Bem-vindo de volta</h2>
				<p class="mt-1 text-sm text-gray-500">Entre com suas credenciais para acessar o sistema</p>
			</div>

			<!-- Success Message -->
			{#if successMessage}
				<div class="mb-6 flex items-start gap-3 rounded-lg border border-green-200 bg-green-50 p-4">
					<svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
					<p class="text-sm text-green-800">{successMessage}</p>
				</div>
			{/if}

			<!-- Error Message -->
			{#if loginError}
				<div class="mb-6 flex items-start gap-3 rounded-lg border border-red-200 bg-red-50 p-4">
					<svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
					<div>
						<p class="text-sm font-medium text-red-800">Falha no login</p>
						<p class="mt-0.5 text-sm text-red-700">{loginError}</p>
					</div>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="space-y-5">
				<!-- Email Field -->
				<div>
					<label for="email" class="mb-1.5 block text-sm font-medium text-gray-700">
						Email
					</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"></path>
							</svg>
						</div>
						<input
							id="email"
							name="email"
							type="email"
							autocomplete="email"
							bind:value={email}
							class="block w-full rounded-lg border py-2.5 pl-10 pr-3 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {emailError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
							placeholder="seu@email.com"
						/>
					</div>
					{#if emailError}
						<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
							<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
							</svg>
							{emailError}
						</p>
					{/if}
				</div>

				<!-- Password Field -->
				<div>
					<label for="password" class="mb-1.5 block text-sm font-medium text-gray-700">
						Senha
					</label>
					<div class="relative">
						<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
							<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
							</svg>
						</div>
						<input
							id="password"
							name="password"
							type={showPassword ? 'text' : 'password'}
							autocomplete="current-password"
							bind:value={password}
							class="block w-full rounded-lg border py-2.5 pl-10 pr-10 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {passwordError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
							placeholder="••••••••"
						/>
						<button
							type="button"
							class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
							on:click={() => showPassword = !showPassword}
						>
							{#if showPassword}
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
								</svg>
							{:else}
								<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
								</svg>
							{/if}
						</button>
					</div>
					{#if passwordError}
						<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
							<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
							</svg>
							{passwordError}
						</p>
					{/if}
				</div>

				<!-- Remember Me & Forgot Password -->
				<div class="flex items-center justify-between">
					<label class="flex cursor-pointer items-center gap-2">
						<input
							type="checkbox"
							class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
						/>
						<span class="text-sm text-gray-600">Lembrar-me</span>
					</label>
					<a href="/forgot-password" class="text-sm font-medium text-blue-700 hover:text-blue-800">
						Esqueceu a senha?
					</a>
				</div>

				<!-- Submit Button -->
				<button
					type="submit"
					disabled={loading}
					class="flex w-full items-center justify-center gap-2 rounded-lg bg-gradient-to-br from-blue-700 to-blue-800 px-4 py-3 text-sm font-semibold text-white shadow-sm transition-all hover:from-blue-800 hover:to-blue-900 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none"
				>
					{#if loading}
						<svg class="h-5 w-5 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
							<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
							<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
						</svg>
						Entrando...
					{:else}
						Entrar
					{/if}
				</button>
			</form>

			<!-- Divider -->
			<div class="relative my-6">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-gray-200"></div>
				</div>
				<div class="relative flex justify-center text-sm">
					<span class="bg-white px-4 text-gray-500">Novo no sistema?</span>
				</div>
			</div>

			<!-- Register Link -->
			<a
				href="/register"
				class="flex w-full items-center justify-center gap-2 rounded-lg border-2 border-blue-700 bg-white px-4 py-2.5 text-sm font-semibold text-blue-700 transition-colors hover:bg-blue-50"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
				</svg>
				Criar uma conta
			</a>
		</div>

		<!-- Footer -->
		<p class="mt-6 text-center text-xs text-gray-500">
			&copy; {new Date().getFullYear()} IME-USP. Todos os direitos reservados.
		</p>
	</div>
</div>
