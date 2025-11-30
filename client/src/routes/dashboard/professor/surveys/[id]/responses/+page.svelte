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
	let responses: any[] = [];
	let loading = true;
	let error = '';
	let statistics: any = {};
	let selectedQuestion: any = null;
	let showResponsesModal = false;

	onMount(async () => {
		const surveyId = $page.params.id;

		if (!surveyId) {
			goto('/dashboard/professor');
			return;
		}

		await loadSurveyResponses(surveyId);
	});

	async function loadSurveyResponses(surveyId: string) {
		try {
			loading = true;

			// Load survey data and responses in parallel
			const [surveysResult, responsesResult] = await Promise.all([
				api.getProfessorSurveys(),
				api.getProfessorSurveyResponses(surveyId)
			]);

			if (!surveysResult.success) {
				throw new Error(surveysResult.error || 'Failed to load survey');
			}

			if (!responsesResult.success) {
				throw new Error(responsesResult.error || 'Failed to load responses');
			}

			const surveys = (surveysResult.data as any)?.surveys || [];
			survey = surveys.find((s: any) => s.id === parseInt(surveyId));

			if (!survey) {
				throw new Error('Survey not found or access denied');
			}

			responses = (responsesResult.data as any)?.responses || [];

			// Calculate statistics
			calculateStatistics();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Error loading responses';
			console.error('Failed to load responses:', err);
		} finally {
			loading = false;
		}
	}

	function calculateStatistics() {
		if (!responses || responses.length === 0) {
			statistics = {
				totalResponses: 0,
				questionStats: {}
			};
			return;
		}

		const questionStats: { [questionId: number]: any } = {};

		survey.questions?.forEach((question: any) => {
			const questionResponses = responses.filter((r) => r.question_id === question.id);
			const responseCount = questionResponses.length;

			if (question.type === 'rating') {
				const ratings = questionResponses.map((r) => parseInt(r.answer)).filter((r) => !isNaN(r));
				const average = ratings.length > 0
					? ratings.reduce((sum, rating) => sum + rating, 0) / ratings.length
					: 0;
				const distribution = [1, 2, 3, 4, 5].map(
					(star) => ratings.filter((r) => r === star).length
				);
				questionStats[question.id] = {
					type: 'rating',
					count: responseCount,
					average,
					distribution
				};
			} else if (question.type === 'nps') {
				const scores = questionResponses.map((r) => parseInt(r.answer)).filter((r) => !isNaN(r));
				const average = scores.length > 0
					? scores.reduce((sum, score) => sum + score, 0) / scores.length
					: 0;
				// NPS calculation: % Promoters (9-10) - % Detractors (0-6)
				const promoters = scores.filter((s) => s >= 9).length;
				const detractors = scores.filter((s) => s <= 6).length;
				const npsScore = scores.length > 0
					? ((promoters - detractors) / scores.length) * 100
					: 0;
				questionStats[question.id] = {
					type: 'nps',
					count: responseCount,
					average,
					npsScore,
					promoters,
					passives: scores.filter((s) => s >= 7 && s <= 8).length,
					detractors
				};
			} else if (question.type === 'multiple_choice') {
				const choices: { [choice: string]: number } = {};
				questionResponses.forEach((r) => {
					choices[r.answer] = (choices[r.answer] || 0) + 1;
				});
				questionStats[question.id] = {
					type: 'multiple_choice',
					count: responseCount,
					choices
				};
			} else {
				// free_text
				questionStats[question.id] = {
					type: 'free_text',
					count: responseCount,
					answers: questionResponses.map((r) => r.answer)
				};
			}
		});

		statistics = {
			totalResponses: responses.length,
			questionStats
		};
	}

	function viewQuestionResponses(question: any) {
		selectedQuestion = question;
		showResponsesModal = true;
	}

	function closeModal() {
		showResponsesModal = false;
		selectedQuestion = null;
	}

	function getQuestionTypeBadge(type: string) {
		switch (type) {
			case 'nps':
				return 'NPS';
			case 'rating':
				return 'Avaliacao';
			case 'multiple_choice':
				return 'Multipla Escolha';
			default:
				return 'Texto Livre';
		}
	}

	function handleBackToSurvey() {
		goto(`/dashboard/professor/surveys/${survey.id}/questions`);
	}

	function handleBackToDashboard() {
		goto('/dashboard/professor');
	}
</script>

<svelte:head>
	<title>Respostas da Pesquisa - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Respostas da Pesquisa</h1>
			<p class="mt-1 text-gray-600">Visualize as respostas anonimas dos alunos</p>
		</div>
		<div class="flex gap-2">
			<Button variant="outline" onclick={handleBackToSurvey} size="sm">
				<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M15 19l-7-7 7-7"
					></path>
				</svg>
				Gerenciar Questoes
			</Button>
			<Button variant="outline" onclick={handleBackToDashboard} size="sm">Dashboard</Button>
		</div>
	</div>

	{#if loading}
		<Card>
			<div class="py-8 text-center">
				<div
					class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-b-2 border-blue-600"
				></div>
				<p class="text-gray-600">Carregando respostas...</p>
			</div>
		</Card>
	{:else if error && !survey}
		<Card class="border-red-200 bg-red-50">
			<div class="py-8 text-center">
				<h3 class="mb-2 text-lg font-medium text-red-800">Erro</h3>
				<p class="mb-4 text-red-600">{error}</p>
				<Button onclick={handleBackToDashboard} variant="outline">Voltar ao Dashboard</Button>
			</div>
		</Card>
	{:else if survey}
		<!-- Survey Info -->
		<Card class="border-blue-200 bg-blue-50">
			<div class="flex items-center justify-between">
				<div>
					<h2 class="text-lg font-semibold text-blue-900">{survey.title}</h2>
					<p class="text-sm text-blue-700">
						{survey.subject?.name} ({survey.subject?.code}) - {survey.questions?.length || 0} questoes
					</p>
				</div>
				<div class="flex items-center gap-2">
					<Badge variant="secondary">
						{survey.is_active ? 'Ativa' : 'Inativa'}
					</Badge>
					<Badge variant="outline">
						Respostas anonimas
					</Badge>
				</div>
			</div>
		</Card>

		<!-- Statistics Overview -->
		{#if responses.length > 0}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-3">
				<Card>
					<div class="text-center">
						<div class="text-2xl font-bold text-blue-600">{statistics.totalResponses}</div>
						<div class="text-sm text-gray-600">Total de Respostas</div>
					</div>
				</Card>
				<Card>
					<div class="text-center">
						<div class="text-2xl font-bold text-green-600">{survey.questions?.length || 0}</div>
						<div class="text-sm text-gray-600">Questoes na Pesquisa</div>
					</div>
				</Card>
				<Card>
					<div class="text-center">
						<div class="text-2xl font-bold text-purple-600">
							{Math.round(statistics.totalResponses / (survey.questions?.length || 1))}
						</div>
						<div class="text-sm text-gray-600">Media de Respostas por Questao</div>
					</div>
				</Card>
			</div>

			<!-- Question-by-Question Results -->
			<Card>
				<h3 class="mb-4 text-lg font-semibold text-gray-900">Resultados por Questao</h3>
				<div class="space-y-6">
					{#each survey.questions || [] as question, index}
						{@const stats = statistics.questionStats[question.id]}
						<div class="border-b border-gray-200 pb-6 last:border-0">
							<div class="mb-3 flex items-start justify-between">
								<div class="flex-1">
									<p class="font-medium text-gray-900">
										{index + 1}. {question.text}
									</p>
									<p class="mt-1 text-sm text-gray-500">
										{stats?.count || 0} respostas
									</p>
								</div>
								<Badge variant="secondary">
									{getQuestionTypeBadge(question.type)}
								</Badge>
							</div>

							{#if stats}
								{#if stats.type === 'rating'}
									<!-- Rating visualization -->
									<div class="mt-3 rounded-lg bg-gray-50 p-4">
										<div class="mb-2 flex items-center gap-2">
											<span class="text-2xl font-bold text-yellow-600">
												{stats.average.toFixed(1)}
											</span>
											<span class="text-gray-600">/ 5 estrelas</span>
										</div>
										<div class="space-y-1">
											{#each [5, 4, 3, 2, 1] as star}
												<div class="flex items-center gap-2">
													<span class="w-8 text-sm text-gray-600">{star}</span>
													<div class="h-4 flex-1 rounded bg-gray-200">
														<div
															class="h-4 rounded bg-yellow-400"
															style="width: {stats.count > 0 ? (stats.distribution[star - 1] / stats.count) * 100 : 0}%"
														></div>
													</div>
													<span class="w-8 text-sm text-gray-600">{stats.distribution[star - 1]}</span>
												</div>
											{/each}
										</div>
									</div>
								{:else if stats.type === 'nps'}
									<!-- NPS visualization -->
									<div class="mt-3 rounded-lg bg-gray-50 p-4">
										<div class="mb-3 flex items-center gap-4">
											<div>
												<span class="text-2xl font-bold text-blue-600">
													{stats.average.toFixed(1)}
												</span>
												<span class="text-gray-600">/ 10 (media)</span>
											</div>
											<div>
												<span class="text-2xl font-bold {stats.npsScore >= 0 ? 'text-green-600' : 'text-red-600'}">
													{stats.npsScore >= 0 ? '+' : ''}{stats.npsScore.toFixed(0)}
												</span>
												<span class="text-gray-600">NPS Score</span>
											</div>
										</div>
										<div class="flex gap-4 text-sm">
											<div class="flex items-center gap-1">
												<span class="h-3 w-3 rounded-full bg-green-500"></span>
												<span>Promotores: {stats.promoters}</span>
											</div>
											<div class="flex items-center gap-1">
												<span class="h-3 w-3 rounded-full bg-yellow-500"></span>
												<span>Neutros: {stats.passives}</span>
											</div>
											<div class="flex items-center gap-1">
												<span class="h-3 w-3 rounded-full bg-red-500"></span>
												<span>Detratores: {stats.detractors}</span>
											</div>
										</div>
									</div>
								{:else if stats.type === 'multiple_choice'}
									<!-- Multiple choice visualization -->
									<div class="mt-3 rounded-lg bg-gray-50 p-4">
										<div class="space-y-2">
											{#each Object.entries(stats.choices) as [choice, count]}
												<div class="flex items-center gap-2">
													<div class="h-4 flex-1 rounded bg-gray-200">
														<div
															class="h-4 rounded bg-blue-500"
															style="width: {stats.count > 0 ? ((count as number) / stats.count) * 100 : 0}%"
														></div>
													</div>
													<span class="w-20 text-sm text-gray-600">{choice}</span>
													<span class="w-12 text-right text-sm font-medium">{count}</span>
												</div>
											{/each}
										</div>
									</div>
								{:else if stats.type === 'free_text'}
									<!-- Free text - show button to view responses -->
									<div class="mt-3">
										<Button
											variant="outline"
											size="sm"
											onclick={() => viewQuestionResponses(question)}
										>
											Ver {stats.count} respostas de texto
										</Button>
									</div>
								{/if}
							{:else}
								<p class="mt-2 text-sm text-gray-500">Nenhuma resposta ainda</p>
							{/if}
						</div>
					{/each}
				</div>
			</Card>
		{:else}
			<Card>
				<div class="py-12 text-center">
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
							d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
						></path>
					</svg>
					<h3 class="mt-2 text-sm font-medium text-gray-900">Nenhuma resposta ainda</h3>
					<p class="mt-1 text-sm text-gray-500">
						Esta pesquisa ainda nao recebeu respostas dos alunos.
					</p>
				</div>
			</Card>
		{/if}
	{/if}

	<!-- Free Text Responses Modal -->
	{#if showResponsesModal && selectedQuestion}
		{@const stats = statistics.questionStats[selectedQuestion.id]}
		<div
			class="bg-opacity-50 fixed inset-0 z-50 h-full w-full overflow-y-auto bg-gray-600"
			onclick={closeModal}
		>
			<div
				class="relative top-20 mx-auto w-11/12 rounded-md border bg-white p-5 shadow-lg md:w-3/4 lg:w-1/2"
				onclick={(e) => e.stopPropagation()}
			>
				<div class="mt-3">
					<!-- Modal Header -->
					<div class="mb-4 flex items-center justify-between">
						<div>
							<h3 class="text-lg font-medium text-gray-900">
								Respostas Anonimas
							</h3>
							<p class="text-sm text-gray-600">{selectedQuestion.text}</p>
						</div>
						<button onclick={closeModal} class="text-gray-400 hover:text-gray-600">
							<svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M6 18L18 6M6 6l12 12"
								></path>
							</svg>
						</button>
					</div>

					<!-- Modal Content -->
					<div class="max-h-96 space-y-3 overflow-y-auto">
						{#if stats?.answers && stats.answers.length > 0}
							{#each stats.answers as answer, i}
								<div class="rounded-md bg-gray-50 p-3">
									<p class="text-sm text-gray-800">{answer}</p>
								</div>
							{/each}
						{:else}
							<p class="py-4 text-center text-gray-500">Nenhuma resposta de texto</p>
						{/if}
					</div>

					<!-- Modal Footer -->
					<div class="mt-6 flex justify-end">
						<Button onclick={closeModal} variant="outline">Fechar</Button>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
