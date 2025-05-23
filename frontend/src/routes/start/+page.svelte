<script lang="ts">
	import { m } from '$lib/paraglide/messages';
	import { onMount } from 'svelte';

	let mood = $state(null);
	let emotions = $state('');
	let diary = $state('');

	let canSubmit = $derived(!!(mood !== null && emotions.trim() && diary.trim()));

	onMount(() => {
		const saved = localStorage.getItem('formData');
		if (saved) {
			try {
				const parsed = JSON.parse(saved);
				mood = parsed.mood ?? null;
				emotions = parsed.emotions || '';
				diary = parsed.diary || '';
			} catch (e) {
				console.error('Failed to parse saved data', e);
			}
		}
	});

	$effect(() => {
		localStorage.setItem('formData', JSON.stringify({ mood, emotions, diary }));
	});
</script>

<svelte:head>
	<title>Sentimenta | Start</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 p-4 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<div class="flex w-full max-w-md flex-col gap-6">
		<h1 class="text-center text-2xl">{m.start_your_day()}</h1>

		<!-- Mood Rating -->
		<div class="flex justify-center gap-4">
			{#each [1, 2, 3, 4, 5] as n (n)}
				<button
					type="button"
					role="radio"
					aria-checked={mood === n}
					class="cursor-pointer border border-current px-3 py-2"
					class:dark:bg-white={mood === n}
					class:bg-stone-400={mood === n}
					class:text-black={mood === n}
					class:text-current={mood !== n}
					onclick={() => (mood = n)}
					aria-label={m.aria_select_mood({ mood: n })}
				>
					{n}
				</button>
			{/each}
		</div>

		<!-- Emotions -->
		<input
			type="text"
			placeholder={m.start_emotions_placeholder()}
			bind:value={emotions}
			class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		/>

		<!-- Diary -->
		<textarea
			placeholder={m.start_today()}
			bind:value={diary}
			class="h-40 w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		></textarea>

		<!-- Submit -->
		<a
			href="/signup"
			class="border border-current px-6 py-2 text-center"
			class:bg-white={canSubmit}
			class:text-black={canSubmit}
			class:opacity-50={!canSubmit}
			class:pointer-events-none={!canSubmit}
		>
			{m.start_process()}</a
		>
	</div>
</div>
