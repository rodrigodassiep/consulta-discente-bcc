<script lang="ts">
	import { onMount } from 'svelte';
	import Card from '$lib/components/ui/Card.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Badge from '$lib/components/ui/Badge.svelte';
	import { api } from '$lib/api.js';

	type Role = 'student' | 'professor' | 'admin';

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

	let loading = true;
	let error = '';
	let updatingId: number | null = null;
	let users: User[] = [];

	const roleLabels: Record<Role, string> = {
		student: 'Estudante',
		professor: 'Professor',
		admin: 'Administrador'
	};

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
			users = ((result.data as any)?.users || []) as User[];
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
			users = users.filter((u) => u.id !== user.id);

			// If there are still other pending users, keep them; otherwise reload in case of drift
			if (users.length === 0) {
				await loadRoleRequests();
			}
		} catch (err) {
			console.error('Failed to approve role', err);
			error = 'Erro ao atualizar papel do usuário';
		} finally {
			updatingId = null;
		}
	}

	onMount(() => {
		loadRoleRequests();
	});
</script>

<svelte:head>
	<title>Solicitações de Acesso - Admin</title>
</svelte:head>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold text-gray-900">Solicitações de acesso</h1>
			<p class="mt-1 text-gray-600">
				Aqui você aprova pedidos de acesso como professor ou administrador. Até aprovação,
				todos continuam com acesso de estudante.
			</p>
		</div>
		<Button variant="outline" onclick={loadRoleRequests} disabled={loading}>
			Atualizar
		</Button>
	</div>

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
				<Button variant="outline" class="mt-4" onclick={loadRoleRequests}>
					Tentar novamente
				</Button>
			</div>
		</Card>
	{:else if users.length === 0}
		<Card>
			<div class="py-8 text-center">
				<p class="text-gray-600">
					Nenhuma solicitação pendente no momento. Novos pedidos aparecerão aqui quando os
					usuários se cadastrarem.
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
						{#each users as user}
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
</div>
