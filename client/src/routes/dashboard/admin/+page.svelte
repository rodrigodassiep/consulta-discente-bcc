<script lang="ts">
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	type Role = 'student' | 'professor' | 'admin';
	type AdminTab = 'roleRequests' | 'semesters' | 'subjects' | 'enrollments' | 'users';

	type User = {
		id: number;
		first_name: string;
		last_name: string;
		email: string;
		role: Role;
		requested_role?: Role;
		created_at?: string;
		updated_at?: string;
	};

	type Semester = {
		id: number;
		name: string;
		year: number;
		period: number;
		start_date: string;
		end_date: string;
		is_active: boolean;
	};

	type Subject = {
		id: number;
		name: string;
		code: string;
		description?: string;
		professor_id: number;
		professor?: User;
	};

	type Enrollment = {
		id: number;
		student_id: number;
		subject_id: number;
		semester_id: number;
		student?: User;
		subject?: Subject;
		semester?: Semester;
	};

	let currentTab: AdminTab = 'roleRequests';

	let loading = true;
	let error = '';
	let updatingId: number | null = null;

	let roleRequests: User[] = [];
	let semesters: Semester[] = [];
	let subjects: Subject[] = [];
	let enrollments: Enrollment[] = [];
	let allUsers: User[] = [];

	let professorUsers: User[] = [];
	let studentUsers: User[] = [];

	let newSemester = {
		name: '',
		year: '',
		period: '',
		startDate: '',
		endDate: ''
	};
	let creatingSemester = false;
	let semesterFormError = '';

	let newSubject = {
		name: '',
		code: '',
		professorId: '',
		description: ''
	};
	let creatingSubject = false;
	let subjectFormError = '';

	let newEnrollment = {
		studentId: '',
		subjectId: '',
		semesterId: ''
	};
	let creatingEnrollment = false;
	let enrollmentFormError = '';

	const roleLabels: Record<Role, string> = {
		student: 'Estudante',
		professor: 'Professor',
		admin: 'Administrador'
	};

	$: professorUsers = allUsers.filter((u) => u.role === 'professor');
	$: studentUsers = allUsers.filter((u) => u.role === 'student');

	function getRoleBadgeVariant(role: Role) {
		if (role === 'admin') return 'primary';
		if (role === 'professor') return 'success';
		return 'secondary';
	}

	function getStatus(user: User) {
		if (!user.requested_role || user.requested_role === user.role) {
			return { text: 'Alinhado', variant: 'secondary' as const };
		}
		return { text: 'Pendente', variant: 'warning' as const };
	}

	async function loadRoleRequests() {
		loading = true;
		error = '';
		try {
			const result = await api.getRoleRequests();
			if (!result.success) {
				error = result.error || 'Erro ao carregar solicitações de acesso';
				return;
			}
			roleRequests = ((result.data as any)?.users || []) as User[];
		} catch (err) {
			console.error('Failed to load role requests', err);
			error = 'Erro ao conectar com o servidor';
		} finally {
			loading = false;
		}
	}

	async function approveRole(user: User, role: Role) {
		updatingId = user.id;
		error = '';
		try {
			const result = await api.updateUserRole(user.id, role);
			if (!result.success) {
				error = result.error || 'Erro ao atualizar papel do usuário';
				return;
			}

			const updated = (result.data as any)?.user as User;
			roleRequests = roleRequests.filter((u) => u.id !== user.id);

			// If there are still other pending users, keep them; otherwise reload in case of drift
			if (roleRequests.length === 0) {
				await loadRoleRequests();
			}
		} catch (err) {
			console.error('Failed to approve role', err);
			error = 'Erro ao atualizar papel do usuário';
		} finally {
			updatingId = null;
		}
	}

	async function loadSemesters() {
		loading = true;
		error = '';
		try {
			const result = await api.getSemesters();
			if (!result.success) {
				error = result.error || 'Erro ao carregar semestres';
				return;
			}
			semesters = ((result.data as any)?.semesters || []) as Semester[];
		} catch (err) {
			console.error('Failed to load semesters', err);
			error = 'Erro ao conectar com o servidor';
		} finally {
			loading = false;
		}
	}

	async function loadSubjects() {
		loading = true;
		error = '';
		try {
			const result = await api.getSubjects();
			if (!result.success) {
				error = result.error || 'Erro ao carregar disciplinas';
				return;
			}
			subjects = ((result.data as any)?.subjects || []) as Subject[];
		} catch (err) {
			console.error('Failed to load subjects', err);
			error = 'Erro ao conectar com o servidor';
		} finally {
			loading = false;
		}
	}

	async function loadEnrollments() {
		loading = true;
		error = '';
		try {
			const result = await api.getEnrollments();
			if (!result.success) {
				error = result.error || 'Erro ao carregar matrículas';
				return;
			}
			enrollments = ((result.data as any)?.enrollments || []) as Enrollment[];
		} catch (err) {
			console.error('Failed to load enrollments', err);
			error = 'Erro ao conectar com o servidor';
		} finally {
			loading = false;
		}
	}

	async function loadUsersList() {
		loading = true;
		error = '';
		try {
			const result = await api.getAllUsers();
			if (!result.success) {
				error = result.error || 'Erro ao carregar usuários';
				return;
			}
			allUsers = ((result.data as any)?.users || []) as User[];
		} catch (err) {
			console.error('Failed to load users', err);
			error = 'Erro ao conectar com o servidor';
		} finally {
			loading = false;
		}
	}

	// Silent loaders used to populate dropdowns without resetting global loading/error state
	async function loadUsersListSilent() {
		try {
			const result = await api.getAllUsers();
			if (result.success) {
				allUsers = ((result.data as any)?.users || []) as User[];
			}
		} catch (err) {
			console.error('Failed to silently load users', err);
		}
	}

	async function loadSubjectsSilent() {
		try {
			const result = await api.getSubjects();
			if (result.success) {
				subjects = ((result.data as any)?.subjects || []) as Subject[];
			}
		} catch (err) {
			console.error('Failed to silently load subjects', err);
		}
	}

	async function loadSemestersSilent() {
		try {
			const result = await api.getSemesters();
			if (result.success) {
				semesters = ((result.data as any)?.semesters || []) as Semester[];
			}
		} catch (err) {
			console.error('Failed to silently load semesters', err);
		}
	}

	async function switchTab(tab: AdminTab) {
		currentTab = tab;
		if (tab === 'roleRequests') {
			await loadRoleRequests();
		} else if (tab === 'semesters') {
			await loadSemesters();
			// Also make sure we have users/subjects for cross-reference (silent)
			await loadUsersListSilent();
			await loadSubjectsSilent();
		} else if (tab === 'subjects') {
			await loadSubjects();
			// Need users list to show professor dropdown
			await loadUsersListSilent();
		} else if (tab === 'enrollments') {
			await loadEnrollments();
			// Need supporting data for dropdowns
			await Promise.all([
				loadUsersListSilent(),
				loadSubjectsSilent(),
				loadSemestersSilent()
			]);
		} else if (tab === 'users') {
			await loadUsersList();
		}
	}

	async function createSemester() {
		semesterFormError = '';
		if (!newSemester.name || !newSemester.year || !newSemester.period || !newSemester.startDate || !newSemester.endDate) {
			semesterFormError = 'Preencha todos os campos do semestre.';
			return;
		}

		const year = Number(newSemester.year);
		const period = Number(newSemester.period);

		if (!year || !period || (period !== 1 && period !== 2)) {
			semesterFormError = 'Ano ou período inválidos.';
			return;
		}

		const payload = {
			name: newSemester.name.trim(),
			year,
			period,
			start_date: new Date(newSemester.startDate).toISOString(),
			end_date: new Date(newSemester.endDate).toISOString()
		};

		creatingSemester = true;
		try {
			const result = await api.createSemester(payload);
			if (!result.success) {
				semesterFormError = result.error || 'Erro ao criar semestre';
				return;
			}
			const created = (result.data as any)?.semester as Semester;
			if (created) {
				semesters = [created, ...semesters];
			} else {
				await loadSemesters();
			}
			newSemester = { name: '', year: '', period: '', startDate: '', endDate: '' };
		} catch (err) {
			console.error('Failed to create semester', err);
			semesterFormError = 'Erro ao criar semestre';
		} finally {
			creatingSemester = false;
		}
	}

	async function activateSemester(semesterId: number) {
		try {
			const result = await api.activateSemester(String(semesterId));
			if (!result.success) {
				error = result.error || 'Erro ao ativar semestre';
				return;
			}
			await loadSemesters();
		} catch (err) {
			console.error('Failed to activate semester', err);
			error = 'Erro ao ativar semestre';
		}
	}

	async function createSubject() {
		subjectFormError = '';
		if (!newSubject.name || !newSubject.code || !newSubject.professorId) {
			subjectFormError = 'Nome, código e ID do professor são obrigatórios.';
			return;
		}

		const professorId = Number(newSubject.professorId);
		if (!professorId) {
			subjectFormError = 'ID do professor inválido.';
			return;
		}

		const payload = {
			name: newSubject.name.trim(),
			code: newSubject.code.trim(),
			description: newSubject.description.trim(),
			professor_id: professorId
		};

		creatingSubject = true;
		try {
			const result = await api.createSubject(payload);
			if (!result.success) {
				subjectFormError = result.error || 'Erro ao criar disciplina';
				return;
			}
			const created = (result.data as any)?.subject as Subject;
			if (created) {
				subjects = [created, ...subjects];
			} else {
				await loadSubjects();
			}
			newSubject = { name: '', code: '', professorId: '', description: '' };
		} catch (err) {
			console.error('Failed to create subject', err);
			subjectFormError = 'Erro ao criar disciplina';
		} finally {
			creatingSubject = false;
		}
	}

	async function createEnrollment() {
		enrollmentFormError = '';
		if (!newEnrollment.studentId || !newEnrollment.subjectId || !newEnrollment.semesterId) {
			enrollmentFormError = 'Informe IDs de estudante, disciplina e semestre.';
			return;
		}

		const studentId = Number(newEnrollment.studentId);
		const subjectId = Number(newEnrollment.subjectId);
		const semesterId = Number(newEnrollment.semesterId);

		if (!studentId || !subjectId || !semesterId) {
			enrollmentFormError = 'IDs informados são inválidos.';
			return;
		}

		const payload = {
			student_id: studentId,
			subject_id: subjectId,
			semester_id: semesterId
		};

		creatingEnrollment = true;
		try {
			const result = await api.createEnrollment(payload);
			if (!result.success) {
				enrollmentFormError = result.error || 'Erro ao criar matrícula';
				return;
			}
			const created = (result.data as any)?.enrollment as Enrollment;
			if (created) {
				enrollments = [created, ...enrollments];
			} else {
				await loadEnrollments();
			}
			newEnrollment = { studentId: '', subjectId: '', semesterId: '' };
		} catch (err) {
			console.error('Failed to create enrollment', err);
			enrollmentFormError = 'Erro ao criar matrícula';
		} finally {
			creatingEnrollment = false;
		}
	}

	async function refreshCurrentTab() {
		await switchTab(currentTab);
	}

	function getTabTitle(tab: AdminTab): string {
		if (tab === 'semesters') return 'Semestres';
		if (tab === 'subjects') return 'Disciplinas';
		if (tab === 'enrollments') return 'Matrículas';
		if (tab === 'users') return 'Usuários';
		return 'Solicitações de acesso';
	}

	function getTabDescription(tab: AdminTab): string {
		if (tab === 'semesters')
			return 'Cadastre e ative semestres letivos usados para vincular disciplinas e pesquisas.';
		if (tab === 'subjects')
			return 'Gerencie disciplinas e vincule cada uma ao professor responsável.';
		if (tab === 'enrollments')
			return 'Relacione estudantes, disciplinas e semestres para habilitar pesquisas.';
		if (tab === 'users')
			return 'Liste todos os usuários cadastrados e seus papéis atuais.';
		return 'Aprova pedidos de acesso como professor ou administrador. Até aprovação, todos continuam com acesso de estudante.';
	}

	onMount(() => {
		switchTab('roleRequests');
	});
</script>

<svelte:head>
	<title>Administração - Sistema de Consulta Discente</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">{getTabTitle(currentTab)}</h1>
			<p class="mt-1 text-gray-600">
				{getTabDescription(currentTab)}
			</p>
		</div>
		<Button variant="outline" onclick={refreshCurrentTab} disabled={loading}>
			Atualizar
		</Button>
	</div>

	<!-- Tabs -->
	<div class="border-b border-gray-200">
		<nav class="-mb-px flex space-x-4">
			<button
				type="button"
				class="whitespace-nowrap border-b-2 border-transparent px-3 py-2 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
				class:border-blue-500={currentTab === 'roleRequests'}
				class:text-blue-600={currentTab === 'roleRequests'}
				on:click={() => switchTab('roleRequests')}
			>
				Solicitações de acesso
			</button>
			<button
				type="button"
				class="whitespace-nowrap border-b-2 border-transparent px-3 py-2 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
				class:border-blue-500={currentTab === 'semesters'}
				class:text-blue-600={currentTab === 'semesters'}
				on:click={() => switchTab('semesters')}
			>
				Semestres
			</button>
			<button
				type="button"
				class="whitespace-nowrap border-b-2 border-transparent px-3 py-2 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
				class:border-blue-500={currentTab === 'subjects'}
				class:text-blue-600={currentTab === 'subjects'}
				on:click={() => switchTab('subjects')}
			>
				Disciplinas
			</button>
			<button
				type="button"
				class="whitespace-nowrap border-b-2 border-transparent px-3 py-2 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
				class:border-blue-500={currentTab === 'enrollments'}
				class:text-blue-600={currentTab === 'enrollments'}
				on:click={() => switchTab('enrollments')}
			>
				Matrículas
			</button>
			<button
				type="button"
				class="whitespace-nowrap border-b-2 border-transparent px-3 py-2 text-sm font-medium text-gray-500 hover:border-gray-300 hover:text-gray-700"
				class:border-blue-500={currentTab === 'users'}
				class:text-blue-600={currentTab === 'users'}
				on:click={() => switchTab('users')}
			>
				Usuários
			</button>
		</nav>
	</div>

	<!-- Role requests tab -->
	{#if currentTab === 'roleRequests'}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center">
					<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
					<p class="mt-4 text-gray-600">Carregando solicitações...</p>
				</div>
			</div>
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-6 text-center">
					<p class="text-red-800">{error}</p>
					<Button variant="outline" class="mt-4" onclick={refreshCurrentTab}>
						Tentar novamente
					</Button>
				</div>
			</Card>
		{:else if roleRequests.length === 0}
			<Card>
				<div class="py-8 text-center">
					<p class="text-gray-600">
						Nenhuma solicitação pendente no momento. Novos pedidos aparecerão aqui
						quando os usuários se cadastrarem.
					</p>
				</div>
			</Card>
		{:else}
			<div class="rounded-lg border border-gray-200 bg-white">
				<div class="max-h-[480px] overflow-y-auto overflow-x-auto">
					<table class="min-w-full divide-y divide-gray-200">
						<thead class="bg-gray-50">
							<tr>
								<th
									scope="col"
									class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Usuário
								</th>
								<th
									scope="col"
									class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Email
								</th>
								<th
									scope="col"
									class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Papel atual
								</th>
								<th
									scope="col"
									class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Papel solicitado
								</th>
								<th
									scope="col"
									class="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Status
								</th>
								<th
									scope="col"
									class="px-6 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500"
								>
									Ações
								</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-200 bg-white">
							{#each roleRequests as user}
								{@const status = getStatus(user)}
								<tr class="hover:bg-gray-50">
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-900">
										<div class="font-medium">
											{user.first_name} {user.last_name}
										</div>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm text-gray-500">
										{user.email}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm">
										<Badge variant={getRoleBadgeVariant(user.role)}>
											{roleLabels[user.role]}
										</Badge>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm">
										{#if user.requested_role && user.requested_role !== user.role}
											<Badge variant={getRoleBadgeVariant(user.requested_role)}>
												{roleLabels[user.requested_role]}
											</Badge>
										{:else}
											<span class="text-gray-500">—</span>
										{/if}
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-sm">
										<Badge variant={status.variant}>
											{status.text}
										</Badge>
									</td>
									<td class="whitespace-nowrap px-6 py-4 text-right text-sm">
										<div class="flex justify-end gap-2">
											<Button
												size="sm"
												variant="outline"
												disabled={updatingId === user.id}
												onclick={() => approveRole(user, 'student')}
											>
												Manter estudante
											</Button>
											{#if user.requested_role && user.requested_role !== 'student'}
												<Button
													size="sm"
													disabled={updatingId === user.id}
													onclick={() => approveRole(user, user.requested_role as Role)}
												>
													Aprovar {roleLabels[user.requested_role as Role]}
												</Button>
											{/if}
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}
	{/if}

	<!-- Semesters tab -->
	{#if currentTab === 'semesters'}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center">
					<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
					<p class="mt-4 text-gray-600">Carregando semestres...</p>
				</div>
			</div>
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-6 text-center">
					<p class="text-red-800">{error}</p>
					<Button variant="outline" class="mt-4" onclick={refreshCurrentTab}>
						Tentar novamente
					</Button>
				</div>
			</Card>
		{:else}
			<div class="space-y-6">
				<Card>
					<div class="space-y-4">
						<h2 class="text-lg font-semibold text-gray-900">Criar novo semestre</h2>
						{#if semesterFormError}
							<p class="text-sm text-red-600">{semesterFormError}</p>
						{/if}
						<div class="grid gap-4 md:grid-cols-2">
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Nome do semestre
								</label>
								<input
									type="text"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									placeholder="Ex: 2025.1"
									bind:value={newSemester.name}
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">Ano</label>
								<input
									type="number"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									placeholder="Ex: 2025"
									bind:value={newSemester.year}
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Período
								</label>
								<select
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSemester.period}
								>
									<option value="">Selecione</option>
									<option value="1">1</option>
									<option value="2">2</option>
								</select>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Data de início
								</label>
								<input
									type="date"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSemester.startDate}
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Data de término
								</label>
								<input
									type="date"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSemester.endDate}
								/>
							</div>
						</div>
						<div class="flex justify-end">
							<Button onclick={createSemester} disabled={creatingSemester}>
								{#if creatingSemester}
									Criando...
								{:else}
									Criar semestre
								{/if}
							</Button>
						</div>
					</div>
				</Card>

				<Card>
					{#if semesters.length === 0}
						<div class="py-6 text-center text-sm text-gray-600">
							Nenhum semestre cadastrado ainda.
						</div>
					{:else}
						<div class="rounded-lg border border-gray-200">
							<div class="max-h-[420px] overflow-y-auto overflow-x-auto">
								<table class="min-w-full divide-y divide-gray-200">
									<thead class="bg-gray-50">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Nome
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Ano/Período
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Início
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Término
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Status
											</th>
											<th class="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-gray-500">
												Ações
											</th>
										</tr>
									</thead>
									<tbody class="divide-y divide-gray-200 bg-white">
										{#each semesters as semester}
											<tr class="hover:bg-gray-50">
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
													{semester.name}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{semester.year}.{semester.period}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{semester.start_date?.slice(0, 10)}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{semester.end_date?.slice(0, 10)}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm">
													{#if semester.is_active}
														<Badge variant="success">Ativo</Badge>
													{:else}
														<Badge variant="secondary">Inativo</Badge>
													{/if}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-right text-sm">
													{#if !semester.is_active}
														<Button
															size="sm"
															variant="outline"
															onclick={() => activateSemester(semester.id)}
														>
															Ativar
														</Button>
													{:else}
														<span class="text-gray-500">Sem ações</span>
													{/if}
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						</div>
					{/if}
				</Card>
			</div>
		{/if}
	{/if}

	<!-- Subjects tab -->
	{#if currentTab === 'subjects'}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center">
					<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
					<p class="mt-4 text-gray-600">Carregando disciplinas...</p>
				</div>
			</div>
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-6 text-center">
					<p class="text-red-800">{error}</p>
					<Button variant="outline" class="mt-4" onclick={refreshCurrentTab}>
						Tentar novamente
					</Button>
				</div>
			</Card>
		{:else}
			<div class="space-y-6">
				<Card>
					<div class="space-y-4">
						<h2 class="text-lg font-semibold text-gray-900">Cadastrar nova disciplina</h2>
						{#if subjectFormError}
							<p class="text-sm text-red-600">{subjectFormError}</p>
						{/if}
						<div class="grid gap-4 md:grid-cols-2">
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Nome da disciplina
								</label>
								<input
									type="text"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSubject.name}
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Código
								</label>
								<input
									type="text"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									placeholder="Ex: MAC0123"
									bind:value={newSubject.code}
								/>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									ID do professor responsável
								</label>
								<select
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSubject.professorId}
								>
									<option value="">Selecione um professor</option>
									{#each professorUsers as professor}
										<option value={String(professor.id)}>
											{professor.first_name} {professor.last_name} (ID {professor.id})
										</option>
									{/each}
								</select>
								<p class="mt-1 text-xs text-gray-500">
									Apenas usuários com papel de professor aparecem aqui.
								</p>
							</div>
							<div class="md:col-span-2">
								<label class="mb-1 block text-sm font-medium text-gray-700">
									Descrição
								</label>
								<textarea
									rows="2"
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newSubject.description}
								/>
							</div>
						</div>
						<div class="flex justify-end">
							<Button onclick={createSubject} disabled={creatingSubject}>
								{#if creatingSubject}
									Cadastrando...
								{:else}
									Cadastrar disciplina
								{/if}
							</Button>
						</div>
					</div>
				</Card>

				<Card>
					{#if subjects.length === 0}
						<div class="py-6 text-center text-sm text-gray-600">
							Nenhuma disciplina cadastrada ainda.
						</div>
					{:else}
						<div class="rounded-lg border border-gray-200">
							<div class="max-h-[420px] overflow-y-auto overflow-x-auto">
								<table class="min-w-full divide-y divide-gray-200">
									<thead class="bg-gray-50">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Nome
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Código
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Professor
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												ID Professor
											</th>
										</tr>
									</thead>
									<tbody class="divide-y divide-gray-200 bg-white">
										{#each subjects as subject}
											<tr class="hover:bg-gray-50">
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
													{subject.name}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{subject.code}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{#if subject.professor}
														{subject.professor.first_name} {subject.professor.last_name}
													{:else}
														<span class="text-gray-400">Professor não carregado</span>
													{/if}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
													{subject.professor_id}
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						</div>
					{/if}
				</Card>
			</div>
		{/if}
	{/if}

	<!-- Enrollments tab -->
	{#if currentTab === 'enrollments'}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center">
					<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
					<p class="mt-4 text-gray-600">Carregando matrículas...</p>
				</div>
			</div>
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-6 text-center">
					<p class="text-red-800">{error}</p>
					<Button variant="outline" class="mt-4" onclick={refreshCurrentTab}>
						Tentar novamente
					</Button>
				</div>
			</Card>
		{:else}
			<div class="space-y-6">
				<Card>
					<div class="space-y-4">
						<h2 class="text-lg font-semibold text-gray-900">Criar nova matrícula</h2>
						{#if enrollmentFormError}
							<p class="text-sm text-red-600">{enrollmentFormError}</p>
						{/if}
						<div class="grid gap-4 md:grid-cols-3">
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									ID do estudante
								</label>
								<select
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newEnrollment.studentId}
								>
									<option value="">Selecione um estudante</option>
									{#each studentUsers as student}
										<option value={String(student.id)}>
											{student.first_name} {student.last_name} (ID {student.id})
										</option>
									{/each}
								</select>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									ID da disciplina
								</label>
								<select
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newEnrollment.subjectId}
								>
									<option value="">Selecione uma disciplina</option>
									{#each subjects as subject}
										<option value={String(subject.id)}>
											{subject.name} ({subject.code})
										</option>
									{/each}
								</select>
							</div>
							<div>
								<label class="mb-1 block text-sm font-medium text-gray-700">
									ID do semestre
								</label>
								<select
									class="focus:ring-primary focus:border-primary block w-full rounded-md border border-gray-300 px-3 py-2 text-sm text-gray-900"
									bind:value={newEnrollment.semesterId}
								>
									<option value="">Selecione um semestre</option>
									{#each semesters as semester}
										<option value={String(semester.id)}>
											{semester.name} ({semester.year}.{semester.period})
										</option>
									{/each}
								</select>
							</div>
						</div>
						<p class="text-xs text-gray-500">
							As listas são preenchidas automaticamente a partir das abas de usuários,
							disciplinas e semestres.
						</p>
						<div class="flex justify-end">
							<Button onclick={createEnrollment} disabled={creatingEnrollment}>
								{#if creatingEnrollment}
									Criando...
								{:else}
									Criar matrícula
								{/if}
							</Button>
						</div>
					</div>
				</Card>

				<Card>
					{#if enrollments.length === 0}
						<div class="py-6 text-center text-sm text-gray-600">
							Nenhuma matrícula cadastrada ainda.
						</div>
					{:else}
						<div class="rounded-lg border border-gray-200">
							<div class="max-h-[420px] overflow-y-auto overflow-x-auto">
								<table class="min-w-full divide-y divide-gray-200">
									<thead class="bg-gray-50">
										<tr>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Estudante
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Disciplina
											</th>
											<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
												Semestre
											</th>
										</tr>
									</thead>
									<tbody class="divide-y divide-gray-200 bg-white">
										{#each enrollments as enrollment}
											<tr class="hover:bg-gray-50">
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
													{#if enrollment.student}
														{enrollment.student.first_name}
														{enrollment.student.last_name}
														<span class="ml-1 text-xs text-gray-500">
															(ID {enrollment.student_id})
														</span>
													{:else}
														<span class="text-gray-700">
															ID {enrollment.student_id}
														</span>
													{/if}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
													{#if enrollment.subject}
														{enrollment.subject.name}
														<span class="ml-1 text-xs text-gray-500">
															({enrollment.subject.code})
														</span>
													{:else}
														<span class="text-gray-700">
															Disciplina {enrollment.subject_id}
														</span>
													{/if}
												</td>
												<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
													{#if enrollment.semester}
														{enrollment.semester.name}
													{:else}
														<span class="text-gray-700">
															Semestre {enrollment.semester_id}
														</span>
													{/if}
												</td>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						</div>
					{/if}
				</Card>
			</div>
		{/if}
	{/if}

	<!-- Users tab -->
	{#if currentTab === 'users'}
		{#if loading}
			<div class="flex items-center justify-center py-12">
				<div class="text-center">
					<div class="mx-auto h-12 w-12 animate-spin rounded-full border-4 border-blue-200 border-t-blue-600"></div>
					<p class="mt-4 text-gray-600">Carregando usuários...</p>
				</div>
			</div>
		{:else if error}
			<Card class="border-red-200 bg-red-50">
				<div class="py-6 text-center">
					<p class="text-red-800">{error}</p>
					<Button variant="outline" class="mt-4" onclick={refreshCurrentTab}>
						Tentar novamente
					</Button>
				</div>
			</Card>
		{:else if allUsers.length === 0}
			<Card>
				<div class="py-8 text-center">
					<p class="text-gray-600">Nenhum usuário cadastrado.</p>
				</div>
			</Card>
		{:else}
			<div class="rounded-lg border border-gray-200 bg-white">
				<div class="max-h-[520px] overflow-y-auto overflow-x-auto">
					<table class="min-w-full divide-y divide-gray-200">
						<thead class="bg-gray-50">
							<tr>
								<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
									ID
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
									Nome
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
									Email
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
									Papel atual
								</th>
								<th class="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">
									Papel solicitado
								</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-gray-200 bg-white">
							{#each allUsers as user}
								<tr class="hover:bg-gray-50">
									<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
										{user.id}
									</td>
									<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-900">
										{user.first_name} {user.last_name}
									</td>
									<td class="whitespace-nowrap px-4 py-3 text-sm text-gray-700">
										{user.email}
									</td>
									<td class="whitespace-nowrap px-4 py-3 text-sm">
										<Badge variant={getRoleBadgeVariant(user.role)}>
											{roleLabels[user.role]}
										</Badge>
									</td>
									<td class="whitespace-nowrap px-4 py-3 text-sm">
										{#if user.requested_role && user.requested_role !== user.role}
											<Badge variant={getRoleBadgeVariant(user.requested_role)}>
												{roleLabels[user.requested_role]}
											</Badge>
										{:else}
											<span class="text-gray-500">—</span>
										{/if}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		{/if}
	{/if}
</div>
