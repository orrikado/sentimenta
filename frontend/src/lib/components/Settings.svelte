<script lang="ts">
	import { browser } from '$app/environment';
	import { env } from '$env/dynamic/public';
	import { m } from '$lib/paraglide/messages';
	import { user } from '$lib/stores/user';
	import { refreshUser } from '$lib/user';

	// Local state for first day of week
	let selectedFirstDay = $state(
		browser ? parseInt(localStorage.getItem('firstDayOfWeek') || '1') : 1
	);

	// Save to localStorage when changed
	$effect(() => {
		localStorage.setItem('firstDayOfWeek', selectedFirstDay.toString());
	});

	// Derived AI toggle state
	let useAi = $derived($user?.use_ai ?? true);
	let error = $state<string | null>(null);

	// Handle AI toggle changes
	async function handleAiToggle(e: Event) {
		const newValue = (e.target as HTMLInputElement).checked;
		const oldValue = useAi;
		useAi = newValue;

		try {
			const response = await fetch('/api/user/update', {
				method: 'PATCH',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ use_ai: newValue })
			});

			if (!response.ok) {
				throw new Error('Server error');
			}
			await refreshUser();
		} catch (err) {
			// Revert on error
			useAi = oldValue;
			console.error(err);
			error = m.ai_update_failed();
		}
	}
</script>

<!-- First Day of Week Setting -->
<div class="mb-6">
	<label for="first-day-select" class="text-stone-500 uppercase dark:text-stone-400">
		{m.first_day_of_week()}
	</label>
	<select
		id="first-day-select"
		bind:value={selectedFirstDay}
		class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
	>
		<option value={0}>{m.sunday()}</option>
		<option value={1}>{m.monday()}</option>
		<option value={2}>{m.tuesday()}</option>
		<option value={3}>{m.wednesday()}</option>
		<option value={4}>{m.thursday()}</option>
		<option value={5}>{m.friday()}</option>
		<option value={6}>{m.saturday()}</option>
	</select>
</div>

<!-- AI Toggle -->
<div class="flex items-center justify-between">
	<span class="text-stone-500 uppercase dark:text-stone-400">
		{m.use_ai()}
	</span>
	<input
		type="checkbox"
		checked={useAi}
		onchange={handleAiToggle}
		disabled={env.PUBLIC_AI_ENABLED !== 'true'}
		class="h-5 w-9 rounded-full border-stone-300 bg-stone-200 transition-colors duration-200 focus:ring-stone-500 disabled:brightness-50 dark:bg-stone-700"
	/>
</div>

{#if env.PUBLIC_AI_ENABLED !== 'true'}
	<span class="text-sm text-red-500 dark:text-red-400">
		{m.ai_disabled()}
	</span>
{/if}

<!-- Error Display -->
{#if error}
	<div class="mt-2 bg-red-100 p-2 text-red-700 dark:bg-red-900 dark:text-red-100">
		{error}
	</div>
{/if}
