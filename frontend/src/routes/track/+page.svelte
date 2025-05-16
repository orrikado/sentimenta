<script lang="ts">
	// External imports
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getMonthDays } from '$lib/calendar-utils';
	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import Modal from '$lib/components/Modal.svelte';
	import * as d3 from 'd3';
	import { browser } from '$app/environment';
	import { moods } from '$lib/stores/moods';
	import { updateMoods, type MoodEntry } from '$lib/moods';
	import { advice } from '$lib/stores/advice';
	import { updateAdvice } from '$lib/advice';

	// State variables
	let today = new Date();
	let currentMonth = $state(today.getMonth());
	let currentYear = $state(today.getFullYear());

	// State for first day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	let firstDayOfWeek = $state(1);
	if (browser) {
		firstDayOfWeek = parseInt(localStorage.getItem('firstDayOfWeek') || '1');
	}

	let dayHeaders = $state<string[]>([]);

	// Rebuild headers whenever firstDayOfWeek or locale messages change
	$effect(() => {
		const allDays = [m.sun(), m.mon(), m.tue(), m.wed(), m.thu(), m.fri(), m.sat()];
		dayHeaders = [...allDays.slice(firstDayOfWeek), ...allDays.slice(0, firstDayOfWeek)];
	});

	let days: Date[] = $derived(getMonthDays(currentYear, currentMonth, firstDayOfWeek));
	let showModal = $state(false);
	let selectedDate: Date = $state(new Date());
	let submitInProcess = $state(false);
	let loading = $state(true);

	let mood = $state<number>(0);
	let emotions = $state('');
	let diary = $state('');
	let formError = $state<string | null>(null);
	let formSuccess = $state<boolean>(false);
	let notificationMessage = $state('');

	// Derived values
	const getDateKey = (dateInput: string | number | Date) => {
		const date = dateInput instanceof Date ? dateInput : new Date(dateInput);
		return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
	};

	let moodMap = $derived(new Map($moods.map((m) => [getDateKey(m.date), m])));
	let adviceMap = $derived(new Map($advice.map((a) => [getDateKey(a.date), a])));
	const canSubmit = $derived(() => {
		const future = new Date(selectedDate) > today;
		// Ensure the selected date is today or in the past
		return !!(mood !== 0 && emotions.trim() && !future);
	});

	let filteredMoods = $derived(
		$moods
			.filter((m) => {
				return m.date.getMonth() === currentMonth && m.date.getFullYear() === currentYear;
			})
			.sort((a, b) => a.date.getTime() - b.date.getTime())
	);

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

	let dimensions = $state({
		width: 600,
		height: 200,
		margin: { top: 20, right: 30, bottom: 40, left: 50 }
	});

	function showTooltip(event: { pageX: number; pageY: number }, d: MoodEntry) {
		const dateStr = d3.timeFormat('%B %d, %Y')(d.date);
		const tooltip = d3.select('#tooltip');

		tooltip.select('.date').text(dateStr);
		tooltip.select('.score').text(d.score);

		const emotionsContainer = tooltip.select('.emotions').html('');
		const emotionList =
			d.emotions
				?.split(',')
				.map((e) => e.trim())
				.filter(Boolean) || [];

		emotionsContainer
			.selectAll('span')
			.data(emotionList)
			.enter()
			.append('span')
			.attr(
				'class',
				'inline-flex items-center mx-px border border-black/10 bg-stone-50 px-1 py-px text-xs dark:border-white/10 dark:bg-stone-800'
			)
			.text((e) => e);

		tooltip.select('.description').text(d.description);

		// Set position instantly
		tooltip
			.classed('hidden', false)
			.style('left', `${event.pageX + 10}px`)
			.style('top', `${event.pageY - 28}px`)
			.style('opacity', 0); // Start invisible

		// Fade in (only opacity transition)
		setTimeout(() => {
			tooltip.style('opacity', 1);
		}, 10); // Small delay to ensure DOM updates
	}

	function hideTooltip() {
		const tooltip = d3.select('#tooltip');

		tooltip.select('.emotions').html('');

		// Fade out
		tooltip.style('opacity', 0);

		// Hide after fade
		setTimeout(() => {
			tooltip.classed('hidden', true);
		}, 200); // Match duration-200
	}

	function updateDimensions() {
		const svgContainer = document.getElementById('mood-chart')?.parentElement;
		if (svgContainer) {
			const containerWidth = svgContainer.clientWidth;
			dimensions.width = Math.min(containerWidth * 0.9, 800); // Max width 800
			dimensions.height = Math.min(containerWidth * 0.4, 250); // 1/3 of width
		}
	}

	// Lifecycle
	onMount(async () => {
		if (!browser) return;
		if (typeof window === 'undefined') return;
		if (!$userId) {
			goto('/login');
			return;
		}
		updateDimensions();
		window.addEventListener('resize', updateDimensions);

		if ($moods.length === 0) {
			await updateMoods();
		}
		if ($advice.length === 0) {
			await updateAdvice();
		}
		loading = false; // Stop loading after data is fetched
	});

	onDestroy(() => {
		if (!browser) return;
		window.removeEventListener('resize', updateDimensions);
	});

	function getDotColor(score: number) {
		return getComputedStyle(document.documentElement)
			.getPropertyValue(`--color-mood-${score}`)
			.trim();
	}

	// draw the chart
	$effect(() => {
		const svg = d3.select('svg');

		// Clear previous chart content
		svg.selectAll('g').remove(); // Remove all <g> elements
		svg.selectAll('.axis').remove(); // Remove old axes

		const g = svg
			.append('g')
			.attr('transform', `translate(${dimensions.margin.left}, ${dimensions.margin.top})`);

		const innerWidth = dimensions.width - dimensions.margin.left - dimensions.margin.right;
		const innerHeight = dimensions.height - dimensions.margin.top - dimensions.margin.bottom;

		if (filteredMoods.length === 0) {
			g.append('text')
				.attr('x', innerWidth / 2)
				.attr('y', innerHeight / 2)
				.attr('text-anchor', 'middle')
				.attr('fill', '#666')
				.text(m.no_data_for_month());
			return;
		}

		// mood gradient
		svg
			.append('defs')
			.append('linearGradient')
			.attr('id', 'moodGradient')
			.attr('x1', '0%')
			.attr('y1', '100%')
			.attr('x2', '0%')
			.attr('y2', '0%')
			.selectAll('stop')
			.data([
				{
					offset: '0%',
					color: getComputedStyle(document.documentElement).getPropertyValue('--color-mood-1')
				},
				{
					offset: '100%',
					color: getComputedStyle(document.documentElement).getPropertyValue('--color-mood-5')
				}
			])
			.enter()
			.append('stop')
			.attr('offset', (d) => d.offset)
			.attr('stop-color', (d) => d.color);

		// Scales
		const xScale = d3
			.scaleTime()
			.domain(d3.extent(filteredMoods, (d) => d.date) as [Date, Date])
			.range([10, innerWidth - 10]);

		const yScale = d3
			.scaleLinear()
			.domain([1, 5])
			.range([innerHeight - 10, 10]);

		// Line generator
		const line = d3
			.line<MoodEntry>()
			.x((d) => xScale(d.date))
			.y((d) => yScale(d.score));

		// Draw path
		g.append('path')
			.datum(filteredMoods)
			.attr('fill', 'none')
			.attr('stroke', 'url(#moodGradient)')
			.attr('stroke-width', 2)
			.attr('stroke-opacity', 0.8)
			.attr('d', line)
			.style('opacity', 0)
			.transition()
			.duration(200)
			.style('opacity', 1)
			.attr('pointer-events', 'none');

		// Axes
		const xAxis = d3
			.axisBottom(xScale)
			.tickValues(filteredMoods.map((d) => d.date))
			.tickFormat((dateObj) => {
				const date = dateObj as Date;
				return d3.timeFormat('%e')(date);
			})
			.tickSizeInner(-innerHeight)
			.tickSizeOuter(0);

		const yAxis = d3.axisLeft(yScale).ticks(5);

		// Draw x-axis with transition
		g.append('g')
			.attr('transform', `translate(0,${innerHeight})`)
			.attr('class', 'axis')
			.call(xAxis)
			.style('opacity', 0)
			.transition()
			.duration(200)
			.style('opacity', 1)
			.attr('pointer-events', 'none');

		// Draw y-axis with transition
		g.append('g')
			.attr('class', 'axis')
			.call(yAxis)
			.style('opacity', 0)
			.transition()
			.duration(200)
			.style('opacity', 1)
			.attr('pointer-events', 'none');

		// Add area under the line
		const area = d3
			.area<MoodEntry>()
			.x((d) => xScale(d.date))
			.y0(innerHeight)
			.y1((d) => yScale(d.score));

		g.append('path')
			.datum(filteredMoods)
			.attr('fill', 'url(#moodGradient)')
			.attr('fill-opacity', 0.16)
			.attr('d', area)
			.attr('pointer-events', 'none');

		// Draw circles with transitions
		const circles = g
			.selectAll<SVGCircleElement, MoodEntry>('circle')
			.data<MoodEntry>(filteredMoods, (d) => d.date.toISOString());

		circles.exit().transition().duration(300).style('opacity', 0).remove();

		circles
			.enter()
			.append('circle')
			.attr('cx', (d) => xScale(d.date))
			.attr('cy', (d) => yScale(d.score))
			.attr('r', 0)
			.attr('fill', (d) => getDotColor(d.score)) // Dynamic color
			.style('cursor', 'pointer')
			.on('mouseover', (event, d) => {
				showTooltip(event, d);
				d3.select(event.currentTarget).transition().duration(100).attr('r', 8);
			})
			.on('mouseout', (event) => {
				hideTooltip();
				d3.select(event.currentTarget).transition().duration(100).attr('r', 5);
			})
			.on('click', (_, d) => {
				selectedDate = d.date;
				showModal = true;
			})
			.transition()
			.duration(200)
			.attr('r', 5);

		// Update existing circles
		circles
			.transition()
			.duration(200)
			.attr('cx', (d) => xScale(d.date))
			.attr('cy', (d) => yScale(d.score));
	});

	// Functions
	function parseEmotions(input: string): string {
		return input
			.toLowerCase()
			.split(',')
			.map((phrase) => phrase.trim())
			.filter((phrase) => phrase.length > 0)
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
						classes += 'bg-mood-1 ';
						break;
					case 2:
						classes += 'bg-mood-2 ';
						break;
					case 3:
						classes += 'bg-mood-3 ';
						break;
					case 4:
						classes += 'bg-mood-4 ';
						break;
					case 5:
						classes += 'bg-mood-5 ';
						break;
				}
			} else {
				classes += 'bg-stone-100 dark:bg-stone-800 ';
			}
		}
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

<main class="mx-auto my-4 md:max-w-3/5 xl:max-w-1/2">
	<div
		class=" flex justify-center border border-stone-300 bg-stone-100 dark:border-white/10 dark:bg-stone-900"
	>
		<svg
			id="mood-chart"
			class="h-auto w-full max-w-2xl"
			viewBox={`0 0 ${dimensions.width} ${dimensions.height}`}
			preserveAspectRatio="xMidYMid meet"
		></svg>
	</div>
	<div
		id="tooltip"
		class="absolute z-10 hidden max-w-md border bg-white p-2 text-sm shadow-md dark:bg-stone-900"
	>
		<div><span class="date mb-2 text-lg font-bold dark:text-stone-300"></span></div>
		<div>
			<span class="font-medium">{m.tooltip_score_label()}</span>
			<span class="score text-xl"></span><span>/5</span>
		</div>
		<!-- Emotions -->
		<div class="flex flex-col gap-1">
			<div class="emotions flex flex-wrap gap-1">
				<!-- Chips will be inserted here -->
			</div>
		</div>
		<!-- Divider -->
		<div class="my-2 border-t border-black/10 dark:border-white/10"></div>
		<!-- Diary -->
		<div class="flex flex-col gap-1">
			<span class="text-xs font-medium text-gray-600 dark:text-gray-300"
				>{m.tooltip_diary_label()}</span
			>
			<span
				class="description max-h-24 overflow-y-auto border border-black/10 bg-stone-50 px-2 py-1 text-sm dark:border-white/10 dark:bg-stone-800"
			></span>
		</div>
	</div>
	<div class="flex items-center justify-between p-4">
		<button
			class="  text-gray-900 hover:text-yellow-600 dark:text-gray-300 dark:hover:text-yellow-300"
			onclick={goToPreviousMonth}>{m.previous()}</button
		>
		<h1 class="my-4 text-center text-2xl font-bold">{currentYear}/{currentMonth + 1}</h1>
		<button
			class="  text-gray-900 hover:text-yellow-600 dark:text-gray-300 dark:hover:text-yellow-300"
			onclick={goToNextMonth}>{m.next()}</button
		>
	</div>

	{#if notificationMessage}
		<div
			class="animate-fade-in-out fixed top-20 right-2 z-50 bg-red-500 px-4 py-2 text-white shadow-md"
		>
			{notificationMessage}
		</div>
	{/if}

	<div
		class="grid grid-cols-7 gap-0.5 border border-stone-300 bg-white p-1 text-sm text-black md:p-2 dark:border-white/10 dark:bg-stone-900 dark:text-white"
		role="grid"
		aria-label={m.aria_calendar_grid()}
	>
		{#each dayHeaders as day (day)}
			<div class="border-b border-white/10 p-2 text-center font-bold" role="columnheader">
				{day}
			</div>
		{/each}

		{#each days as date (date)}
			{#if date instanceof Date && !isNaN(date.getDate())}
				<button
					class={`${getDayClass(date, moodMap)} flex aspect-square items-center justify-center p-2 text-center text-xl transition-all duration-200 `}
					class:border={date instanceof Date && !isNaN(date.getDate())}
					class:hover:brightness-85={!(date > today)}
					tabindex="0"
					class:cursor-pointer={!(date > today)}
					class:cursor-not-allowed={date > today}
					aria-disabled={date > today}
					role="gridcell"
					aria-label={m.aria_select_date({ date: date.toLocaleDateString() })}
					aria-selected={moodMap.has(getDateKey(date)) ? 'true' : 'false'}
					onclick={() => {
						if (date instanceof Date && !isNaN(date.getDate())) {
							if (date > today) {
								notificationMessage = m.cant_submit_future();
								setTimeout(() => (notificationMessage = ''), 3000);
								return;
							} else {
								selectedDate = date;
								showModal = true;
							}
						}
					}}
				>
					{date.getDate()}
				</button>
			{:else}
				<div class="aspect-square" role="gridcell" aria-hidden="true">
					<!-- Blank cell -->
				</div>
			{/if}
		{/each}
	</div>
	{#if loading}
		<div
			class="bg-opacity-50 fixed inset-0 z-40 bg-black backdrop-blur-sm transition-opacity duration-300"
		></div>
		<div class="fixed inset-0 z-50 flex items-center justify-center">
			<svg
				class="h-12 w-12 animate-spin text-white"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
			>
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
				></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
			</svg>
		</div>
	{/if}
</main>

<Modal bind:showModal>
	{#snippet header()}
		<h2 id="modal-title" class="text-center font-bold">
			{m.log_for()}
			{selectedDate.toLocaleDateString()}
		</h2>
	{/snippet}

	{#snippet sideContent()}
		{#if adviceMap.has(getDateKey(selectedDate))}
			<div class="justify-center text-center">
				<h3 class="text-lg font-bold">{m.advice()}</h3>
				<p><i>{adviceMap.get(getDateKey(selectedDate))?.text}</i></p>
			</div>
		{/if}
	{/snippet}

	<form
		class="flex w-full max-w-md flex-col gap-6"
		onsubmit={async () => {
			if (!canSubmit() || submitInProcess) return;

			submitInProcess = true; // Show spinner right away

			let nextDay = new Date(selectedDate);
			nextDay.setDate(nextDay.getDate() + 1); // workaround because js is stupid

			if (moodMap.has(getDateKey(selectedDate))) {
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
					$moods.map((m) => {
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

					moods.set([...$moods, newMood]);

					formError = null;
					formSuccess = true;
					showModal = false;
				}
			}
			submitInProcess = false;
		}}
	>
		<h1 class="text-center text-2xl">{m.start_your_day()}</h1>

		<!-- Mood Rating -->
		<div role="radiogroup" class="flex justify-center gap-4">
			{#each [1, 2, 3, 4, 5] as n (n)}
				<button
					type="button"
					role="radio"
					aria-checked={mood === n}
					class="cursor-pointer border border-current px-3 py-2"
					class:bg-white={mood === n}
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
		<label for="emotions" class="sr-only">{m.emotions_label()}</label>
		<input
			type="text"
			id="emotions"
			placeholder={m.start_emotions_placeholder()}
			bind:value={emotions}
			class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		/>

		<!-- Diary -->
		<label for="diary" class="sr-only">{m.diary_label()}</label>
		<textarea
			placeholder={m.start_today()}
			id="diary"
			bind:value={diary}
			class="h-40 w-full border border-current bg-transparent px-3 py-2 placeholder-current"
		></textarea>

		<!-- Submit -->
		<button
			type="submit"
			class="relative flex items-center justify-center border border-current px-6 py-2 text-center transition-colors duration-300"
			class:bg-white={canSubmit()}
			class:text-black={canSubmit()}
			class:opacity-50={!canSubmit()}
			class:pointer-events-none={!canSubmit() || submitInProcess}
			aria-disabled={!canSubmit() || submitInProcess}
			disabled={submitInProcess}
		>
			{#if submitInProcess}
				<svg
					class="mr-2 h-5 w-5 animate-spin"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
				>
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"
					></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
					></path>
				</svg>
			{/if}
			<span>{m.start_process()}</span>
		</button>

		<div aria-live="polite" aria-atomic="true">
			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
			{#if formSuccess}
				<p class="text-sm text-green-500 dark:text-green-400">{m.mood_upload_success()}</p>
			{/if}
		</div>
	</form>
</Modal>

<style>
	/* Base axis styles */
	:global(.axis text) {
		font-size: 10px;
		fill: #222;
		font-family: inherit;
		-webkit-user-select: none;
		user-select: none;
	}

	/* Dark mode axis text */
	@media (prefers-color-scheme: dark) {
		:global(.axis text) {
			fill: #fff; /* Light text for dark mode */
		}
	}

	/* Responsive font sizing */
	@media (min-width: 768px) {
		:global(.axis text) {
			font-size: 12px;
		}
	}

	@media (min-width: 1536px) {
		:global(.axis text) {
			font-size: 14px;
		}
	}

	/* Optional: Improve tick spacing for larger screens */
	@media (min-width: 1024px) {
		:global(.axis line, .axis path) {
			stroke-width: 1.5px;
		}
	}
	@keyframes fade-in-out {
		0% {
			opacity: 0;
			transform: translateY(-10px);
		}
		10% {
			opacity: 1;
			transform: translateY(0);
		}
		90% {
			opacity: 1;
			transform: translateY(0);
		}
		100% {
			opacity: 0;
			transform: translateY(-10px);
		}
	}
	.animate-fade-in-out {
		animation: fade-in-out 3s ease-in-out;
	}

	:global(.axis .tick line) {
		stroke-opacity: 0.1;
		stroke-dasharray: 2 2;
	}
</style>
