<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	// State
	let survey: any = null;
	let questions: any[] = [];
	let loading = true;
	let submitting = false;
	let error = '';

	// Question form state
	let showAddForm = false;
	let editingQuestion: any = null; // Track question being edited
	let questionForm = {
		type: '',
		text: '',
		required: false,
		options: [''] // For multiple choice
	};

	// Question type options
	const questionTypes = [
		{ value: 'free_text', label: 'Texto Livre', description: 'Resposta em texto aberto' },
		{ value: 'nps', label: 'NPS (0-10)', description: 'Escala de recomendação de 0 a 10' },
		{ value: 'rating', label: 'Avaliação (1-5)', description: 'Avaliação em estrelas de 1 a 5' },
		{
			value: 'multiple_choice',
			label: 'Múltipla Escolha',
			description: 'Seleção entre opções predefinidas'
		}
	];

	onMount(async () => {
		const surveyId = $page.params.id;

		if (!surveyId) {
			goto('/dashboard/professor');
			return;
		}

		await loadSurveyAndQuestions(surveyId);
	});

	async function loadSurveyAndQuestions(surveyId: string) {
		try {
			loading = true;

			// Get professor's surveys to find this one
			const surveysResult = await api.getProfessorSurveys();

			if (!surveysResult.success) {
				throw new Error(surveysResult.error || 'Failed to load surveys');
			}

			const surveys = (surveysResult.data as any)?.surveys || [];
			survey = surveys.find((s: any) => s.id === parseInt(surveyId));

			if (!survey) {
				throw new Error('Survey not found or access denied');
			}

			// Questions are included in the survey data
			questions = survey.questions || [];
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error loading survey';
			console.error('Failed to load survey:', err);
		} finally {
			loading = false;
		}
	}

	function resetQuestionForm() {
		questionForm = {
			type: '',
			text: '',
			required: false,
			options: ['']
		};
	}

	function addOption() {
		questionForm.options = [...questionForm.options, ''];
	}

	function removeOption(index: number) {
		if (questionForm.options.length > 1) {
			questionForm.options = questionForm.options.filter((_, i) => i !== index);
		}
	}

	function updateOption(index: number, value: string) {
		questionForm.options = questionForm.options.map((option, i) => (i === index ? value : option));
	}

	async function saveQuestion() {
		// Validate form
		if (!questionForm.text.trim()) {
			error = 'Question text is required';
			return;
		}

		if (questionForm.type === 'multiple_choice') {
			const validOptions = questionForm.options.filter((opt) => opt.trim() !== '');
			if (validOptions.length < 2) {
				error = 'Multiple choice questions need at least 2 options';
				return;
			}
		}

		submitting = true;
		error = '';

		try {
			// Prepare question data
			const questionData = {
				type: questionForm.type,
				text: questionForm.text.trim(),
				required: questionForm.required,
				order: questions.length + 1 // Add at the end
			};

			// Add options for multiple choice
			if (questionForm.type === 'multiple_choice') {
				const validOptions = questionForm.options.filter((opt) => opt.trim() !== '');
				(questionData as any).options = JSON.stringify(validOptions);
			}

			const result = await api.addQuestionToSurvey(survey.id.toString(), questionData);

			if (!result.success) {
				throw new Error(result.error || 'Failed to save question');
			}

			// Reload questions
			await loadSurveyAndQuestions(survey.id.toString());

			// Reset form
			resetQuestionForm();
			showAddForm = false;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error saving question';
		} finally {
			submitting = false;
		}
	}

	function cancelAddQuestion() {
		resetQuestionForm();
		showAddForm = false;
		editingQuestion = null;
		error = '';
	}

	function startEditQuestion(question: any) {
		editingQuestion = question;
		questionForm = {
			type: question.type,
			text: question.text,
			required: question.required,
			options: question.type === 'multiple_choice' && question.options
				? JSON.parse(question.options)
				: ['']
		};
		showAddForm = true;
		error = '';
	}

	async function updateQuestion() {
		// Validate form
		if (!questionForm.text.trim()) {
			error = 'O texto da questão é obrigatório';
			return;
		}

		if (questionForm.type === 'multiple_choice') {
			const validOptions = questionForm.options.filter((opt) => opt.trim() !== '');
			if (validOptions.length < 2) {
				error = 'Questões de múltipla escolha precisam de pelo menos 2 opções';
				return;
			}
		}

		submitting = true;
		error = '';

		try {
			const questionData: any = {
				type: questionForm.type,
				text: questionForm.text.trim(),
				required: questionForm.required,
				order: editingQuestion.order
			};

			if (questionForm.type === 'multiple_choice') {
				const validOptions = questionForm.options.filter((opt) => opt.trim() !== '');
				questionData.options = JSON.stringify(validOptions);
			}

			const result = await api.updateQuestion(
				survey.id.toString(),
				editingQuestion.id.toString(),
				questionData
			);

			if (!result.success) {
				throw new Error(result.error || 'Falha ao atualizar questão');
			}

			// Reload questions
			await loadSurveyAndQuestions(survey.id.toString());

			// Reset form
			resetQuestionForm();
			showAddForm = false;
			editingQuestion = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Erro ao atualizar questão';
		} finally {
			submitting = false;
		}
	}

	async function deleteQuestion(question: any) {
		if (!confirm(`Tem certeza que deseja excluir a questão "${question.text}"?`)) {
			return;
		}

		try {
			const result = await api.deleteQuestion(
				survey.id.toString(),
				question.id.toString()
			);

			if (!result.success) {
				throw new Error(result.error || 'Falha ao excluir questão');
			}

			// Reload questions
			await loadSurveyAndQuestions(survey.id.toString());
		} catch (err) {
			error = err instanceof Error ? err.message : 'Erro ao excluir questão';
		}
	}

	function getQuestionTypeLabel(type: string): string {
		const questionType = questionTypes.find((qt) => qt.value === type);
		return questionType ? questionType.label : type;
	}

	function handleBackToSurvey() {
		goto('/dashboard/professor');
	}
</script>

<svelte:head>
	<title>Gerenciar Questões - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="flex items-center space-x-4">
			<Button variant="outline" onclick={handleBackToSurvey} size="sm">
				<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"
					></path>
				</svg>
				Voltar
			</Button>
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Gerenciar Questões</h1>
				{#if survey}
					<p class="mt-1 text-gray-600">{survey.title}</p>
				{/if}
			</div>
		</div>
	</div>

	<div class="mx-auto max-w-4xl space-y-6">

		<!-- Loading State -->
		{#if loading}
			<Card>
				<div class="animate-pulse space-y-4">
					<div class="h-4 w-1/4 rounded bg-gray-200"></div>
					<div class="h-20 w-full rounded bg-gray-200"></div>
					<div class="h-4 w-1/4 rounded bg-gray-200"></div>
					<div class="h-20 w-full rounded bg-gray-200"></div>
				</div>
			</Card>

			<!-- Error State -->
		{:else if error && !survey}
			<Card class="border-red-200 bg-red-50">
				<div class="py-8 text-center">
					<h3 class="mb-2 text-lg font-medium text-red-800">Erro</h3>
					<p class="mb-4 text-red-600">{error}</p>
					<Button onclick={handleBackToSurvey} variant="outline">Voltar ao Dashboard</Button>
				</div>
			</Card>
		{:else if survey}
			<!-- Survey Info -->
			<Card class="border-blue-200 bg-blue-50">
				<div class="flex items-center justify-between">
					<div>
						<h2 class="text-lg font-semibold text-blue-900">{survey.title}</h2>
						<p class="text-sm text-blue-700">
							{survey.subject?.name} ({survey.subject?.code}) • {questions.length} questões
						</p>
					</div>
					<Badge variant="secondary">
						{survey.is_active ? 'Ativa' : 'Inativa'}
					</Badge>
				</div>
			</Card>

			<!-- Questions List -->
			{#if questions.length > 0}
				<Card>
					<h2 class="mb-4 text-lg font-semibold text-gray-900">Questões Existentes</h2>
					<div class="space-y-4">
						{#each questions as question, index}
							<div class="rounded-lg border border-gray-200 bg-gray-50 p-4">
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<div class="mb-2 flex items-center space-x-2">
											<span class="text-sm font-medium text-gray-500"
												>#{question.order || index + 1}</span
											>
											<Badge variant="secondary">
												{getQuestionTypeLabel(question.type)}
											</Badge>
											{#if question.required}
												<Badge variant="primary">Obrigatória</Badge>
											{/if}
										</div>
										<p class="font-medium text-gray-900">{question.text}</p>

										<!-- Show options for multiple choice -->
										{#if question.type === 'multiple_choice' && question.options}
											<div class="mt-2 ml-4">
												{#each JSON.parse(question.options) as option}
													<div class="text-sm text-gray-600">• {option}</div>
												{/each}
											</div>
										{/if}
									</div>

									<!-- Question Actions -->
									<div class="flex items-center space-x-2">
										<Button size="sm" variant="outline" onclick={() => startEditQuestion(question)}>Editar</Button>
										<Button size="sm" variant="outline" onclick={() => deleteQuestion(question)}>
											<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
												></path>
											</svg>
										</Button>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</Card>
			{/if}

			<!-- Add Question Section -->
			<Card>
				{#if !showAddForm}
					<div class="py-8 text-center">
						<div class="mb-4">
							<svg
								class="mx-auto h-12 w-12 text-gray-400"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
								></path>
							</svg>
						</div>
						<h3 class="mb-2 text-lg font-medium text-gray-900">
							{questions.length === 0 ? 'Adicione a primeira questão' : 'Adicionar nova questão'}
						</h3>
						<p class="mb-4 text-gray-600">
							{questions.length === 0
								? 'Sua pesquisa precisa de pelo menos uma questão para ser utilizada pelos alunos.'
								: 'Adicione mais questões para enriquecer sua pesquisa.'}
						</p>
						<Button onclick={() => (showAddForm = true)}>
							<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 4v16m8-8H4"
								></path>
							</svg>
							Adicionar Questão
						</Button>
					</div>
				{:else}
					<!-- Add/Edit Question Form -->
					<div class="space-y-6">
						<h2 class="text-lg font-semibold text-gray-900">
							{editingQuestion ? 'Editar Questão' : 'Nova Questão'}
						</h2>

						<!-- Question Type Selection -->
						<div>
							<label class="mb-2 block text-sm font-medium text-gray-700">Tipo de Questão</label>
							<div class="grid grid-cols-1 gap-3 md:grid-cols-2">
								{#each questionTypes as type}
									<label class="relative">
										<input
											type="radio"
											name="questionType"
											value={type.value}
											bind:group={questionForm.type}
											class="peer sr-only"
										/>
										<div
											class="cursor-pointer rounded-lg border border-gray-300 p-3 peer-checked:border-blue-500 peer-checked:bg-blue-50 hover:bg-gray-50"
										>
											<div class="font-medium text-gray-900">{type.label}</div>
											<div class="text-sm text-gray-600">{type.description}</div>
										</div>
									</label>
								{/each}
							</div>
						</div>

						{#if questionForm.type}
							<!-- Question Text -->
							<div>
								<label for="questionText" class="mb-1 block text-sm font-medium text-gray-700">
									Texto da Questão *
								</label>
								<textarea
									id="questionText"
									bind:value={questionForm.text}
									placeholder="Digite o texto da questão..."
									rows="3"
									class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
								></textarea>
							</div>

							<!-- Multiple Choice Options -->
							{#if questionForm.type === 'multiple_choice'}
								<div>
									<label class="mb-2 block text-sm font-medium text-gray-700"
										>Opções de Resposta</label
									>
									<div class="space-y-2">
										{#each questionForm.options as option, index}
											<div class="flex items-center space-x-2">
												<input
													type="text"
													value={option}
													onchange={(e) =>
														updateOption(index, (e.target as HTMLInputElement).value)}
													placeholder={`Opção ${index + 1}`}
													class="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
												/>
												{#if questionForm.options.length > 1}
													<Button size="sm" variant="outline" onclick={() => removeOption(index)}>
														<svg
															class="h-4 w-4"
															fill="none"
															stroke="currentColor"
															viewBox="0 0 24 24"
														>
															<path
																stroke-linecap="round"
																stroke-linejoin="round"
																stroke-width="2"
																d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
															></path>
														</svg>
													</Button>
												{/if}
											</div>
										{/each}
									</div>
									<Button size="sm" variant="outline" onclick={addOption}>
										<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M12 4v16m8-8H4"
											></path>
										</svg>
										Adicionar Opção
									</Button>
								</div>
							{/if}

							<!-- Required Toggle -->
							<div>
								<label class="flex items-center space-x-3">
									<input
										type="checkbox"
										bind:checked={questionForm.required}
										class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
									/>
									<span class="text-sm font-medium text-gray-700">Questão obrigatória</span>
								</label>
							</div>

							<!-- Error Message -->
							{#if error}
								<div class="rounded-md border border-red-200 bg-red-50 p-3">
									<p class="text-sm text-red-800">{error}</p>
								</div>
							{/if}

							<!-- Form Actions -->
							<div class="flex items-center justify-end space-x-3 border-t pt-4">
								<Button variant="outline" onclick={cancelAddQuestion} disabled={submitting}>
									Cancelar
								</Button>
								<Button onclick={editingQuestion ? updateQuestion : saveQuestion} disabled={submitting}>
									{submitting ? 'Salvando...' : (editingQuestion ? 'Atualizar Questão' : 'Salvar Questão')}
								</Button>
							</div>
						{/if}
					</div>
				{/if}
			</Card>
		{/if}
	</div>
</div>
