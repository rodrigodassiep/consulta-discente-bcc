<script>
  import { createEventDispatcher } from 'svelte';
  
  const dispatch = createEventDispatcher();
  
  let email = '';
  let password = '';
  let loading = false;
  let emailError = '';
  let passwordError = '';

  async function login() {
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
      console.log('Login successful:', data);
      dispatch('login', { email, password });
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
      setTimeout(login, 1000);
    }
  }
  
  function handleForgotPassword() {
    dispatch('forgotPassword');
  }
  
  function handleSignUp() {
    dispatch('register');
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-xl shadow-md">
    <div>
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
        Faça login na sua conta
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Ou
        <a
          class="font-medium text-primary hover:text-primary-focus"
          href="/register"
        >
          Cadastre-se
      </a>
      </p>
    </div>
    
    <form class="mt-8 space-y-6" on:submit|preventDefault={handleSubmit}>
      <div class="rounded-md -space-y-px">
        <div class="mb-4">
          <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            id="email"
            name="email"
            type="email"
            autocomplete="email"
            bind:value={email}
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
            placeholder="Email"
          />
          {#if emailError}
            <p class="mt-1 text-sm text-red-600">{emailError}</p>
          {/if}
        </div>
        
        <div class="mb-2">
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Senha</label>
          <input
            id="password"
            name="password"
            type="password"
            autocomplete="current-password"
            bind:value={password}
            class="appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-primary focus:border-primary focus:z-10 sm:text-sm"
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
            class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
          />
          <label for="remember-me" class="ml-2 block text-sm text-gray-900">
            Lembrar-me
          </label>
        </div>

        <div class="text-sm">
          <button
            type="button"
            on:click={handleForgotPassword}
            class="font-medium text-primary hover:text-primary-focus"
          >
            Esqueceu sua senha?
          </button>
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading}
          on:click={login}
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-primary hover:bg-primary-focus focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
        >
          {#if loading}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
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