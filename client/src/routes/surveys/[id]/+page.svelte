<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Layout from '$lib/components/Layout.svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	let survey: any = null;
	let loading = true;
	let responses: { [questionId: number]: string } = {};
	let submitting = false;
	let submitted = false;
	let error = '';
	let studentResponses: any[] = [];
	let hasAnswered = false;

	// Check if survey is active and within date range
	function isSurveyActive() {
		if (!survey?.is_active) return false;
		
		const now = new Date();
		const openDate = new Date(survey.open_date);
		const closeDate = new Date(survey.close_date);
		
		return now >= openDate && now <= closeDate;
	}

	// Format date for display
	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('pt-BR', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	// Handle form submission
	async function submitSurvey(event: Event) {
		event.preventDefault();
		
		// Validate required questions
		const requiredQuestions = survey.questions.filter((q: any) => q.required);
		const missingRequired = requiredQuestions.filter((q: any) => !responses[q.id]);
		
		if (missingRequired.length > 0) {
			error = `Por favor, responda todas as questões obrigatórias (${missingRequired.length} faltando)`;
			return;
		}

		submitting = true;
		error = '';

		try {
			// Submit each response individually as per backend API design
			for (const [questionId, answer] of Object.entries(responses)) {
				if (answer.trim() === '') continue; // Skip empty answers

				const response = {
					survey_id: survey.id,
					question_id: parseInt(questionId),
					answer: answer.trim()
				};

				const result = await api.submitResponse(response);
				if (!result.success) {
					throw new Error(result.error || 'Erro ao enviar resposta');
				}
			}

			submitted = true;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Erro ao enviar respostas';
		} finally {
			submitting = false;
		}
	}

	// Parse multiple choice options
	function parseOptions(optionsString: string): string[] {
		try {
			return JSON.parse(optionsString);
		} catch {
			return [];
		}
	}

	// Handle NPS rating
	function selectNPS(questionId: number, value: number) {
		responses[questionId] = value.toString();
	}

	// Handle star rating
	function selectRating(questionId: number, value: number) {
		responses[questionId] = value.toString();
	}

	onMount(async () => {
		const surveyId = $page.params.id;
		
		if (!surveyId) {
			goto('/dashboard/student');
			return;
		}

		try {
			// Load survey details
			const surveyResult = await api.getSurveyById(surveyId);
			
			if (!surveyResult.success) {
				error = surveyResult.error || 'Pesquisa não encontrada';
				loading = false;
				return;
			}

			survey = (surveyResult.data as any)?.survey;
			
			if (!survey) {
				error = 'Pesquisa não encontrada';
				loading = false;
				return;
			}

			// Check if student has already answered this survey
			const responsesResult = await api.getSurveyResponses(surveyId);
			
			if (responsesResult.success) {
				studentResponses = (responsesResult.data as any)?.responses || [];
				hasAnswered = studentResponses.length > 0;

				// If student has answered, populate the responses object for display
				if (hasAnswered) {
					studentResponses.forEach((response: any) => {
						responses[response.question_id] = response.answer;
					});
				}
			}
		} catch (err) {
			error = 'Erro ao carregar pesquisa';
			console.error('Failed to load survey:', err);
		} finally {
			loading = false;
		}
	});
</script>

<Layout title="Responder Pesquisa">
	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="text-center">
				<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
				<p class="mt-4 text-gray-600">Carregando pesquisa...</p>
			</div>
		</div>
	{:else if error}
		<Card class="border-red-200 bg-red-50">
			<div class="text-center">
				<p class="text-red-800">{error}</p>
				<Button variant="outline" onclick={() => goto('/dashboard/student')} class="mt-4">
					Voltar ao Dashboard
				</Button>
			</div>
		</Card>
	{:else if !survey}
		<Card class="border-red-200 bg-red-50">
			<div class="text-center">
				<p class="text-red-800">Pesquisa não encontrada</p>
				<Button variant="outline" onclick={() => goto('/dashboard/student')} class="mt-4">
					Voltar ao Dashboard
				</Button>
			</div>
		</Card>
	{:else if submitted || hasAnswered}
		<Card class="border-green-200 bg-green-50">
			<div class="text-center space-y-4">
				<div class="mx-auto w-16 h-16 bg-green-100 rounded-full flex items-center justify-center">
					<svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
					</svg>
				</div>
				<h2 class="text-2xl font-bold text-green-800">
					{submitted ? 'Respostas Enviadas!' : 'Pesquisa Já Respondida'}
				</h2>
				<p class="text-green-700">
					{submitted 
						? `Obrigado por participar da pesquisa sobre ${survey.subject.name}.`
						: `Você já respondeu esta pesquisa sobre ${survey.subject.name}.`
					}
				</p>
				<Button onclick={() => goto('/dashboard/student')}>
					Voltar ao Dashboard
				</Button>
			</div>
		</Card>

		<!-- Show answered questions if already responded -->
		{#if hasAnswered && !submitted}
			<div class="space-y-6">
				<Card>
					<h3 class="text-lg font-semibold text-gray-900 mb-4">Suas Respostas:</h3>
				</Card>

				{#each survey.questions as question, index}
					<Card>
						<div class="space-y-4">
							<!-- Question Header -->
							<div class="flex items-start justify-between">
								<div class="flex-1">
									<h3 class="text-lg font-medium text-gray-900">
										{index + 1}. {question.text}
									</h3>
								</div>
								<Badge variant="secondary" class="text-xs">
									{question.type === 'nps' ? 'NPS' : 
									 question.type === 'free_text' ? 'Texto Livre' :
									 question.type === 'rating' ? 'Avaliação' : 'Múltipla Escolha'}
								</Badge>
							</div>

							<!-- Display Answer -->
							<div class="bg-gray-50 rounded-lg p-4">
								{#if question.type === 'free_text'}
									<p class="text-gray-800">{responses[question.id] || 'Não respondido'}</p>

								{:else if question.type === 'nps'}
									<div class="flex items-center space-x-3">
										<span class="font-medium text-gray-700">Nota:</span>
										<span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
											{responses[question.id] || 'N/A'}/10
										</span>
									</div>

								{:else if question.type === 'rating'}
									<div class="flex items-center space-x-3">
										<span class="font-medium text-gray-700">Avaliação:</span>
										<div class="flex gap-1">
											{#each Array(5) as _, i}
												<svg class="w-5 h-5 {parseInt(responses[question.id]) > i ? 'text-yellow-400' : 'text-gray-300'}" 
													 fill="currentColor" viewBox="0 0 20 20">
													<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
												</svg>
											{/each}
										</div>
										<span class="text-sm text-gray-600">({responses[question.id] || 0}/5)</span>
									</div>

								{:else if question.type === 'multiple_choice'}
									<div class="flex items-center space-x-3">
										<span class="font-medium text-gray-700">Resposta:</span>
										<span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
											{responses[question.id] || 'Não respondido'}
										</span>
									</div>
								{/if}
							</div>
						</div>
					</Card>
				{/each}
			</div>
		{/if}
	{:else if !isSurveyActive()}
		<Card class="border-yellow-200 bg-yellow-50">
			<div class="text-center">
				<h2 class="text-xl font-semibold text-yellow-800 mb-2">Pesquisa Indisponível</h2>
				<p class="text-yellow-700">Esta pesquisa não está ativa ou está fora do período de participação.</p>
				<div class="mt-4 text-sm text-yellow-600">
					<p>Período: {formatDate(survey.open_date)} a {formatDate(survey.close_date)}</p>
				</div>
				<Button variant="outline" onclick={() => goto('/dashboard/student')} class="mt-4">
					Voltar ao Dashboard
				</Button>
			</div>
		</Card>
	{:else if !hasAnswered}
		<div class="space-y-6">
			<!-- Survey Header -->
			<Card>
				<div class="space-y-4">
					<div class="flex items-start justify-between">
						<div>
							<h1 class="text-2xl font-bold text-gray-900">{survey.title}</h1>
							<p class="text-gray-600">{survey.subject.name} ({survey.subject.code})</p>
						</div>
						<Badge variant="success">Ativa</Badge>
					</div>
					
					{#if survey.description}
						<p class="text-gray-700">{survey.description}</p>
					{/if}
					
					<div class="flex flex-wrap gap-4 text-sm text-gray-600">
						<div>
							<span class="font-medium">Professor:</span> 
							{survey.subject.professor?.first_name} {survey.subject.professor?.last_name}
						</div>
						<div>
							<span class="font-medium">Semestre:</span> {survey.semester.name}
						</div>
						<div>
							<span class="font-medium">Questões:</span> {survey.questions.length}
						</div>
					</div>
				</div>
			</Card>

			<!-- Error Message -->
			{#if error}
				<Card class="border-red-200 bg-red-50">
					<p class="text-red-800">{error}</p>
				</Card>
			{/if}

			<!-- Questions Form -->
			<form onsubmit={submitSurvey}>
				<div class="space-y-6">
					{#each survey.questions as question, index}
						<Card>
							<div class="space-y-4">
								<!-- Question Header -->
								<div class="flex items-start justify-between">
									<div class="flex-1">
										<h3 class="text-lg font-medium text-gray-900">
											{index + 1}. {question.text}
											{#if question.required}
												<span class="text-red-500">*</span>
											{/if}
										</h3>
									</div>
									<Badge variant="secondary" class="text-xs">
										{question.type === 'nps' ? 'NPS' : 
										 question.type === 'free_text' ? 'Texto Livre' :
										 question.type === 'rating' ? 'Avaliação' : 'Múltipla Escolha'}
									</Badge>
								</div>

								<!-- Question Input Based on Type -->
								{#if question.type === 'free_text'}
									<textarea
										bind:value={responses[question.id]}
										placeholder="Digite sua resposta..."
										rows="4"
										class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
									></textarea>

								{:else if question.type === 'nps'}
									<div class="space-y-3">
										<p class="text-sm text-gray-600">Em uma escala de 0 a 10, o quanto você recomendaria?</p>
										<div class="flex flex-wrap gap-2">
											{#each Array(11) as _, i}
												<button
													type="button"
													onclick={() => selectNPS(question.id, i)}
													class="w-12 h-12 rounded-lg border-2 font-medium transition-colors
														{responses[question.id] === i.toString() 
															? 'border-blue-500 bg-blue-500 text-white' 
															: 'border-gray-300 bg-white text-gray-700 hover:border-blue-300'}"
												>
													{i}
												</button>
											{/each}
										</div>
										<div class="flex justify-between text-xs text-gray-500">
											<span>Muito improvável</span>
											<span>Muito provável</span>
										</div>
									</div>

								{:else if question.type === 'rating'}
									<div class="space-y-3">
										<p class="text-sm text-gray-600">Avalie de 1 a 5 estrelas:</p>
										<div class="flex gap-1">
											{#each Array(5) as _, i}
												<button
													type="button"
													onclick={() => selectRating(question.id, i + 1)}
													class="w-8 h-8 transition-colors"
												>
													<svg class="w-full h-full {parseInt(responses[question.id]) > i ? 'text-yellow-400' : 'text-gray-300'}" 
														 fill="currentColor" viewBox="0 0 20 20">
														<path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
													</svg>
												</button>
											{/each}
										</div>
									</div>

								{:else if question.type === 'multiple_choice'}
									<div class="space-y-2">
										{#each parseOptions(question.options) as option}
											<label class="flex items-center space-x-3 cursor-pointer">
												<input
													type="radio"
													name="question_{question.id}"
													value={option}
													bind:group={responses[question.id]}
													class="h-4 w-4 text-blue-600 border-gray-300 focus:ring-blue-500"
												/>
												<span class="text-gray-700">{option}</span>
											</label>
										{/each}
									</div>
								{/if}
							</div>
						</Card>
					{/each}

					<!-- Submit Button -->
					<Card>
						<div class="flex items-center justify-between">
							<div class="text-sm text-gray-600">
								<p>Questões obrigatórias marcadas com <span class="text-red-500">*</span></p>
							</div>
							<div class="flex gap-3">
								<Button 
									variant="outline" 
									onclick={() => goto('/dashboard/student')}
									disabled={submitting}
								>
									Cancelar
								</Button>
								<Button 
									type="submit" 
									disabled={submitting}
									class="px-8"
								>
									{submitting ? 'Enviando...' : 'Enviar Respostas'}
								</Button>
							</div>
						</div>
					</Card>
				</div>
			</form>
		</div>
	{/if}
</Layout> 