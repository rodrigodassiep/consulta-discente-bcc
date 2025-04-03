<script>
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	let firstName = '';
	let lastName = '';
	let email = '';
	let password = '';
	let confirmPassword = '';
	let loading = false;

	let firstNameError = '';
	let lastNameError = '';
	let emailError = '';
	let passwordError = '';
	let confirmPasswordError = '';

	async function register() {
		const response = await fetch('http://localhost:3030/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ firstName, lastName, email, password })
		});
		if (response.ok) {
			const data = await response.json();
			// Handle successful login
			console.log('Login successful:', data);
			dispatch('login', { firstName, lastName, email, password });
		} else {
			const error = await response.json();
			// Handle error
			console.error('Login failed:', error);
		}
	}

	function validateEmail(email) {
		const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		return re.test(email);
	}

	function handleSubmit() {
		// Reset errors
		firstNameError = '';
		lastNameError = '';
		emailError = '';
		passwordError = '';
		confirmPasswordError = '';

		// Validate inputs
		let isValid = true;

		if (!firstName) {
			firstNameError = 'Nome é obrigatório';
			isValid = false;
		}

		if (!lastName) {
			lastNameError = 'Sobrenome é obrigatório';
			isValid = false;
		}

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

		if (!confirmPassword) {
			confirmPasswordError = 'Por favor, confirme sua senha';
			isValid = false;
		} else if (password !== confirmPassword) {
			confirmPasswordError = 'Senhas não correspondem';
			isValid = false;
		}

		if (isValid) {
			loading = true;

			// Simulate API call
			setTimeout(() => {
				dispatch('signup', { firstName, lastName, email, password });
				loading = false;
			}, 1000);
		}
	}

	function handleLogin() {
		dispatch('login');
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-md space-y-8 rounded-xl bg-white p-8 shadow-md">
		<div>
			<h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Crie a sua conta</h2>
			<p class="mt-2 text-center text-sm text-gray-600">
				Ou
				<a class="text-primary hover:text-primary-focus font-medium" href="/login">
					faça login na sua conta
        </a>
			</p>
		</div>

		<form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
			<div class="space-y-4 rounded-md">
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<label for="first-name" class="mb-1 block text-sm font-medium text-gray-700">Nome</label
						>
						<input
							id="first-name"
							name="first-name"
							type="text"
							bind:value={firstName}
							class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
							placeholder="Nome"
						/>
						{#if firstNameError}
							<p class="mt-1 text-sm text-red-600">{firstNameError}</p>
						{/if}
					</div>

					<div>
						<label for="last-name" class="mb-1 block text-sm font-medium text-gray-700"
							>Sobrenome</label
						>
						<input
							id="last-name"
							name="last-name"
							type="text"
							bind:value={lastName}
							class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
							placeholder="Sobrenome"
						/>
						{#if lastNameError}
							<p class="mt-1 text-sm text-red-600">{lastNameError}</p>
						{/if}
					</div>
				</div>

				<div>
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

				<div>
					<label for="password" class="mb-1 block text-sm font-medium text-gray-700">Senha</label>
					<input
						id="password"
						name="password"
						type="password"
						bind:value={password}
						class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
						placeholder="Senha"
					/>
					{#if passwordError}
						<p class="mt-1 text-sm text-red-600">{passwordError}</p>
					{/if}
				</div>

				<div>
					<label for="confirm-password" class="mb-1 block text-sm font-medium text-gray-700"
						>Confirme sua senha</label
					>
					<input
						id="confirm-password"
						name="confirm-password"
						type="password"
						bind:value={confirmPassword}
						class="focus:ring-primary focus:border-primary relative block w-full appearance-none rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-500 focus:z-10 focus:outline-none sm:text-sm"
						placeholder="Confirme sua senha"
					/>
					{#if confirmPasswordError}
						<p class="mt-1 text-sm text-red-600">{confirmPasswordError}</p>
					{/if}
				</div>
			</div>

			<div class="flex items-center">
				<input
					id="terms"
					name="terms"
					type="checkbox"
					class="text-primary focus:ring-primary h-4 w-4 rounded border-gray-300"
					required
				/>
				<label for="terms" class="ml-2 block text-sm text-gray-900">
					Eu aceito os <a href="#" class="text-primary hover:text-primary-focus"
						>Termos de Serviço</a
					>
					ou a <a href="#" class="text-primary hover:text-primary-focus">Política de Privacidade</a>
				</label>
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
						Criando a conta...
					{:else}
						Criar conta
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
