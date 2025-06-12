<script lang="ts">
	// External imports
	import { onDestroy, onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { getMonthDays } from '$lib/calendar-utils';
	import { m } from '$lib/paraglide/messages';
	import { user, userId } from '$lib/stores/user';
	import Modal from '$lib/components/Modal.svelte';
	import * as d3 from 'd3';
	import { browser } from '$app/environment';
	import { moods } from '$lib/stores/moods';
	import { updateMoods, type MoodEntry } from '$lib/moods';
	import { advice } from '$lib/stores/advice';
	import { updateAdvice } from '$lib/advice';
	import RegistrationModal from '$lib/components/RegistrationModal.svelte';
	import { refreshUser } from '$lib/user';
	import { server_status } from '$lib/stores/server_status';
	import { refreshServerStatus } from '$lib/status';
	import { env } from '$env/dynamic/public';

	// Regular variables
	let today = new Date();
	const fullEmotionPool = [
		m.emotion_joy(),
		m.emotion_sadness(),
		m.emotion_anger(),
		m.emotion_anxiety(),
		m.emotion_stress(),
		m.emotion_gratitude(),
		m.emotion_hope(),
		m.emotion_fatigue(),
		m.emotion_excitement(),
		m.emotion_calm(),
		m.emotion_loneliness(),
		m.emotion_confusion(),
		m.emotion_pride(),
		m.emotion_disappointment(),
		m.emotion_motivation(),
		m.emotion_nostalgia(),
		m.emotion_boredom(),
		m.emotion_relief(),
		m.emotion_frustration(),
		m.emotion_optimism(),
		m.emotion_guilt(),
		m.emotion_surprise(),
		m.emotion_love(),
		m.emotion_ambition(),
		m.emotion_curiosity(),
		m.emotion_apathy(),
		m.emotion_empathy(),
		m.emotion_envy(),
		m.emotion_contentment()
	];
	let socket: WebSocket;

	// State variables
	let currentMonth = $state(today.getMonth());
	let currentYear = $state(today.getFullYear());
	let animate_width = $state(false);
	let firstDayOfWeek = $state(1);
	if (browser) {
		firstDayOfWeek = parseInt(localStorage.getItem('firstDayOfWeek') || '1');
	}
	let dayHeaders = $derived(() => {
		const allDays = [m.sun(), m.mon(), m.tue(), m.wed(), m.thu(), m.fri(), m.sat()];
		return [...allDays.slice(firstDayOfWeek), ...allDays.slice(0, firstDayOfWeek)];
	});

	let showModal = $state(false);
	let selectedDate = $state(new Date());
	let submitInProcess = $state(false);
	let loading = $state(true);
	let is_put = $state(false);
	let mood = $state<number>(0);
	let emotions = $state('');
	let emotionSubset = $state([
		m.emotion_joy(),
		m.emotion_sadness(),
		m.emotion_fatigue(),
		m.emotion_boredom(),
		m.emotion_calm(),
		m.emotion_motivation()
	]);
	let diary = $state('');
	let formError = $state<string | null>(null);
	let formSuccess = $state<boolean>(false);
	let notificationMessage = $state('');
	let dimensions = $state({
		width: 600,
		height: 200,
		margin: { top: 20, right: 30, bottom: 40, left: 50 }
	});
	let showRegistationModal = $state(false);

	// Derived values
	let days: Date[] = $derived(getMonthDays(currentYear, currentMonth, firstDayOfWeek));
	let moodMap = $derived(new Map($moods.map((m) => [getDateKey(m.date), m])));
	let adviceMap = $derived(new Map($advice.map((a) => [getDateKey(a.date), a])));
	const canSubmit = $derived(() => {
		const future = new Date(selectedDate) > today;
		if (emotions.length > parseInt(env.PUBLIC_MOOD_EMOTES_LENGTH_MAX || '120')) {
			return false;
		}
		if (diary.length > parseInt(env.PUBLIC_MOOD_DESC_LENGTH_MAX || '320')) {
			return false;
		}
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
		if (emotions.length > parseInt(env.PUBLIC_MOOD_EMOTES_LENGTH_MAX || '120')) {
			formError = m.emotions_too_long();
		} else if (diary.length > parseInt(env.PUBLIC_MOOD_DESC_LENGTH_MAX || '320')) {
			formError = m.diary_too_long();
		} else {
			formError = null;
		}
	});
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
				const saved = localStorage.getItem('formData');
				if (saved) {
					try {
						const parsed = JSON.parse(saved);
						mood = parsed.mood ?? 0;
						emotions = parsed.emotions || '';
						diary = parsed.diary || '';
					} catch (e) {
						console.error('Failed to parse saved data', e);
					}
				} else {
					mood = 0;
					emotions = '';
					diary = '';
				}
			}
		}
	});
	$effect(() => {
		console.log(moodMap);
		const svg = d3.select('svg');
		svg.selectAll('g').remove();
		svg.selectAll('.axis').remove();
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
		const xScale = d3
			.scaleTime()
			.domain(d3.extent(filteredMoods, (d) => d.date) as [Date, Date])
			.range([10, innerWidth - 10]);
		const yScale = d3
			.scaleLinear()
			.domain([1, 5])
			.range([innerHeight - 10, 10]);
		const line = d3
			.line<MoodEntry>()
			.x((d) => xScale(d.date))
			.y((d) => yScale(d.score));
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
		g.append('g')
			.attr('transform', `translate(0,${innerHeight})`)
			.attr('class', 'axis')
			.call(xAxis)
			.style('opacity', 0)
			.transition()
			.duration(200)
			.style('opacity', 1)
			.attr('pointer-events', 'none');
		g.append('g')
			.attr('class', 'axis')
			.call(yAxis)
			.style('opacity', 0)
			.transition()
			.duration(200)
			.style('opacity', 1)
			.attr('pointer-events', 'none');
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
			.attr('fill', (d) => getDotColor(d.score))
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
		circles
			.transition()
			.duration(200)
			.attr('cx', (d) => xScale(d.date))
			.attr('cy', (d) => yScale(d.score));
	});

	// Functions
	function refreshEmotions() {
		const selected = getEmotionsArray();
		const available = fullEmotionPool.filter(
			(e) => !selected.some((sel) => sel.toLowerCase() === e.toLowerCase())
		);
		const shuffled = [...available].sort(() => 0.5 - Math.random());
		emotionSubset = shuffled.slice(0, 7);
	}
	function getEmotionsArray(): string[] {
		return emotions
			.split(',')
			.map((e) => e.trim())
			.filter((e) => e !== '');
	}
	function isSelected(emotion: string): boolean {
		const lowerEmotions = getEmotionsArray().map((e) => e.toLowerCase());
		return lowerEmotions.includes(emotion.toLowerCase());
	}
	function addEmotion(emotion: string): void {
		const current = getEmotionsArray();
		const lowerCurrent = current.map((e) => e.toLowerCase());
		if (!lowerCurrent.includes(emotion.toLowerCase())) {
			emotions = [...current, emotion].join(',');
		}
	}
	function getDateKey(dateInput: string | number | Date) {
		const date = dateInput instanceof Date ? dateInput : new Date(dateInput);
		return `${date.getFullYear()}-${date.getMonth() + 1}-${date.getDate()}`;
	}
	function isToday(date: {
		getDate: () => number;
		getMonth: () => number;
		getFullYear: () => number;
	}) {
		const now = new Date();
		return (
			date.getDate() === now.getDate() &&
			date.getMonth() === now.getMonth() &&
			date.getFullYear() === now.getFullYear()
		);
	}
	function isYesterday(date: {
		getDate: () => number;
		getMonth: () => number;
		getFullYear: () => number;
	}) {
		if (!(date instanceof Date)) {
			throw new Error('Invalid argument: you must provide a "date" instance');
		}
		const yesterday = new Date();
		yesterday.setDate(yesterday.getDate() - 1);
		return (
			date.getDate() === yesterday.getDate() &&
			date.getMonth() === yesterday.getMonth() &&
			date.getFullYear() === yesterday.getFullYear()
		);
	}
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
		if (d.description && d.description.trim() !== '') {
			tooltip.select('.description').text(d.description).style('display', 'block');
			tooltip.select('.my-2').style('display', 'block');
			tooltip.select('.description-box').style('display', 'block');
		} else {
			tooltip.select('.description').style('display', 'none');
			tooltip.select('.my-2').style('display', 'none');
			tooltip.select('.description-box').style('display', 'none');
		}
		tooltip
			.classed('hidden', false)
			.style('left', `${event.pageX + 10}px`)
			.style('top', `${event.pageY - 28}px`)
			.style('opacity', 0);
		setTimeout(() => {
			tooltip.style('opacity', 1);
		}, 10);
	}
	function hideTooltip() {
		const tooltip = d3.select('#tooltip');
		tooltip.select('.emotions').html('');
		tooltip.style('opacity', 0);
		setTimeout(() => {
			tooltip.classed('hidden', true);
		}, 200);
	}
	function updateDimensions() {
		const svgContainer = document.getElementById('mood-chart')?.parentElement;
		if (svgContainer) {
			const containerWidth = svgContainer.clientWidth;
			if (containerWidth < 640) {
				dimensions.margin = { top: 20 / 3, right: 30 / 3, bottom: 40 / 3, left: 50 / 3 };
			} else {
				dimensions.margin = { top: 20, right: 30, bottom: 40, left: 50 };
			}
			dimensions.width = Math.min(containerWidth * 0.9, 800);
			dimensions.height = Math.min(containerWidth * 0.4, 250);
		}
	}
	function getDotColor(score: number) {
		return getComputedStyle(document.documentElement)
			.getPropertyValue(`--color-mood-${score}`)
			.trim();
	}
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
	function goToPreviousMonth() {
		if (currentMonth === 0) {
			currentMonth = 11;
			currentYear -= 1;
		} else {
			currentMonth -= 1;
		}
	}
	function goToNextMonth() {
		if (currentMonth === 11) {
			currentMonth = 0;
			currentYear += 1;
		} else {
			currentMonth += 1;
		}
	}

	// Lifecycle
	onMount(async () => {
		if (!browser) return;
		if (typeof window === 'undefined') return;
		await refreshServerStatus();
		if (!$userId && $server_status) {
			goto('/login');
			return;
		}
		if (localStorage.getItem('justRegistered') === 'true') {
			localStorage.removeItem('justRegistered');
			showRegistationModal = true;
		}
		refreshUser();
		updateDimensions();
		window.addEventListener('resize', updateDimensions);
		if ($moods.length === 0) {
			await updateMoods();
		}
		if ($advice.length === 0) {
			await updateAdvice();
		}
		loading = false;
	});
	onDestroy(() => {
		if (!browser) return;
		window.removeEventListener('resize', updateDimensions);
	});
</script>

<svelte:head>
	<title>Sentimenta | Track moods</title>
</svelte:head>

<RegistrationModal
	bind:showModal={showRegistationModal}
	onClose={() => (firstDayOfWeek = parseInt(localStorage.getItem('firstDayOfWeek') || '1'))}
/>

<main class="mx-auto my-4 md:max-w-3/5 xl:max-w-1/2">
	<div
		class="flex justify-center border border-stone-300 bg-stone-100 dark:border-white/10 dark:bg-stone-900"
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
		<div class="description-box flex flex-col gap-1">
			<span class="text-xs font-medium text-gray-600 dark:text-gray-300"
				>{m.tooltip_diary_label()}</span
			>
			<span
				class="description max-h-24 overflow-y-auto border border-black/10 bg-stone-50 px-2 py-1 text-sm dark:border-white/10 dark:bg-stone-800"
			></span>
		</div>
	</div>
	<div class="flex items-center justify-between px-4 sm:py-4">
		<button class="  hover:text-accent text-gray-900 dark:text-gray-300" onclick={goToPreviousMonth}
			>{m.previous()}</button
		>
		<h1 class="my-4 text-center text-2xl font-bold">{currentYear}/{currentMonth + 1}</h1>
		<button class="  hover:text-accent text-gray-900 dark:text-gray-300" onclick={goToNextMonth}
			>{m.next()}</button
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

<Modal bind:showModal {animate_width}>
	{#snippet header()}
		<h2 id="modal-title" class="text-center font-bold">
			{m.log_for()}
			{selectedDate.toLocaleDateString()}
		</h2>
	{/snippet}

	{#snippet sideContent()}
		{#if adviceMap.has(getDateKey(selectedDate))}
			<div class="justify-center text-left">
				<h3 class="text-lg font-bold">{m.advice()}</h3>
				<p class="text-lg"><i>{adviceMap.get(getDateKey(selectedDate))?.text}</i></p>
			</div>
		{/if}
	{/snippet}

	<form
		class="flex w-full max-w-md flex-col gap-6"
		onsubmit={async () => {
			if (!canSubmit() || submitInProcess) return;

			localStorage.removeItem('formData');

			submitInProcess = true; // Show spinner right away

			let nextDay = new Date(selectedDate);
			nextDay.setDate(nextDay.getDate() + 1); // workaround because js is stupid

			if (moodMap.has(getDateKey(selectedDate))) {
				if (moodMap.get(getDateKey(selectedDate))?.uid == undefined) {
					await updateMoods();
				}
				is_put = true;
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
					$moods = $moods.map((m) => {
						if (m.uid == moodMap.get(getDateKey(selectedDate))?.uid) {
							return {
								...m,
								score: mood,
								description: diary,
								emotions: parseEmotions(emotions)
							};
						}
						return m;
					});

					formError = null;
					formSuccess = true;
					showModal = false;
					updateDimensions();
				}
			} else {
				is_put = false;
				if ($user?.use_ai === true && (isToday(selectedDate) || isYesterday(selectedDate))) {
					const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
					socket = new WebSocket(wsProtocol + window.location.host + '/ws');

					socket.addEventListener('message', (event) => {
						console.log('Received:', event.data);
						let newAdvice = JSON.parse(event.data);
						newAdvice.date = new Date(newAdvice.date).getTime() - 1 * 24 * 60 * 60 * 1000;
						newAdvice.generated_by_websocket = true;
						advice.set([...$advice, newAdvice]);

						// Force modal to re-render if open
						if (showModal && getDateKey(selectedDate) === getDateKey(newAdvice.date)) {
							animate_width = true;
							showModal = false;
							setTimeout(() => (showModal = true), 0);
							setTimeout(() => (animate_width = false), 350);
						}
						submitInProcess = false;
						socket.close();
					});
				}

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
				}
			}
			if (
				$user?.use_ai !== true ||
				is_put ||
				!(isToday(selectedDate) || isYesterday(selectedDate))
			) {
				submitInProcess = false;
			}
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
		<div>
			<label for="emotions" class="sr-only">{m.emotions_label()}</label>
			<input
				type="text"
				id="emotions"
				placeholder={m.start_emotions_placeholder()}
				bind:value={emotions}
				class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
			/>

			<!-- Emotion Buttons -->
			<div class="mt-2 flex flex-wrap gap-1">
				{#each emotionSubset as emotion (emotion)}
					{#if !isSelected(emotion)}
						<button
							type="button"
							onclick={() => addEmotion(emotion)}
							class=" cursor-pointer border bg-stone-100 px-2 py-1 hover:brightness-90 dark:bg-stone-800 dark:text-white"
						>
							+ {emotion}
						</button>
					{/if}
				{/each}

				<button
					type="button"
					onclick={refreshEmotions}
					class="border p-1 text-gray-700 hover:text-gray-900 dark:text-gray-300 dark:hover:text-white"
					aria-label="Refresh emotions"
				>
					<svg
						class="h-6 w-6"
						viewBox="0 0 489.533 489.533"
						fill="currentColor"
						stroke="currentColor"
						stroke-width="2.5"
					>
						<path
							d="M268.175,488.161c98.2-11,176.9-89.5,188.1-187.7c14.7-128.4-85.1-237.7-210.2-239.1v-57.6c0-3.2-4-4.9-6.7-2.9l-118.6,87.1c-2,1.5-2,4.4,0,5.9l118.6,87.1c2.7,2,6.7,0.2,6.7-2.9v-57.5c87.9,1.4,158.3,76.2,152.3,165.6c-5.1,76.9-67.8,139.3-144.7,144.2c-81.5,5.2-150.8-53-163.2-130c-2.3-14.3-14.8-24.7-29.2-24.7c-17.9,0-31.9,15.9-29.1,33.6C49.575,418.961,150.875,501.261,268.175,488.161z"
						/>
					</svg>
				</button>
			</div>
		</div>

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
			<span>{m.button_save()}</span>
		</button>

		<div aria-live="polite" aria-atomic="true">
			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
			{#if formSuccess}
				{#if $user?.use_ai === true && !is_put && (isToday(selectedDate) || isYesterday(selectedDate))}
					<p class="text-sm text-green-500 dark:text-green-400">{m.advice_generating()}</p>
				{:else}
					<p class="text-sm text-green-500 dark:text-green-400">{m.mood_upload_success()}</p>
				{/if}
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
