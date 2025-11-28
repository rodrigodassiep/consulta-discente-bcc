<script lang="ts">
	let firstName = '';
	let lastName = '';
	let email = '';
	let password = '';
	let confirmPassword = '';
	let loading = false;
	let desiredRole: 'student' | 'professor' | 'admin' = 'student';
	let acceptedTerms = false;

	let firstNameError = '';
	let lastNameError = '';
	let emailError = '';
	let passwordError = '';
	let confirmPasswordError = '';
	let roleError = '';
	let registerError = '';

	let showPassword = false;
	let showConfirmPassword = false;

	async function register() {
		try {
			const response = await fetch('http://localhost:3030/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					first_name: firstName,
					last_name: lastName,
					email,
					password,
					role: 'student',
					requested_role: desiredRole
				})
			});

			if (response.ok) {
				window.location.href = '/login?registered=true';
			} else {
				const error = await response.json();

				if (response.status === 409) {
					emailError = 'Este email já está cadastrado';
				} else if (error.error) {
					if (error.error.includes('First name')) {
						firstNameError = error.error;
					} else if (error.error.includes('Last name')) {
						lastNameError = error.error;
					} else if (error.error.includes('Email')) {
						emailError = error.error;
					} else if (error.error.includes('Password')) {
						passwordError = error.error;
					} else {
						registerError = error.error;
					}
				} else {
					registerError = 'Erro ao criar conta. Tente novamente.';
				}
			}
		} catch (error) {
			console.error('Network error:', error);
			registerError = 'Erro de conexão. Verifique sua internet e tente novamente.';
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
		firstNameError = '';
		lastNameError = '';
		emailError = '';
		passwordError = '';
		confirmPasswordError = '';
		roleError = '';
		registerError = '';

		let isValid = true;

		if (!firstName.trim()) {
			firstNameError = 'Nome é obrigatório';
			isValid = false;
		}

		if (!lastName.trim()) {
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
			passwordError = 'Senha é obrigatória';
			isValid = false;
		} else if (password.length < 6) {
			passwordError = 'Senha deve ter pelo menos 6 caracteres';
			isValid = false;
		}

		if (!confirmPassword) {
			confirmPasswordError = 'Por favor, confirme sua senha';
			isValid = false;
		} else if (password !== confirmPassword) {
			confirmPasswordError = 'As senhas não correspondem';
			isValid = false;
		}

		if (!desiredRole) {
			roleError = 'Selecione um tipo de acesso';
			isValid = false;
		}

		if (isValid) {
			loading = true;
			register();
		}
	}

	const roleOptions = [
		{
			value: 'student',
			label: 'Estudante',
			description: 'Avaliar disciplinas e professores',
			icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>`
		},
		{
			value: 'professor',
			label: 'Professor',
			description: 'Criar e gerenciar pesquisas',
			icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"></path>`
		},
		{
			value: 'admin',
			label: 'Administrador',
			description: 'Gerenciar todo o sistema',
			icon: `<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>`
		}
	];
</script>

<svelte:head>
	<title>Criar Conta - Sistema de Consulta Discente</title>
</svelte:head>

<div class="flex min-h-screen items-center justify-center bg-gradient-to-br from-gray-50 to-blue-50 px-4 py-12 sm:px-6 lg:px-8">
	<div class="w-full max-w-lg">
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

		<!-- Register Card -->
		<div class="rounded-xl border border-gray-200 bg-white p-8 shadow-lg">
			<div class="mb-6">
				<h2 class="text-xl font-semibold text-gray-900">Crie sua conta</h2>
				<p class="mt-1 text-sm text-gray-500">Preencha os dados abaixo para se cadastrar no sistema</p>
			</div>

			<!-- Error Message -->
			{#if registerError}
				<div class="mb-6 flex items-start gap-3 rounded-lg border border-red-200 bg-red-50 p-4">
					<svg class="mt-0.5 h-5 w-5 flex-shrink-0 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
					</svg>
					<div>
						<p class="text-sm font-medium text-red-800">Erro no cadastro</p>
						<p class="mt-0.5 text-sm text-red-700">{registerError}</p>
					</div>
				</div>
			{/if}

			<form on:submit|preventDefault={handleSubmit} class="space-y-5">
				<!-- Name Fields -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<label for="first-name" class="mb-1.5 block text-sm font-medium text-gray-700">
							Nome
						</label>
						<div class="relative">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
								<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
								</svg>
							</div>
							<input
								id="first-name"
								name="first-name"
								type="text"
								bind:value={firstName}
								class="block w-full rounded-lg border py-2.5 pl-10 pr-3 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {firstNameError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
								placeholder="Seu nome"
							/>
						</div>
						{#if firstNameError}
							<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
								<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
								</svg>
								{firstNameError}
							</p>
						{/if}
					</div>

					<div>
						<label for="last-name" class="mb-1.5 block text-sm font-medium text-gray-700">
							Sobrenome
						</label>
						<div class="relative">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
								<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
								</svg>
							</div>
							<input
								id="last-name"
								name="last-name"
								type="text"
								bind:value={lastName}
								class="block w-full rounded-lg border py-2.5 pl-10 pr-3 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {lastNameError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
								placeholder="Seu sobrenome"
							/>
						</div>
						{#if lastNameError}
							<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
								<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
								</svg>
								{lastNameError}
							</p>
						{/if}
					</div>
				</div>

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

				<!-- Role Selection -->
				<div>
					<label class="mb-2 block text-sm font-medium text-gray-700">
						Tipo de acesso desejado
					</label>
					<div class="grid grid-cols-1 gap-3 sm:grid-cols-3">
						{#each roleOptions as role}
							<label class="relative cursor-pointer">
								<input
									type="radio"
									name="role"
									value={role.value}
									bind:group={desiredRole}
									class="peer sr-only"
								/>
								<div class="rounded-lg border-2 p-3 text-center transition-all peer-checked:border-blue-600 peer-checked:bg-blue-50 {desiredRole === role.value ? 'border-blue-600 bg-blue-50' : 'border-gray-200 hover:border-gray-300 hover:bg-gray-50'}">
									<svg class="mx-auto h-6 w-6 {desiredRole === role.value ? 'text-blue-600' : 'text-gray-400'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
										{@html role.icon}
									</svg>
									<p class="mt-1 text-sm font-medium {desiredRole === role.value ? 'text-blue-900' : 'text-gray-900'}">{role.label}</p>
									<p class="mt-0.5 text-xs {desiredRole === role.value ? 'text-blue-700' : 'text-gray-500'}">{role.description}</p>
								</div>
							</label>
						{/each}
					</div>
					{#if desiredRole !== 'student'}
						<div class="mt-2 flex items-start gap-2 rounded-lg bg-amber-50 p-3">
							<svg class="mt-0.5 h-4 w-4 flex-shrink-0 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
							</svg>
							<p class="text-xs text-amber-800">
								Pedidos de acesso como {desiredRole === 'professor' ? 'professor' : 'administrador'} serão revisados. Até a aprovação, sua conta terá acesso como estudante.
							</p>
						</div>
					{/if}
					{#if roleError}
						<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
							<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
							</svg>
							{roleError}
						</p>
					{/if}
				</div>

				<!-- Password Fields -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
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
								bind:value={password}
								class="block w-full rounded-lg border py-2.5 pl-10 pr-10 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {passwordError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
								placeholder="Min. 6 caracteres"
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

					<div>
						<label for="confirm-password" class="mb-1.5 block text-sm font-medium text-gray-700">
							Confirmar senha
						</label>
						<div class="relative">
							<div class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
								<svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path>
								</svg>
							</div>
							<input
								id="confirm-password"
								name="confirm-password"
								type={showConfirmPassword ? 'text' : 'password'}
								bind:value={confirmPassword}
								class="block w-full rounded-lg border py-2.5 pl-10 pr-10 text-gray-900 placeholder-gray-400 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0 sm:text-sm {confirmPasswordError ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : 'border-gray-300 focus:border-blue-500 focus:ring-blue-200'}"
								placeholder="Repita a senha"
							/>
							<button
								type="button"
								class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
								on:click={() => showConfirmPassword = !showConfirmPassword}
							>
								{#if showConfirmPassword}
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
						{#if confirmPasswordError}
							<p class="mt-1.5 flex items-center gap-1 text-sm text-red-600">
								<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01"></path>
								</svg>
								{confirmPasswordError}
							</p>
						{/if}
					</div>
				</div>

				<!-- Terms -->
				<div class="flex items-start gap-3">
					<input
						id="terms"
						name="terms"
						type="checkbox"
						bind:checked={acceptedTerms}
						required
						class="mt-1 h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
					/>
					<label for="terms" class="text-sm text-gray-600">
						Eu li e aceito os Termos de Uso e a Política de Privacidade
					</label>
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
						Criando conta...
					{:else}
						<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"></path>
						</svg>
						Criar conta
					{/if}
				</button>
			</form>

			<!-- Divider -->
			<div class="relative my-6">
				<div class="absolute inset-0 flex items-center">
					<div class="w-full border-t border-gray-200"></div>
				</div>
				<div class="relative flex justify-center text-sm">
					<span class="bg-white px-4 text-gray-500">Já tem uma conta?</span>
				</div>
			</div>

			<!-- Login Link -->
			<a
				href="/login"
				class="flex w-full items-center justify-center gap-2 rounded-lg border-2 border-blue-700 bg-white px-4 py-2.5 text-sm font-semibold text-blue-700 transition-colors hover:bg-blue-50"
			>
				<svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path>
				</svg>
				Fazer login
			</a>
		</div>

		<!-- Footer -->
		<p class="mt-6 text-center text-xs text-gray-500">
			&copy; {new Date().getFullYear()} IME-USP. Todos os direitos reservados.
		</p>
	</div>
</div>
