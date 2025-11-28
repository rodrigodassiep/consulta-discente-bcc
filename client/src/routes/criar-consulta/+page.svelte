<script lang="ts">
	import { env } from '$env/dynamic/public';

	const API_BASE_URL = env.PUBLIC_API_URL || 'http://localhost:3030';

	let nome = '';
	let dataInicio = '';
	let dataFim = '';
	let submitted = false;
	let error = '';

	async function createConsultation() {
		try {
			const response = await fetch(`${API_BASE_URL}/consulta`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ name: nome, startDate: dataInicio, endDate: dataFim })
			});

			if (response.ok) {
				const data = await response.json();
				// Handle successful creation
			} else {
				const errorData = await response.json();
				// Handle error
				console.error('Error creating consultation:', errorData);
			}
		} catch (err) {
			console.error('Error during consultation creation:', err);
		}
	}

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

    createConsultation();
		// Form is valid, you can process the data here
		// alert(
		// 	`Formulário enviado com sucesso!\nNome: ${nome}\nData de Início: ${dataInicio}\nData de Fim: ${dataFim}`
		// );


		// Reset form after successful submission
		nome = '';
		dataInicio = '';
		dataFim = '';
		submitted = false;
	}
</script>

<div class="min-w-lg mx-auto rounded-lg bg-white p-6 shadow-md">
	<h2 class="mb-6 text-2xl font-bold text-gray-800">Criar consulta discente</h2>

	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		<div class="space-y-2">
			<label for="nome" class="block text-sm font-medium text-gray-700">Nome</label>
			<input
				type="text"
				id="nome"
				bind:value={nome}
				required
				class="focus:ring-primary focus:border-primary w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none focus:ring-2"
				placeholder="Digite o nome da consulta"
			/>
		</div>

		<div class="space-y-2">
			<label for="dataInicio" class="block text-sm font-medium text-gray-700">Data de Início</label>
			<input
				type="date"
				id="dataInicio"
				bind:value={dataInicio}
				required
				class="focus:ring-primary focus:border-primary w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none focus:ring-2"
			/>
		</div>

		<div class="space-y-2">
			<label for="dataFim" class="block text-sm font-medium text-gray-700">Data de Fim</label>
			<input
				type="date"
				id="dataFim"
				bind:value={dataFim}
				required
				class="focus:ring-primary focus:border-primary w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:outline-none focus:ring-2"
			/>
		</div>

		{#if error && submitted}
			<div class="rounded border border-red-400 bg-red-100 p-3 text-red-700">
				{error}
			</div>
		{/if}

		<button
			type="submit"
			class="bg-primary hover:bg-primary-focus focus:ring-primary w-full rounded-md border border-transparent px-4 py-2 text-sm font-medium text-white shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2"
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
