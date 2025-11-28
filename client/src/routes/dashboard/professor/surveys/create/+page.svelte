<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	// Form data
	let formData = {
		title: '',
		description: '',
		subject_id: 0,
		semester_id: 3,
		open_date: '',
		close_date: '',
		is_active: true
	};

	// State
	let selectedSubject: any = null;
	let subjects: any[] = [];
	let currentSemester: any = null;
	let loading = true;
	let submitting = false;
	let errors: { [key: string]: string } = {};

	// Load data on mount
	onMount(async () => {
		try {
			// Get subject ID from URL params
			const subjectId = $page.url.searchParams.get('subject');

			// Load professor's subjects and current semester in parallel
			const [subjectsResult, semesterResult] = await Promise.all([
				api.getProfessorSubjects(),
				api.getCurrentSemester()
			]);

			if (!subjectsResult.success) {
				throw new Error(subjectsResult.error || 'Failed to load subjects');
			}

			subjects = (subjectsResult.data as any)?.subjects || [];

			// Get current active semester
			if (semesterResult.success) {
				currentSemester = (semesterResult.data as any)?.semester;
				if (currentSemester) {
					formData.semester_id = currentSemester.id;
				}
			}

			// If subject ID is provided, pre-select it
			if (subjectId) {
				const subject = subjects.find((s) => s.id === parseInt(subjectId));
				if (subject) {
					selectedSubject = subject;
					formData.subject_id = subject.id;
				}
			}

			// Set default dates (open tomorrow, close in 2 weeks)
			setDefaultDates();
		} catch (err) {
			console.error('Failed to initialize:', err);
			// Redirect back to dashboard on error
			goto('/dashboard/professor');
		} finally {
			loading = false;
		}
	});

	function setDefaultDates() {
		const now = new Date();

		// Default open date: tomorrow at 9 AM
		const openDate = new Date(now);
		openDate.setDate(openDate.getDate() + 1);
		openDate.setHours(9, 0, 0, 0);

		// Default close date: 2 weeks from open date at 11:59 PM
		const closeDate = new Date(openDate);
		closeDate.setDate(closeDate.getDate() + 14);
		closeDate.setHours(23, 59, 0, 0);

		// Format for datetime-local input
		formData.open_date = formatDateForInput(openDate);
		formData.close_date = formatDateForInput(closeDate);
	}

	function formatDateForInput(date: Date): string {
		// Format: YYYY-MM-DDTHH:MM
		const year = date.getFullYear();
		const month = String(date.getMonth() + 1).padStart(2, '0');
		const day = String(date.getDate()).padStart(2, '0');
		const hours = String(date.getHours()).padStart(2, '0');
		const minutes = String(date.getMinutes()).padStart(2, '0');

		return `${year}-${month}-${day}T${hours}:${minutes}`;
	}

	function validateForm(): boolean {
		errors = {};

		// Title validation
		if (!formData.title.trim()) {
			errors.title = 'Título é obrigatório';
		} else if (formData.title.length > 200) {
			errors.title = 'Título deve ter no máximo 200 caracteres';
		}

		// Description validation
		if (formData.description.length > 1000) {
			errors.description = 'Descrição deve ter no máximo 1000 caracteres';
		}

		// Subject validation
		if (!formData.subject_id) {
			errors.subject_id = 'Disciplina é obrigatória';
		}

		// Date validation
		if (!formData.open_date) {
			errors.open_date = 'Data de abertura é obrigatória';
		}

		if (!formData.close_date) {
			errors.close_date = 'Data de fechamento é obrigatória';
		}

		if (formData.open_date && formData.close_date) {
			const openDate = new Date(formData.open_date);
			const closeDate = new Date(formData.close_date);

			if (closeDate <= openDate) {
				errors.close_date = 'Data de fechamento deve ser posterior à data de abertura';
			}
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit() {
		if (!validateForm()) {
			return;
		}

		submitting = true;

		try {
			// Prepare survey data for API
			const surveyData = {
				title: formData.title.trim(),
				description: formData.description.trim(),
				subject_id: formData.subject_id,
				semester_id: formData.semester_id,
				open_date: new Date(formData.open_date).toISOString(),
				close_date: new Date(formData.close_date).toISOString(),
				is_active: formData.is_active
			};

			const result = await api.createSurvey(surveyData);

			if (!result.success) {
				throw new Error(result.error || 'Erro ao criar pesquisa');
			}

			const createdSurvey = (result.data as any)?.survey;

			// Success! Navigate to question management
			if (createdSurvey?.id) {
				// Navigate to question management for the new survey
				goto(`/dashboard/professor/surveys/${createdSurvey.id}/questions`);
			} else {
				goto('/dashboard/professor');
			}
		} catch (err) {
			console.error('Failed to create survey:', err);
			errors.submit = err instanceof Error ? err.message : 'Erro ao criar pesquisa';
		} finally {
			submitting = false;
		}
	}

	function handleCancel() {
		goto('/dashboard/professor');
	}

	function handleSubjectChange(event: Event) {
		const select = event.target as HTMLSelectElement;
		const subjectId = parseInt(select.value);

		formData.subject_id = subjectId;

		// Update selected subject (semester_id already set from current active semester)
		selectedSubject = subjects.find((s) => s.id === subjectId) || null;
	}
</script>

<svelte:head>
	<title>Criar Nova Pesquisa - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center space-x-4">
		<Button variant="outline" onclick={handleCancel} size="sm">
			<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"
				></path>
			</svg>
			Voltar
		</Button>
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Criar Nova Pesquisa</h1>
			<p class="mt-1 text-gray-600">Configure uma nova pesquisa de feedback para seus alunos</p>
		</div>
	</div>

	<div class="mx-auto max-w-2xl space-y-6">

		<!-- Loading State -->
		{#if loading}
			<Card>
				<div class="animate-pulse space-y-4">
					<div class="h-4 w-1/4 rounded bg-gray-200"></div>
					<div class="h-10 w-full rounded bg-gray-200"></div>
					<div class="h-4 w-1/4 rounded bg-gray-200"></div>
					<div class="h-20 w-full rounded bg-gray-200"></div>
					<div class="flex space-x-4">
						<div class="h-10 w-24 rounded bg-gray-200"></div>
						<div class="h-10 w-24 rounded bg-gray-200"></div>
					</div>
				</div>
			</Card>
		{:else}
			<!-- Survey Creation Form -->
			<form
				onsubmit={(e) => {
					e.preventDefault();
					handleSubmit();
				}}
			>
				<div class="space-y-6">
					<!-- Basic Information -->
					<Card>
						<div class="space-y-4">
							<h2 class="text-lg font-semibold text-gray-900">Informações Básicas</h2>

							<!-- Title -->
							<div>
								<label for="title" class="mb-1 block text-sm font-medium text-gray-700">
									Título da Pesquisa *
								</label>
								<input
									id="title"
									type="text"
									bind:value={formData.title}
									placeholder="Ex: Avaliação do curso de Algoritmos I"
									maxlength="200"
									class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none
										{errors.title ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
								/>
								{#if errors.title}
									<p class="mt-1 text-sm text-red-600">{errors.title}</p>
								{/if}
								<p class="mt-1 text-xs text-gray-500">{formData.title.length}/200 caracteres</p>
							</div>

							<!-- Description -->
							<div>
								<label for="description" class="mb-1 block text-sm font-medium text-gray-700">
									Descrição (opcional)
								</label>
								<textarea
									id="description"
									bind:value={formData.description}
									placeholder="Descreva o objetivo da pesquisa e instruções adicionais..."
									rows="3"
									maxlength="1000"
									class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none
										{errors.description ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
								></textarea>
								{#if errors.description}
									<p class="mt-1 text-sm text-red-600">{errors.description}</p>
								{/if}
								<p class="mt-1 text-xs text-gray-500">
									{formData.description.length}/1000 caracteres
								</p>
							</div>

							<!-- Subject Selection -->
							<div>
								<label for="subject" class="mb-1 block text-sm font-medium text-gray-700">
									Disciplina *
								</label>
								<select
									id="subject"
									bind:value={formData.subject_id}
									onchange={handleSubjectChange}
									class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none
										{errors.subject_id ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
								>
									<option value={0}>Selecione uma disciplina</option>
									{#each subjects as subject}
										<option value={subject.id}>
											{subject.name} ({subject.code})
										</option>
									{/each}
								</select>
								{#if errors.subject_id}
									<p class="mt-1 text-sm text-red-600">{errors.subject_id}</p>
								{/if}

								<!-- Show selected subject info -->
								{#if selectedSubject}
									<div class="mt-2 rounded-md bg-blue-50 p-3">
										<div class="flex items-center space-x-2">
											<Badge variant="secondary">{selectedSubject.code}</Badge>
											<span class="text-sm font-medium text-blue-900">{selectedSubject.name}</span>
										</div>
										{#if currentSemester}
											<p class="mt-1 text-xs text-blue-700">
												Semestre: {currentSemester.name}
											</p>
										{/if}
									</div>
								{/if}
							</div>
						</div>
					</Card>

					<!-- Schedule -->
					<Card>
						<div class="space-y-4">
							<h2 class="text-lg font-semibold text-gray-900">Cronograma</h2>

							<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
								<!-- Open Date -->
								<div>
									<label for="open_date" class="mb-1 block text-sm font-medium text-gray-700">
										Data de Abertura *
									</label>
									<input
										id="open_date"
										type="datetime-local"
										bind:value={formData.open_date}
										class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none
											{errors.open_date ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
									/>
									{#if errors.open_date}
										<p class="mt-1 text-sm text-red-600">{errors.open_date}</p>
									{/if}
								</div>

								<!-- Close Date -->
								<div>
									<label for="close_date" class="mb-1 block text-sm font-medium text-gray-700">
										Data de Fechamento *
									</label>
									<input
										id="close_date"
										type="datetime-local"
										bind:value={formData.close_date}
										class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none
											{errors.close_date ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''}"
									/>
									{#if errors.close_date}
										<p class="mt-1 text-sm text-red-600">{errors.close_date}</p>
									{/if}
								</div>
							</div>

							<!-- Status -->
							<div>
								<label class="flex items-center space-x-3">
									<input
										type="checkbox"
										bind:checked={formData.is_active}
										class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
									/>
									<span class="text-sm font-medium text-gray-700"
										>Ativar pesquisa imediatamente</span
									>
								</label>
								<p class="mt-1 text-xs text-gray-500">
									Se desmarcado, a pesquisa será criada como rascunho
								</p>
							</div>
						</div>
					</Card>

					<!-- Submit Error -->
					{#if errors.submit}
						<Card class="border-red-200 bg-red-50">
							<p class="text-red-800">{errors.submit}</p>
						</Card>
					{/if}

					<!-- Actions -->
					<Card>
						<div class="flex items-center justify-between">
							<Button variant="outline" onclick={handleCancel} disabled={submitting}>
								Cancelar
							</Button>
							<Button type="submit" disabled={submitting}>
								{submitting ? 'Criando...' : 'Criar Pesquisa'}
							</Button>
						</div>
					</Card>
				</div>
			</form>
		{/if}
	</div>
</div>
