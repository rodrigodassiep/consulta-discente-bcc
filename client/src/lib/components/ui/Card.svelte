<script lang="ts">
	interface Props {
		class?: string;
		children: any;
	}

	let { class: className = '', children }: Props = $props();
</script>

<div class="rounded-lg border border-gray-200 bg-white p-6 shadow-sm {className}">
	{@render children()}
</div> 