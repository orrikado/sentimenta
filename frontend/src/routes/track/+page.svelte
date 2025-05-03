<script lang="ts">
	// External imports
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getMonthDays } from '$lib/calendar-utils';
	import { m } from '$lib/paraglide/messages';
	import { refreshUserId } from '$lib/user';
	import { userId } from '$lib/stores/user';
	import Modal from '$lib/components/Modal.svelte';

	// State variables
	let today = new Date();
	let currentMonth = $state(today.getMonth());
	let currentYear = $state(today.getFullYear());

	let days: Date[] = $derived(getMonthDays(currentYear, currentMonth));
	let showModal = $state(false);
	let selectedDate: Date = $state(new Date());
	let submitInProcess = $state(false);
	let moods = $state<MoodEntry[]>([]);

	let mood = $state<number>(0);
	let emotions = $state('');
	let diary = $state('');
	let formError = $state<string | null>(null);
	let formSuccess = $state<boolean>(false);

	// Types
	type MoodEntry = {
		uid: number | undefined;
		date: Date;
		score: number;
		description: string;
		emotions: string;
	};

	// Derived values
	const getDateKey = (date: Date) =>
		`${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;

	let moodMap = $derived(new Map(moods.map((m) => [getDateKey(m.date), m])));
	const canSubmit = $derived(() => {
		const future = new Date(selectedDate) > today;
		// Ensure the selected date is today or in the past
		return !!(mood !== 0 && diary.trim() && emotions.trim() && !future);
	});

	// Effects
	$effect(() => {
		if (showModal && !submitInProcess) {
			formSuccess = false;
			formError = null;

			if (moodMap.has(getDateKey(selectedDate))) {
				const entry = moodMap.get(getDateKey(selectedDate));
				mood = entry?.score || 0;
				diary = entry?.description || '';
				emotions = entry?.emotions || '';
			} else {
				mood = 0;
				diary = '';
				emotions = '';
			}
		}
	});

	// Lifecycle
	onMount(async () => {
		if (!$userId) {
			goto('/login');
			return;
		}
		await updateMoods();
	});

	// Functions
	async function updateMoods() {
		try {
			const res = await fetch('/api/moods/get');
			if (res.ok) {
				const data = await res.json();
				const parsed = data.map((m: { date: string | number | Date }) => ({
					...m,
					date: new Date(m.date)
				}));
				moods = parsed;
			} else {
				console.error('Failed to fetch moods');
				refreshUserId();
				if (!$userId) goto('/login');
			}
		} catch (e) {
			console.error('Network error:', e);
		}
	}

	function parseEmotions(input: string): string {
		return input
			.toLowerCase()
			.split(/[\s,|./\\;:!?&]+|(?:\band\b)/gi)
			.map((word) => word.trim())
			.filter((word) => word.length > 0)
			.join(',');
	}

	function getDayClass(date: Date | null, moods: Map<string, MoodEntry>) {
		let classes = '';
		if (date instanceof Date && !isNaN(date.getDate())) {
			const dateKey = getDateKey(date);
			const todayKey = getDateKey(new Date());

			if (dateKey === todayKey) {
				classes += 'dark:border-white border-black ';
			} else {
				classes += 'border-black/10  dark:border-white/10 ';
			}

			if (moods.has(getDateKey(date))) {
				switch (moods.get(getDateKey(date))?.score) {
					case 1:
						classes += 'bg-red-300 dark:bg-red-900 ';
						break;
					case 2:
						classes += 'bg-orange-300 dark:bg-orange-900 ';
						break;
					case 3:
						classes += 'bg-yellow-400 dark:bg-yellow-900 ';
						break;
					case 4:
						classes += 'bg-green-300 dark:bg-green-900 ';
						break;
					case 5:
						classes += 'bg-blue-300 dark:bg-blue-900 ';
						break;
				}
			} else {
				classes += 'bg-stone-100 dark:bg-stone-800 ';
			}
		}
		console.log(classes);
		return classes;
	}

	// Event handlers for "previous" and "next" buttons
	function goToPreviousMonth() {
		if (currentMonth === 0) {
			currentMonth = 11; // Wrap to December
			currentYear -= 1; // Move to the previous year
		} else {
			currentMonth -= 1; // Decrement the month
		}
	}

	function goToNextMonth() {
		if (currentMonth === 11) {
			currentMonth = 0; // Wrap to January
			currentYear += 1; // Move to the next year
		} else {
			currentMonth += 1; // Increment the month
		}
	}
</script>

<svelte:head>
	<title>Sentimenta | Track moods</title>
</svelte:head>

<div class="mx-auto my-4 md:max-w-3/5">
	<div class="flex items-center justify-between p-4">
		<button class="text-gray-300 hover:text-yellow-300" onclick={goToPreviousMonth}
			>{m.previous()}</button
		>
		<h1 class="my-4 text-center text-2xl font-bold">{currentYear}/{currentMonth + 1}</h1>
		<button class="text-gray-300 hover:text-yellow-300" onclick={goToNextMonth}>{m.next()}</button>
	</div>

	<div
		class="grid grid-cols-7 gap-0.5 border border-stone-300 bg-white p-1 text-sm text-black md:p-2 dark:border-white/10 dark:bg-stone-900 dark:text-white"
	>
		{#each [m.sun(), m.mon(), m.tue(), m.wed(), m.thu(), m.fri(), m.sat()] as day (day)}
			<div class="border-b border-white/10 p-2 text-center font-bold">{day}</div>
		{/each}

		{#each days as date (date)}
			<button
				class={`${getDayClass(date, moodMap)} flex aspect-square items-center justify-center p-2 text-center text-xl`}
				class:border={date instanceof Date && !isNaN(date.getDate())}
				onclick={() => {
					if (date instanceof Date && !isNaN(date.getDate())) {
						selectedDate = date;
						showModal = true;
					}
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
				// TEMP: Add 1 day to match backend timezone handling
				// TODO: Fix backend to handle timezones
				let nextDay = new Date(selectedDate);
				nextDay.setDate(nextDay.getDate() + 1);

				if (moodMap.has(getDateKey(selectedDate))) {
					submitInProcess = true;
					if (moodMap.get(getDateKey(selectedDate))?.uid == undefined) {
						await updateMoods();
					}
					let result = await fetch('/api/moods/update', {
						method: 'PUT',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({
							uid: moodMap.get(getDateKey(selectedDate))?.uid,
							score: mood,
							description: diary,
							date: nextDay,
							emotions: parseEmotions(emotions)
						})
					});
					if (!result.ok) {
						const errorData = result;
						console.error('Error:', errorData);
						formError = (await result.text()) || m.error_occured();
					} else {
						moods.map((m) => {
							if (m.uid == moodMap.get(getDateKey(selectedDate))?.uid) {
								m.score = mood;
								m.description = diary;
								m.emotions = parseEmotions(emotions);
							}
						});

						formError = null;
						formSuccess = true;
						showModal = false;
					}
				} else {
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
							uid: undefined,
							date: selectedDate,
							score: mood,
							description: diary,
							emotions: parseEmotions(emotions)
						};

						moods = [...moods, newMood];

						formError = null;
						formSuccess = true;
						showModal = false;
					}
				}
				submitInProcess = false;
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
					aria-label={`select mood ${n}`}
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
			aria-describedby="emotions-description"
			class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		/>

		<!-- Diary -->
		<textarea
			placeholder={m.start_today()}
			bind:value={diary}
			aria-describedby="diary-description"
			class="h-40 w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		></textarea>

		<!-- Submit -->
		<button
			type="submit"
			class="border border-current px-6 py-2 text-center transition-colors duration-300"
			class:bg-white={canSubmit()}
			class:text-black={canSubmit()}
			class:opacity-50={!canSubmit()}
			class:pointer-events-none={!canSubmit()}
			class:cursor-pointer={submitInProcess || canSubmit()}
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
