<script lang="ts">
	import { m } from '$lib/paraglide/messages';

	let mood: number | null = null; // 1–5
	let emotions = '';
	let diary = '';

	$: canSubmit = mood !== null && emotions.trim() && diary.trim();
</script>

<div class="flex min-h-screen flex-col items-center gap-6 px-4 py-8 font-mono text-white">
	<h1 class="text-center text-2xl">{m.start_your_day()}</h1>

	<!-- Mood Rating -->
	<div class="flex gap-4">
		{#each [1, 2, 3, 4, 5] as n}
			<button
				class="border border-white px-4 py-2 transition-none"
				class:bg-white={mood === n}
				class:text-black={mood === n}
				class:text-white={mood !== n}
				on:click={() => (mood = n)}
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
		class="w-full max-w-md border border-white bg-black px-4 py-2 text-white placeholder-white"
	/>

	<!-- Diary -->
	<textarea
		placeholder={m.start_today()}
		bind:value={diary}
		class="h-40 w-full max-w-md border border-white bg-black px-4 py-2 text-white placeholder-white"
	></textarea>

	<!-- Submit -->
	<a
		href="/signup"
		class="my-2 border border-white px-6 py-2 transition-none"
		class:bg-white={canSubmit}
		class:text-black={canSubmit}
		class:text-white={!canSubmit}
		class:opacity-50={!canSubmit}
		class:pointer-events-none={!canSubmit}
		class:hover:bg-black={canSubmit}
		class:hover:text-white={canSubmit}
	>
		{m.start_process()}</a
	>
</div>
