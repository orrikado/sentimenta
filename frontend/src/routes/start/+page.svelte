<script lang="ts">
	import { m } from '$lib/paraglide/messages';

	let mood: number | null = null; // 1–5
	let emotions = '';
	let diary = '';

	$: canSubmit = mood !== null && emotions.trim() && diary.trim();
</script>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 p-4 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<div class="flex w-full max-w-md flex-col gap-6">
		<h1 class="text-center text-2xl">{m.start_your_day()}</h1>

		<!-- Mood Rating -->
		<div class="flex justify-center gap-4">
			{#each [1, 2, 3, 4, 5] as n}
				<button
					class="border border-current px-3 py-2"
					class:bg-white={mood === n}
					class:text-black={mood === n}
					class:text-current={mood !== n}
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
