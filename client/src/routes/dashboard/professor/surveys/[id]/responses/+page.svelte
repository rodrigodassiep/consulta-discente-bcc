<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Layout from '$lib/components/Layout.svelte';
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
	let selectedStudent: any = null;
	let showStudentModal = false;

	// Filter states
	let searchTerm = '';
	let sortBy = 'submitted_at';
	let sortDirection = 'desc';

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
				uniqueStudents: 0,
				averageRatings: {},
				npsScores: {},
				completionRate: 0
			};
			return;
		}

		// Group responses by student
		const studentResponses = groupResponsesByStudent();
		const uniqueStudents = Object.keys(studentResponses).length;

		// Calculate averages for rating and NPS questions
		const averageRatings: { [questionId: number]: number } = {};
		const npsScores: { [questionId: number]: number } = {};

		survey.questions?.forEach((question: any) => {
			const questionResponses = responses.filter((r) => r.question_id === question.id);

			if (question.type === 'rating') {
				const ratings = questionResponses.map((r) => parseInt(r.answer)).filter((r) => !isNaN(r));
				if (ratings.length > 0) {
					averageRatings[question.id] =
						ratings.reduce((sum, rating) => sum + rating, 0) / ratings.length;
				}
			} else if (question.type === 'nps') {
				const scores = questionResponses.map((r) => parseInt(r.answer)).filter((r) => !isNaN(r));
				if (scores.length > 0) {
					npsScores[question.id] = scores.reduce((sum, score) => sum + score, 0) / scores.length;
				}
			}
		});

		// Calculate completion rate (students who answered at least one question)
		const expectedResponses = uniqueStudents * (survey.questions?.length || 0);
		const completionRate = expectedResponses > 0 ? (responses.length / expectedResponses) * 100 : 0;

		statistics = {
			totalResponses: responses.length,
			uniqueStudents,
			averageRatings,
			npsScores,
			completionRate
		};
	}

	function groupResponsesByStudent() {
		const grouped: { [studentId: number]: any[] } = {};

		responses.forEach((response) => {
			if (!grouped[response.student_id]) {
				grouped[response.student_id] = [];
			}
			grouped[response.student_id].push(response);
		});

		return grouped;
	}

	function getFilteredStudents() {
		const studentResponses = groupResponsesByStudent();
		let students = Object.entries(studentResponses).map(([studentId, studentResponses]) => ({
			id: parseInt(studentId),
			student: studentResponses[0]?.student,
			responses: studentResponses,
			responseCount: studentResponses.length,
			lastSubmitted: new Date(
				Math.max(...studentResponses.map((r) => new Date(r.submitted_at).getTime()))
			)
		}));

		// Apply search filter
		if (searchTerm) {
			const term = searchTerm.toLowerCase();
			students = students.filter(
				(student) =>
					student.student?.first_name?.toLowerCase().includes(term) ||
					student.student?.last_name?.toLowerCase().includes(term) ||
					student.student?.email?.toLowerCase().includes(term)
			);
		}

		// Apply sorting
		students.sort((a, b) => {
			let aValue: any, bValue: any;

			switch (sortBy) {
				case 'name':
					aValue = `${a.student?.first_name} ${a.student?.last_name}`.toLowerCase();
					bValue = `${b.student?.first_name} ${b.student?.last_name}`.toLowerCase();
					break;
				case 'response_count':
					aValue = a.responseCount;
					bValue = b.responseCount;
					break;
				case 'submitted_at':
				default:
					aValue = a.lastSubmitted.getTime();
					bValue = b.lastSubmitted.getTime();
					break;
			}

			if (sortDirection === 'asc') {
				return aValue > bValue ? 1 : -1;
			} else {
				return aValue < bValue ? 1 : -1;
			}
		});

		return students;
	}

	function viewStudentResponses(student: any) {
		selectedStudent = student;
		showStudentModal = true;
	}

	function closeStudentModal() {
		showStudentModal = false;
		selectedStudent = null;
	}

	function formatDate(date: string | Date) {
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString('pt-BR', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getAnswerDisplay(question: any, answer: string) {
		if (question.type === 'rating') {
			const rating = parseInt(answer);
			return isNaN(rating) ? answer : `${rating}/5 estrelas`;
		} else if (question.type === 'nps') {
			const score = parseInt(answer);
			return isNaN(score) ? answer : `${score}/10`;
		} else if (question.type === 'multiple_choice') {
			return answer;
		} else {
			return answer;
		}
	}

	function handleBackToSurvey() {
		goto(`/dashboard/professor/surveys/${survey.id}/questions`);
	}

	function handleBackToDashboard() {
		goto('/dashboard/professor');
	}
</script>

<Layout title="Respostas da Pesquisa">
	<div class="mx-auto max-w-7xl space-y-6 p-6">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-2xl font-bold text-gray-900">Respostas da Pesquisa</h1>
				<p class="text-gray-600">Visualize e analise as respostas dos alunos</p>
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
					Gerenciar Questões
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
							{survey.subject?.name} ({survey.subject?.code}) • {survey.questions?.length || 0} questões
						</p>
					</div>
					<Badge variant="secondary">
						{survey.is_active ? 'Ativa' : 'Inativa'}
					</Badge>
				</div>
			</Card>

			<!-- Statistics -->
			{#if responses.length > 0}
				<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
					<Card>
						<div class="text-center">
							<div class="text-2xl font-bold text-blue-600">{statistics.totalResponses}</div>
							<div class="text-sm text-gray-600">Total de Respostas</div>
						</div>
					</Card>
					<Card>
						<div class="text-center">
							<div class="text-2xl font-bold text-green-600">{statistics.uniqueStudents}</div>
							<div class="text-sm text-gray-600">Alunos Participantes</div>
						</div>
					</Card>
					<Card>
						<div class="text-center">
							<div class="text-2xl font-bold text-purple-600">
								{statistics.completionRate.toFixed(1)}%
							</div>
							<div class="text-sm text-gray-600">Taxa de Completude</div>
						</div>
					</Card>
					<Card>
						<div class="text-center">
							<div class="text-2xl font-bold text-orange-600">{survey.questions?.length || 0}</div>
							<div class="text-sm text-gray-600">Questões na Pesquisa</div>
						</div>
					</Card>
				</div>

				<!-- Question Statistics -->
				{#if Object.keys(statistics.averageRatings).length > 0 || Object.keys(statistics.npsScores).length > 0}
					<Card>
						<h3 class="mb-4 text-lg font-semibold text-gray-900">Estatísticas por Questão</h3>
						<div class="space-y-4">
							{#each survey.questions || [] as question}
								{#if question.type === 'rating' && statistics.averageRatings[question.id]}
									<div class="border-l-4 border-yellow-400 pl-4">
										<p class="font-medium text-gray-900">{question.text}</p>
										<p class="text-sm text-gray-600">
											Avaliação média: <span class="font-semibold text-yellow-600">
												{statistics.averageRatings[question.id].toFixed(1)}/5 estrelas
											</span>
										</p>
									</div>
								{:else if question.type === 'nps' && statistics.npsScores[question.id]}
									<div class="border-l-4 border-blue-400 pl-4">
										<p class="font-medium text-gray-900">{question.text}</p>
										<p class="text-sm text-gray-600">
											NPS médio: <span class="font-semibold text-blue-600">
												{statistics.npsScores[question.id].toFixed(1)}/10
											</span>
										</p>
									</div>
								{/if}
							{/each}
						</div>
					</Card>
				{/if}

				<!-- Response List -->
				<Card>
					<div class="mb-4 flex items-center justify-between">
						<h3 class="text-lg font-semibold text-gray-900">Respostas dos Alunos</h3>
						<div class="flex items-center gap-4">
							<!-- Search -->
							<input
								type="text"
								placeholder="Buscar aluno..."
								bind:value={searchTerm}
								class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							/>
							<!-- Sort -->
							<select
								bind:value={sortBy}
								class="rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-1 focus:ring-blue-500 focus:outline-none"
							>
								<option value="submitted_at">Data de Submissão</option>
								<option value="name">Nome</option>
								<option value="response_count">Quantidade de Respostas</option>
							</select>
							<button
								onclick={() => (sortDirection = sortDirection === 'asc' ? 'desc' : 'asc')}
								class="rounded-md border border-gray-300 p-2 hover:bg-gray-50"
							>
								<svg
									class="h-4 w-4 transform {sortDirection === 'desc' ? 'rotate-180' : ''}"
									fill="none"
									stroke="currentColor"
									viewBox="0 0 24 24"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M5 15l7-7 7 7"
									></path>
								</svg>
							</button>
						</div>
					</div>

					{#if getFilteredStudents().length === 0}
						<div class="py-8 text-center">
							{#if searchTerm}
								<p class="text-gray-600">Nenhum aluno encontrado para "{searchTerm}"</p>
							{:else}
								<p class="text-gray-600">Nenhuma resposta encontrada ainda.</p>
							{/if}
						</div>
					{:else}
						<div class="overflow-x-auto">
							<table class="min-w-full divide-y divide-gray-200">
								<thead class="bg-gray-50">
									<tr>
										<th
											class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
										>
											Aluno
										</th>
										<th
											class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
										>
											Respostas
										</th>
										<th
											class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
										>
											Última Submissão
										</th>
										<th
											class="px-6 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase"
										>
											Ações
										</th>
									</tr>
								</thead>
								<tbody class="divide-y divide-gray-200 bg-white">
									{#each getFilteredStudents() as student}
										<tr class="hover:bg-gray-50">
											<td class="px-6 py-4 whitespace-nowrap">
												<div>
													<div class="text-sm font-medium text-gray-900">
														{student.student?.first_name}
														{student.student?.last_name}
													</div>
													<div class="text-sm text-gray-500">{student.student?.email}</div>
												</div>
											</td>
											<td class="px-6 py-4 whitespace-nowrap">
												<div class="flex items-center">
													<span class="text-sm text-gray-900">{student.responseCount}</span>
													<span class="ml-1 text-xs text-gray-500"
														>/ {survey.questions?.length || 0}</span
													>
													{#if student.responseCount === (survey.questions?.length || 0)}
														<Badge variant="success">Completo</Badge>
													{:else}
														<Badge variant="secondary">Parcial</Badge>
													{/if}
												</div>
											</td>
											<td class="px-6 py-4 text-sm whitespace-nowrap text-gray-900">
												{formatDate(student.lastSubmitted)}
											</td>
											<td class="px-6 py-4 text-sm font-medium whitespace-nowrap">
												<Button
													size="sm"
													variant="outline"
													onclick={() => viewStudentResponses(student)}
												>
													Ver Respostas
												</Button>
											</td>
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					{/if}
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
							Esta pesquisa ainda não recebeu respostas dos alunos.
						</p>
					</div>
				</Card>
			{/if}
		{/if}
	</div>

	<!-- Student Response Modal -->
	{#if showStudentModal && selectedStudent}
		<div
			class="bg-opacity-50 fixed inset-0 z-50 h-full w-full overflow-y-auto bg-gray-600"
			onclick={closeStudentModal}
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
								Respostas de {selectedStudent.student?.first_name}
								{selectedStudent.student?.last_name}
							</h3>
							<p class="text-sm text-gray-600">{selectedStudent.student?.email}</p>
						</div>
						<button onclick={closeStudentModal} class="text-gray-400 hover:text-gray-600">
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
					<div class="max-h-96 space-y-4 overflow-y-auto">
						{#each survey.questions || [] as question}
							{@const studentResponse = selectedStudent.responses.find(
								(r: any) => r.question_id === question.id
							)}
							<div class="border-b border-gray-200 pb-4">
								<div class="mb-2 flex items-start justify-between">
									<p class="text-sm font-medium text-gray-900">{question.text}</p>
									<Badge variant="secondary">
										{question.type === 'nps'
											? 'NPS'
											: question.type === 'rating'
												? 'Avaliação'
												: question.type === 'multiple_choice'
													? 'Múltipla Escolha'
													: 'Texto Livre'}
									</Badge>
								</div>

								{#if studentResponse}
									<div class="rounded-md bg-gray-50 p-3">
										<p class="text-sm text-gray-800">
											{getAnswerDisplay(question, studentResponse.answer)}
										</p>
										<p class="mt-1 text-xs text-gray-500">
											Respondido em: {formatDate(studentResponse.submitted_at)}
										</p>
									</div>
								{:else}
									<div class="rounded-md bg-red-50 p-3">
										<p class="text-sm text-red-600">Não respondido</p>
									</div>
								{/if}
							</div>
						{/each}
					</div>

					<!-- Modal Footer -->
					<div class="mt-6 flex justify-end">
						<Button onclick={closeStudentModal} variant="outline">Fechar</Button>
					</div>
				</div>
			</div>
		</div>
	{/if}
</Layout>
