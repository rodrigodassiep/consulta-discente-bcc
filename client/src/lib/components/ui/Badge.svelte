<script lang="ts">
	interface Props {
		variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger';
		size?: 'sm' | 'md';
		children: any;
	}

	let { variant = 'primary', size = 'md', children }: Props = $props();

	const baseClasses = 'inline-flex items-center rounded-full font-medium';
	
	const variantClasses = {
		primary: 'bg-blue-100 text-blue-800',
		secondary: 'bg-gray-100 text-gray-800',
		success: 'bg-green-100 text-green-800',
		warning: 'bg-yellow-100 text-yellow-800',
		danger: 'bg-red-100 text-red-800'
	};

	const sizeClasses = {
		sm: 'px-2 py-1 text-xs',
		md: 'px-3 py-1 text-sm'
	};

	const classes = `${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]}`;
</script>

<span class={classes}>
	{@render children()}
</span> 