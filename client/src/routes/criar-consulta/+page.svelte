<script lang="ts">
  let nome = '';
  let dataInicio = '';
  let dataFim = '';
  let submitted = false;
  let error = '';

  function handleSubmit() {
    submitted = true;
    
    if (!nome || !dataInicio || !dataFim) {
      error = 'Por favor, preencha todos os campos obrigatórios.';
      return;
    }
    
    // Validate that end date is after start date
    if (new Date(dataFim) < new Date(dataInicio)) {
      error = 'A data de fim deve ser posterior à data de início.';
      return;
    }
    
    error = '';
    
    // Form is valid, you can process the data here
    alert(`Formulário enviado com sucesso!\nNome: ${nome}\nData de Início: ${dataInicio}\nData de Fim: ${dataFim}`);
    
    // Reset form after successful submission
    nome = '';
    dataInicio = '';
    dataFim = '';
    submitted = false;
  }
</script>

<div class="max-w-md mx-auto p-6 bg-white rounded-lg shadow-md">
  <h2 class="text-2xl font-bold mb-6 text-gray-800">Formulário de Período</h2>
  
  <form on:submit|preventDefault={handleSubmit} class="space-y-4">
    <div class="space-y-2">
      <label for="nome" class="block text-sm font-medium text-gray-700">Nome</label>
      <input
        type="text"
        id="nome"
        bind:value={nome}
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
        placeholder="Digite seu nome"
      />
    </div>
    
    <div class="space-y-2">
      <label for="dataInicio" class="block text-sm font-medium text-gray-700">Data de Início</label>
      <input
        type="date"
        id="dataInicio"
        bind:value={dataInicio}
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
      />
    </div>
    
    <div class="space-y-2">
      <label for="dataFim" class="block text-sm font-medium text-gray-700">Data de Fim</label>
      <input
        type="date"
        id="dataFim"
        bind:value={dataFim}
        required
        class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-primary"
      />
    </div>
    
    {#if error && submitted}
      <div class="p-3 bg-red-100 border border-red-400 text-red-700 rounded">
        {error}
      </div>
    {/if}
    
    <button
      type="submit"
      class="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary hover:bg-primary-focus focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
    >
      Enviar
    </button>
  </form>
</div>

<style>
  /* Use the CSS variables from the layout */
  :global(:root) {
    --color-primary: #234aac;
    --color-primary-focus: #2a88c6;
  }
  
  /* Add utility classes for the primary colors */
  :global(.bg-primary) {
    background-color: var(--color-primary);
  }
  
  :global(.bg-primary-focus) {
    background-color: var(--color-primary-focus);
  }
  
  :global(.focus\:ring-primary:focus) {
    --tw-ring-color: var(--color-primary);
  }
  
  :global(.focus\:border-primary:focus) {
    border-color: var(--color-primary);
  }
  
  :global(.focus\:ring-primary-focus:focus) {
    --tw-ring-color: var(--color-primary-focus);
  }
  
  :global(.focus\:ring-offset-2:focus) {
    --tw-ring-offset-width: 2px;
  }
</style>