<script lang="ts">
	interface Props {
		variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		loading?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: () => void;
		children: any;
	}

	let {
		variant = 'primary',
		size = 'md',
		disabled = false,
		loading = false,
		type = 'button',
		onclick,
		children,
		...restProps
	}: Props = $props();

	const baseClasses = 'inline-flex items-center justify-center rounded-lg font-semibold transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:pointer-events-none shadow-sm hover:shadow-md';

	const variantClasses = {
		primary: 'bg-gradient-to-br from-blue-700 to-blue-800 text-white hover:from-blue-800 hover:to-blue-900 focus:ring-blue-500 border border-blue-800',
		secondary: 'bg-gradient-to-br from-amber-500 to-amber-600 text-white hover:from-amber-600 hover:to-amber-700 focus:ring-amber-500 border border-amber-600',
		outline: 'border-2 border-blue-700 bg-white text-blue-700 hover:bg-blue-50 focus:ring-blue-500',
		ghost: 'text-gray-700 hover:bg-gray-100 focus:ring-gray-400 shadow-none hover:shadow-none'
	};

	const sizeClasses = {
		sm: 'px-3 py-1.5 text-sm',
		md: 'px-5 py-2.5 text-sm',
		lg: 'px-8 py-3.5 text-base'
	};

	const classes = `${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]}`;
</script>

<button
	{type}
	class={classes}
	{disabled}
	on:click={onclick}
	{...restProps}
>
	{#if loading}
		<svg
			class="-ml-1 mr-2 h-4 w-4 animate-spin"
			xmlns="http://www.w3.org/2000/svg"
			fill="none"
			viewBox="0 0 24 24"
		>
			<circle
				class="opacity-25"
				cx="12"
				cy="12"
				r="10"
				stroke="currentColor"
				stroke-width="4"
			></circle>
			<path
				class="opacity-75"
				fill="currentColor"
				d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
			></path>
		</svg>
	{/if}
	{@render children()}
</button> 