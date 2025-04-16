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
	type MoodEntry = {
		date: Date;
		score: number;
		description: string;
		emotions: string;
	};

	let moods = $state<MoodEntry[]>([]);
	let moodMap = $state<Map<string, MoodEntry>>(new Map());

	onMount(async () => {
		if (!$userId) {
			goto('/login');
			return;
		}
		days = getMonthDays(currentYear, currentMonth);

		try {
			const res = await fetch('/api/moods/get');
			if (res.ok) {
				const data = await res.json();
				const parsed = data.map((m: { date: string | number | Date }) => ({
					...m,
					date: new Date(m.date)
				}));

				moods = parsed;
				moodMap = new Map(moods.map((m) => [m.date.toString(), m]));
			} else {
				console.error('Failed to fetch moods');
			}
		} catch (e) {
			console.error('Network error:', e);
		}
	});

	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import { goto } from '$app/navigation';

	let mood = $state<number>(0);
	let emotions = $state('');
	let diary = $state('');
	let formError = $state<string | null>(null);
	let formSuccess = $state<boolean>(false);

	const canSubmit = $derived(() => !!(mood !== null && diary.trim() && emotions.trim()));

	$effect(() => {
		if (showModal) {
			formSuccess = false;
			mood = 0;
			diary = '';
			formError = null;
			emotions = '';
		}
	});
	$inspect(moods);
	$inspect(moodMap);

	function parseEmotions(input: string): string {
		return input
			.toLowerCase()
			.split(/[\s,|./\\;:!?&]+|(?:\band\b)/gi) // split on all kinds of delimiters
			.map((word) => word.trim())
			.filter((word) => word.length > 0)
			.join(',');
	}

	function getDayClass(date: Date | null, moods: Map<string, MoodEntry>) {
		if (date instanceof Date && !isNaN(date.getDate())) {
			if (moods.has(date.toString())) {
				switch (moods.get(date.toString())?.score) {
					case 1:
						return 'bg-red-300 dark:bg-red-900';
					case 2:
						return 'bg-orange-300 dark:bg-orange-900';
					case 3:
						return 'bg-yellow-400 dark:bg-yellow-900';
					case 4:
						return 'bg-green-300 dark:bg-green-900';
					case 5:
						return 'bg-blue-300 dark:bg-blue-900';
				}
			} else {
				return 'bg-stone-100 dark:bg-stone-800';
			}
		} else {
			return '';
		}
	}
</script>

<h1 class="my-4 text-center text-2xl font-bold">{currentYear}/{currentMonth + 1}</h1>
<div class="mx-auto my-4 md:max-w-3/5">
	<div
		class="grid grid-cols-7 gap-0.5 border border-stone-300 bg-white p-1 text-sm text-black md:p-2 dark:border-white/10 dark:bg-stone-900 dark:text-white"
	>
		{#each [m.sun(), m.mon(), m.tue(), m.wed(), m.thu(), m.fri(), m.sat()] as day (day)}
			<div class="border-b border-white/10 p-2 text-center font-bold">{day}</div>
		{/each}

		{#each days as date (date)}
			<button
				class={`${getDayClass(date, moodMap)} flex aspect-square items-center justify-center border-black/10 p-2 text-center text-xl dark:border-white/10`}
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
		<h2 class="text-center font-bold">{m.log_for()} {selectedDate.toLocaleDateString()}</h2>
	{/snippet}

	<form
		class="flex w-full max-w-md flex-col gap-6"
		onsubmit={async () => {
			if (canSubmit()) {
				let nextDay = new Date(selectedDate);
				nextDay.setDate(nextDay.getDate() + 1);
				let result = await fetch('/api/moods/add', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						score: mood,
						description: diary,
						date: nextDay,
						emotions: parseEmotions(emotions)
					})
				});
				if (!result.ok) {
					const errorData = result;
					console.error('Error:', errorData);
					formError = m.error_occured();
				} else {
					const newMood: MoodEntry = {
						date: selectedDate,
						score: mood,
						description: diary,
						emotions: parseEmotions(emotions)
					};

					moods = [...moods, newMood];
					moodMap = new Map(moodMap).set(selectedDate.toString(), newMood);

					formError = null;
					formSuccess = true;
					showModal = false;
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
