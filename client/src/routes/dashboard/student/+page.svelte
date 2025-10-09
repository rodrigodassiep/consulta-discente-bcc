<script lang="ts">
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	let surveys: any[] = [];
	let enrollments: any[] = [];
	let pastResponses: any[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		loading = true;
		error = '';

		try {
			// Load surveys, enrollments, and past responses in parallel
			const [surveysResult, enrollmentsResult, responsesResult] = await Promise.all([
				api.getStudentSurveys(),
				api.getStudentSubjects(),
				api.getStudentResponses()
			]);

			if (!surveysResult.success) {
				throw new Error(surveysResult.error || 'Erro ao carregar pesquisas');
			}
			if (!enrollmentsResult.success) {
				throw new Error(enrollmentsResult.error || 'Erro ao carregar disciplinas');
			}
			if (!responsesResult.success) {
				throw new Error(responsesResult.error || 'Erro ao carregar respostas');
			}

			surveys = surveysResult.data?.surveys || [];
			enrollments = enrollmentsResult.data?.enrollments || [];
			pastResponses = responsesResult.data?.responses || [];

		} catch (err) {
			error = err instanceof Error ? err.message : 'Erro desconhecido';
		} finally {
			loading = false;
		}
	}

	function formatDate(dateString: string) {
		return new Date(dateString).toLocaleDateString('pt-BR', {
			day: '2-digit',
			month: '2-digit',
			year: 'numeric'
		});
	}

	function isSurveyActive(survey: any) {
		if (!survey.is_active) return false;
		
		const now = new Date();
		const openDate = new Date(survey.open_date);
		const closeDate = new Date(survey.close_date);
		
		return now >= openDate && now <= closeDate;
	}

	function getSurveyStatus(survey: any) {
		if (!survey.is_active) return { text: 'Inativa', variant: 'secondary' as const };
		
		const now = new Date();
		const openDate = new Date(survey.open_date);
		const closeDate = new Date(survey.close_date);
		
		if (now < openDate) return { text: 'Em breve', variant: 'warning' as const };
		if (now > closeDate) return { text: 'Encerrada', variant: 'danger' as const };
		return { text: 'Ativa', variant: 'success' as const };
	}

	function hasAnsweredSurvey(surveyId: number) {
		return pastResponses.some(response => response.survey_id === surveyId);
	}

	function goToSurvey(surveyId: number) {
		window.location.href = `/surveys/${surveyId}`;
	}
</script>

<svelte:head>
	<title>Dashboard do Estudante - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-8">
	<div>
		<h2 class="text-3xl font-bold text-gray-900">Dashboard do Estudante</h2>
	</div>

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="text-center">
				<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
				<p class="mt-4 text-gray-600">Carregando...</p>
			</div>
		</div>
	{:else if error}
		<Card class="border-red-200 bg-red-50">
			<div class="text-center">
				<p class="text-red-800">{error}</p>
				<Button variant="outline" onclick={loadData} class="mt-4">
					Tentar novamente
				</Button>
			</div>
		</Card>
	{:else}
		<div class="space-y-8">
			<!-- Enrolled Subjects Overview -->
			<section>
				<h3 class="mb-4 text-xl font-semibold text-gray-900">Disciplinas Matriculadas</h3>
				{#if enrollments.length === 0}
					<Card>
						<p class="text-center text-gray-500">Você não está matriculado em nenhuma disciplina.</p>
					</Card>
				{:else}
					<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
						{#each enrollments as enrollment}
							<Card class="h-fit">
								<div class="space-y-2">
									<h4 class="font-semibold text-gray-900">{enrollment.subject.name}</h4>
									<p class="text-sm text-gray-600">Código: {enrollment.subject.code}</p>
									<p class="text-sm text-gray-600">
										Professor: {enrollment.subject.professor.first_name} {enrollment.subject.professor.last_name}
									</p>
									<p class="text-sm text-gray-600">
										Semestre: {enrollment.semester.name}
									</p>
								</div>
							</Card>
						{/each}
					</div>
				{/if}
			</section>

			<!-- Available Surveys -->
			<section>
				<h3 class="mb-4 text-xl font-semibold text-gray-900">Pesquisas Disponíveis</h3>
				{#if surveys.length === 0}
					<Card>
						<p class="text-center text-gray-500">Não há pesquisas disponíveis no momento.</p>
					</Card>
				{:else}
					<div class="space-y-4">
						{#each surveys as survey}
							{@const status = getSurveyStatus(survey)}
							{@const isActive = isSurveyActive(survey)}
							{@const hasAnswered = hasAnsweredSurvey(survey.id)}
							
							<Card class="border-l-4 {isActive ? 'border-l-green-500' : 'border-l-gray-300'}">
								<div class="flex items-start justify-between">
									<div class="flex-1 space-y-3">
										<div class="flex items-start justify-between">
											<div>
												<h4 class="text-lg font-semibold text-gray-900">{survey.title}</h4>
												<p class="text-sm text-gray-600">{survey.subject.name} ({survey.subject.code})</p>
											</div>
											<Badge variant={status.variant}>{status.text}</Badge>
										</div>
										
										{#if survey.description}
											<p class="text-gray-700">{survey.description}</p>
										{/if}
										
										<div class="flex flex-wrap gap-4 text-sm text-gray-600">
											<div>
												<span class="font-medium">Abertura:</span> {formatDate(survey.open_date)}
											</div>
											<div>
												<span class="font-medium">Fechamento:</span> {formatDate(survey.close_date)}
											</div>
											<div>
												<span class="font-medium">Questões:</span> {survey.questions?.length || 0}
											</div>
										</div>

										{#if hasAnswered}
											<div class="flex items-center space-x-2">
												<Badge variant="success">Respondida</Badge>
												<span class="text-sm text-gray-600">Você já respondeu esta pesquisa</span>
											</div>
										{/if}
									</div>
									
									<div class="ml-4 flex-shrink-0">
										{#if isActive && !hasAnswered}
											<Button onclick={() => goToSurvey(survey.id)}>
												Responder
											</Button>
										{:else if hasAnswered}
											<Button variant="outline" onclick={() => goToSurvey(survey.id)}>
												Ver Respostas
											</Button>
										{:else}
											<Button variant="ghost" disabled>
												Indisponível
											</Button>
										{/if}
									</div>
								</div>
							</Card>
						{/each}
					</div>
				{/if}
			</section>

			<!-- Recent Activity -->
			<section>
				<h3 class="mb-4 text-xl font-semibold text-gray-900">Atividade Recente</h3>
				{#if pastResponses.length === 0}
					<Card>
						<p class="text-center text-gray-500">Você ainda não respondeu nenhuma pesquisa.</p>
					</Card>
				{:else}
					<Card>
						<div class="space-y-3">
							{#each pastResponses.slice(0, 5) as response}
								<div class="flex items-center justify-between border-b border-gray-100 pb-3 last:border-b-0 last:pb-0">
									<div>
										<p class="font-medium text-gray-900">{response.survey.title}</p>
										<p class="text-sm text-gray-600">{response.survey.subject.name}</p>
									</div>
									<div class="text-sm text-gray-500">
										{formatDate(response.submitted_at)}
									</div>
								</div>
							{/each}
							{#if pastResponses.length > 5}
								<div class="text-center">
									<Button variant="ghost" size="sm">
										Ver todas as respostas
									</Button>
								</div>
							{/if}
						</div>
					</Card>
				{/if}
			</section>
		</div>
	{/if}
</div> 