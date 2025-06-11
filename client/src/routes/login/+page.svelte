<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { onMount } from 'svelte';

	const dispatch = createEventDispatcher();

	let email = '';
	let password = '';
	let loading = false;
	let emailError = '';
	let passwordError = '';
	let successMessage = '';

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
			if (response.ok) {
				const data = await response.json();
				// Handle successful login
				const user = data;

				// Store user data and ID in localStorage
				localStorage.setItem('user', JSON.stringify(user));
				localStorage.setItem('userId', user.id.toString());

				// Redirect based on user role
				const roleRedirects = {
					student: '/dashboard/student',
					professor: '/dashboard/professor',
					admin: '/dashboard/admin'
				};

				const redirectPath = roleRedirects[user.role as keyof typeof roleRedirects] || '/';
				window.location.href = redirectPath;
			} else {
				const error = await response.json();
				// Handle error
				if (response.status === 401) {
					emailError = 'Email ou senha inválidos';
				} else {
					console.error('Login failed:', error);
				}
			}
		} catch (error) {
			console.log(error);
			emailError = 'Email ou senha inválidos';
			console.error('Error during login:', emailError);
		} finally {
			loading = false;
		}
	}

	function validateEmail(email: string) {
		const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return re.test(email);
	}

	function handleSubmit() {
		// Reset errors
		emailError = '';
		passwordError = '';

		// Validate inputs
		let isValid = true;

		if (!email) {
			emailError = 'Email é obrigatório';
			isValid = false;
		} else if (!validateEmail(email)) {
			emailError = 'Por favor, insira um email válido';
			isValid = false;
		}

		if (!password) {
			passwordError = 'Senha é obrigatório';
			isValid = false;
		} else if (password.length < 6) {
			passwordError = 'Senha deve ter pelo menos 6 caracteres';
			isValid = false;
		}

		if (isValid) {
			loading = true;

			// Simulate API call
			login();
		}
	}

	function handleForgotPassword() {
		dispatch('forgotPassword');
	}

	function handleSignUp() {
		dispatch('register');
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8 rounded-xl bg-white p-8 shadow-md">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
				Faça login na sua conta
			</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Ou
				<a class="text-primary hover:text-primary-focus font-medium" href="/register">
					Cadastre-se
				</a>
			</p>
		</div>
		{#if successMessage}
			<div class="rounded-md bg-green-50 p-4">
				<p class="text-sm text-green-800">{successMessage}</p>
			</div>
		{/if}

		<form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="-space-y-px rounded-md">
				<div class="mb-4">
					<label for="email" class="mb-1 block text-sm font-medium text-gray-700">Email</label>
					<input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						bind:value={email}
						class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
						placeholder="Email"
					/>
					{#if emailError}
						<p class="mt-1 text-sm text-red-600">{emailError}</p>
					{/if}
				</div>

				<div class="mb-2">
					<label for="password" class="mb-1 block text-sm font-medium text-gray-700">Senha</label>
					<input
						id="password"
						name="password"
						type="password"
						autocomplete="current-password"
						bind:value={password}
						class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
						placeholder="Senha"
					/>
					{#if passwordError}
						<p class="mt-1 text-sm text-red-600">{passwordError}</p>
					{/if}
				</div>
			</div>

			<div class="flex items-center justify-between">
				<div class="flex items-center">
					<input
						id="remember-me"
						name="remember-me"
						type="checkbox"
						class="text-primary focus:ring-primary h-4 w-4 rounded border-gray-300"
					/>
					<label for="remember-me" class="ml-2 block text-sm text-gray-900"> Lembrar-me </label>
				</div>

				<div class="text-sm">
					<button
						type="button"
						on:click={handleForgotPassword}
						class="text-primary hover:text-primary-focus font-medium"
					>
						Esqueceu sua senha?
					</button>
				</div>
			</div>

			<div>
				<button
					type="submit"
					disabled={loading}
					class="bg-primary hover:bg-primary-focus focus:ring-primary group relative flex w-full justify-center rounded-md border border-transparent px-4 py-2 text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2"
				>
					{#if loading}
						<svg
							class="-ml-1 mr-3 h-5 w-5 animate-spin text-white"
							xmlns="http://www.w3.org/2000/svg"
							fill="none"
							viewBox="0 0 24 24"
						>
							<circle
								class="opacity-25"
								cx="12"
								cy="12"
								r="10"
								stroke="currentColor"
								stroke-width="4"
							></circle>
							<path
								class="opacity-75"
								fill="currentColor"
								d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
							></path>
						</svg>
						Fazendo login
					{:else}
						Fazer login
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>

<style>
	.text-primary {
		color: var(--color-primary);
	}

	.text-primary-focus {
		color: var(--color-primary-focus);
	}

	.bg-primary {
		background-color: var(--color-primary);
	}

	.bg-primary-focus {
		background-color: var(--color-primary-focus);
	}

	.hover\:text-primary-focus:hover {
		color: var(--color-primary-focus);
	}

	.hover\:bg-primary-focus:hover {
		background-color: var(--color-primary-focus);
	}

	.focus\:ring-primary:focus {
		--tw-ring-color: var(--color-primary);
	}

	.focus\:border-primary:focus {
		border-color: var(--color-primary);
	}
</style>
