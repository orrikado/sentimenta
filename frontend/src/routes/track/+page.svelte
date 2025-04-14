<script lang="ts">
	import { onMount } from 'svelte';
	import { getMonthDays } from '$lib/calendar-utils';
	import Modal from '$lib/components/Modal.svelte';

	let today = new Date();
	let currentMonth = today.getMonth();
	let currentYear = today.getFullYear();

	let days: Date[] = $state([]);
	let showModal = $state(false);
	let selectedDate: Date = $state(new Date());

	onMount(() => {
		days = getMonthDays(currentYear, currentMonth);
	});

	import { m } from '$lib/paraglide/messages';

	let mood = $state<number | null>(null);
	// let emotions = '';
	let diary = $state('');
	let formError = $state<string | null>(null);
	let formSuccess = $state<boolean>(false);

	const canSubmit = $derived(() => !!(mood !== null && diary.trim()));

	$effect(() => {
		if (showModal) {
			formSuccess = false;
			mood = null;
			diary = '';
			formError = null;
		}
	});
</script>

<h1 class="my-4 text-center text-2xl font-bold">{currentYear}/{currentMonth + 1}</h1>
<div class="mx-auto my-4 md:max-w-3/5">
	<div
		class="grid grid-cols-7 gap-0.5 border border-stone-300 bg-white p-1 text-sm text-black md:p-2 dark:border-white/10 dark:bg-stone-900 dark:text-white"
	>
		{#each ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'] as day (day)}
			<div class="border-b border-white/10 p-2 text-center font-bold">{day}</div>
		{/each}

		{#each days as date (date)}
			<button
				class="flex aspect-square items-center justify-center border-black/10 p-2 text-center text-xl dark:border-white/10"
				class:dark:bg-stone-800={date instanceof Date && !isNaN(date.getDate())}
				class:bg-stone-100={date instanceof Date && !isNaN(date.getDate())}
				class:border={date instanceof Date && !isNaN(date.getDate())}
				onclick={() => {
					selectedDate = date;
					showModal = true;
				}}
			>
				{#if date instanceof Date && !isNaN(date.getDate())}
					{date.getDate()}
				{/if}
			</button>
		{/each}
	</div>
</div>

<Modal bind:showModal>
	{#snippet header()}
		<h2 class="text-center font-bold">log for {selectedDate.toDateString()}</h2>
	{/snippet}

	<form
		class="flex w-full max-w-md flex-col gap-6"
		onsubmit={async () => {
			if (canSubmit()) {
				let result = await fetch('/api/moods/add', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({ score: mood, description: diary, date: selectedDate })
				});
				if (!result.ok) {
					const errorData = result;
					console.error('Error:', errorData);
					formError = m.error_occured();
				} else {
					formSuccess = true;
				}
			}
		}}
	>
		<h1 class="text-center text-2xl">{m.start_your_day()}</h1>

		<!-- Mood Rating -->
		<div class="flex justify-center gap-4">
			{#each [1, 2, 3, 4, 5] as n (n)}
				<button
					type="button"
					class="border border-current px-3 py-2"
					class:bg-white={mood === n}
					class:text-black={mood === n}
					class:text-current={mood !== n}
					onclick={() => (mood = n)}
				>
					{n}
				</button>
			{/each}
		</div>

		<!-- Diary -->
		<textarea
			placeholder={m.start_today()}
			bind:value={diary}
			class="h-40 w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		></textarea>

		<!-- Submit -->
		<button
			type="submit"
			class="border border-current px-6 py-2 text-center"
			class:bg-white={canSubmit()}
			class:text-black={canSubmit()}
			class:opacity-50={!canSubmit()}
			class:pointer-events-none={!canSubmit()}
		>
			{m.start_process()}
		</button>

		{#if formError}
			<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
		{/if}
		{#if formSuccess}
			<p class="text-sm text-green-500 dark:text-green-400">{m.mood_upload_success()}</p>
		{/if}
	</form>
</Modal>
