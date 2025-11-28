<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	let subjects: any[] = [];
	let surveys: any[] = [];
	let loading = true;
	let error = '';

	onMount(async () => {
		try {
			// Load both subjects and surveys in parallel
			const [subjectsResult, surveysResult] = await Promise.all([
				api.getProfessorSubjects(),
				api.getProfessorSurveys()
			]);

			if (!subjectsResult.success) {
				error = subjectsResult.error || 'Erro ao carregar disciplinas';
				return;
			}

			if (!surveysResult.success) {
				error = surveysResult.error || 'Erro ao carregar pesquisas';
				return;
			}

			subjects = (subjectsResult.data as any)?.subjects || [];
			surveys = (surveysResult.data as any)?.surveys || [];
		} catch (err) {
			error = 'Erro ao conectar com o servidor';
			console.error('Failed to load professor data:', err);
		} finally {
			loading = false;
		}
	});

	function getSurveyStatus(survey: any): { status: string; color: string } {
		const now = new Date();
		const openDate = new Date(survey.open_date);
		const closeDate = new Date(survey.close_date);

		if (!survey.is_active) {
			return { status: 'Inativa', color: 'gray' };
		}

		if (survey.questions?.length === 0) {
			return { status: 'Sem questões', color: 'yellow' };
		}

		if (now < openDate) {
			return { status: 'Agendada', color: 'blue' };
		}

		if (now >= openDate && now <= closeDate) {
			return { status: 'Ativa', color: 'green' };
		}

		return { status: 'Encerrada', color: 'red' };
	}

	// We'll show current semester info separately since subjects don't include semester data
</script>

<svelte:head>
	<title>Dashboard do Professor - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-6">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Minhas Disciplinas</h1>
				<p class="mt-1 text-gray-600">Gerencie suas disciplinas e pesquisas</p>
			</div>
		</div>

		<!-- Loading State -->
		{#if loading}
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
				{#each Array(3) as _}
					<Card>
						<div class="animate-pulse">
							<div class="mb-3 h-4 w-3/4 rounded bg-gray-200"></div>
							<div class="mb-4 h-3 w-1/2 rounded bg-gray-200"></div>
							<div class="h-8 w-full rounded bg-gray-200"></div>
						</div>
					</Card>
				{/each}
			</div>

			<!-- Error State -->
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-8 text-center">
					<div class="mb-2 text-red-600">
						<svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.962-.833-2.732 0L4.082 15.5c-.77.833.192 2.5 1.732 2.5z"
							></path>
						</svg>
					</div>
					<h3 class="mb-2 text-lg font-medium text-red-800">Erro ao carregar disciplinas</h3>
					<p class="mb-4 text-red-600">{error}</p>
					<Button onclick={() => window.location.reload()} variant="outline">
						Tentar novamente
					</Button>
				</div>
			</Card>

			<!-- Empty State -->
		{:else if subjects.length === 0}
			<Card class="border-gray-200">
				<div class="py-12 text-center">
					<div class="mb-4 text-gray-400">
						<svg class="mx-auto h-16 w-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
							></path>
						</svg>
					</div>
					<h3 class="mb-2 text-lg font-medium text-gray-900">Nenhuma disciplina encontrada</h3>
					<p class="mb-4 text-gray-600">
						Você não possui disciplinas atribuídas no sistema. Entre em contato com o administrador
						para verificar suas atribuições.
					</p>
				</div>
			</Card>

			<!-- Content Grid -->
		{:else}
			<!-- Subjects Section -->
			<div>
				<div class="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
					{#each subjects as subject}
						<Card class="transition-shadow hover:shadow-md">
							<div class="space-y-4">
								<!-- Subject Header -->
								<div>
									<div class="mb-2 flex items-start justify-between">
										<h3 class="text-lg font-semibold text-gray-900">{subject.name}</h3>
										<Badge variant="secondary">
											{subject.code}
										</Badge>
									</div>

									{#if subject.description}
										<p class="line-clamp-2 text-sm text-gray-600">{subject.description}</p>
									{/if}
								</div>

								<!-- Subject Info -->
								<div class="space-y-2 text-sm text-gray-600">
									<div class="flex items-center">
										<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
											></path>
										</svg>
										<span>Semestre ativo</span>
									</div>
								</div>

								<!-- Actions -->
								<div class="flex gap-2 pt-2">
									<Button
										size="sm"
										onclick={() =>
											goto(`/dashboard/professor/surveys/create?subject=${subject.id}`)}
									>
										<svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M12 4v16m8-8H4"
											></path>
										</svg>
										Criar Pesquisa
									</Button>
									<Button
										variant="outline"
										size="sm"
										onclick={() => goto(`/dashboard/professor/subjects/${subject.id}`)}
									>
										<svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path
												stroke-linecap="round"
												stroke-linejoin="round"
												stroke-width="2"
												d="M9 5l7 7-7 7"
											></path>
										</svg>
									</Button>
								</div>
							</div>
						</Card>
					{/each}
				</div>
			</div>

			<!-- Surveys Section -->
			{#if surveys.length > 0}
				<div class="mt-8">
					<h2 class="mb-4 text-xl font-semibold text-gray-900">Minhas Pesquisas</h2>
					<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
						{#each surveys as survey}
							{@const statusInfo = getSurveyStatus(survey)}
							<Card class="transition-shadow hover:shadow-md">
								<div class="space-y-3">
									<!-- Survey Header -->
									<div class="flex items-start justify-between">
										<div class="flex-1">
											<h3 class="text-lg font-semibold text-gray-900">{survey.title}</h3>
											<p class="text-sm text-gray-600">
												{survey.subject?.name} ({survey.subject?.code})
											</p>
										</div>
										<Badge variant={statusInfo.color === 'green' ? 'primary' : 'secondary'}>
											{statusInfo.status}
										</Badge>
									</div>

									<!-- Survey Description -->
									{#if survey.description}
										<p class="line-clamp-2 text-sm text-gray-600">{survey.description}</p>
									{/if}

									<!-- Survey Stats -->
									<div class="flex items-center space-x-4 text-sm text-gray-500">
										<div class="flex items-center">
											<svg
												class="mr-1 h-4 w-4"
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
											{survey.questions?.length || 0} questões
										</div>
										<div class="flex items-center">
											<svg
												class="mr-1 h-4 w-4"
												fill="none"
												stroke="currentColor"
												viewBox="0 0 24 24"
											>
												<path
													stroke-linecap="round"
													stroke-linejoin="round"
													stroke-width="2"
													d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
												></path>
											</svg>
											{new Date(survey.open_date).toLocaleDateString('pt-BR')}
										</div>
									</div>

									<!-- Actions -->
									<div class="flex gap-2 border-t pt-2">
										{#if survey.questions?.length === 0}
											<Button
												size="sm"
												onclick={() => goto(`/dashboard/professor/surveys/${survey.id}/questions`)}
											>
												<svg
													class="mr-2 h-4 w-4"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M12 4v16m8-8H4"
													></path>
												</svg>
												Adicionar Questões
											</Button>
										{:else}
											<Button
												size="sm"
												variant="outline"
												onclick={() => goto(`/dashboard/professor/surveys/${survey.id}/questions`)}
											>
												<svg
													class="mr-2 h-4 w-4"
													fill="none"
													stroke="currentColor"
													viewBox="0 0 24 24"
												>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
													></path>
													<path
														stroke-linecap="round"
														stroke-linejoin="round"
														stroke-width="2"
														d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
													></path>
												</svg>
												Gerenciar
											</Button>
											{#if statusInfo.status === 'Ativa' || statusInfo.status === 'Encerrada'}
												<Button
													size="sm"
													variant="outline"
													onclick={() =>
														goto(`/dashboard/professor/surveys/${survey.id}/responses`)}
												>
													<svg
														class="mr-2 h-4 w-4"
														fill="none"
														stroke="currentColor"
														viewBox="0 0 24 24"
													>
														<path
															stroke-linecap="round"
															stroke-linejoin="round"
															stroke-width="2"
															d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
														></path>
													</svg>
													Respostas
												</Button>
											{/if}
										{/if}
									</div>
								</div>
							</Card>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Summary Stats -->
			<Card class="border-blue-200 bg-blue-50">
				<div class="flex items-center space-x-4">
					<div class="flex-shrink-0">
						<div class="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100">
							<svg
								class="h-6 w-6 text-blue-600"
								fill="none"
								stroke="currentColor"
								viewBox="0 0 24 24"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.746 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
								></path>
							</svg>
						</div>
					</div>
					<div>
						<h3 class="text-sm font-medium text-blue-800">Total de Disciplinas</h3>
						<p class="text-2xl font-bold text-blue-900">{subjects.length}</p>
					</div>
				</div>
			</Card>
		{/if}
</div>
