<script lang="ts">
	import { goto } from '$app/navigation';
	import { m } from '$lib/paraglide/messages.js';
	import { userId } from '$lib/stores/user';
	import { onMount } from 'svelte';

	onMount(() => {
		if ($userId) {
			goto('/track');
		}
	});
</script>

<svelte:head>
	<title>Sentimenta</title>
</svelte:head>

<div
	class="flex min-h-screen flex-col items-center justify-center gap-4 bg-stone-100 px-4 text-center font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<h1 class="text-2xl">{m.ai_judge()}</h1>
	<p class="max-w-md text-base">
		{m.frontpage_ai_notice()}
	</p>
	<a
		href="/start"
		class="border border-black px-4 py-2 transition-none hover:bg-black hover:text-white dark:border-white dark:hover:bg-white dark:hover:text-black"
	>
		{m.ill_try()}
	</a>
</div>
