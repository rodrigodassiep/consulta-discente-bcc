<script lang="ts">
	import '../app.css';
	import { pwaInfo } from 'virtual:pwa-info';

	let { children } = $props();
</script>

<svelte:head>
	{#if pwaInfo}
		<link rel="manifest" href={pwaInfo.webManifest.href} />
	{/if}
</svelte:head>

<div class="app">
	<main>
		{@render children()}
	</main>
</div>

<style>
	:root {
    --color-primary: #234aac;
    --color-primary-focus: #2a88c6;
  }
	.app {
		display: flex;
		flex-direction: column;
		/* min-height: 100vh; */
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 64rem;
		margin: 0 auto;
		box-sizing: border-box;
	}

</style>
