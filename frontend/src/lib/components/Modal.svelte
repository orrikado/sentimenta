<script lang="ts">
	import { m } from '$lib/paraglide/messages';

	let { showModal = $bindable(), header, children, sideContent, animate_width, ...rest } = $props();

	// State with proper null handling
	let modalContent = $state<HTMLElement | null>(null);
	let sideContentEl = $state<HTMLElement | null>(null);
	let firstFocusableElement = $state<HTMLElement | null>(null);
	let modalContainer = $state<HTMLElement | null>(null);

	$effect(() => {
		if (animate_width && modalContainer) {
			setTimeout(() => {
				modalContainer!.scrollTop = modalContainer!.scrollHeight;
			}, 50);
		}
	});

	function trapFocus(e: KeyboardEvent) {
		if (!modalContent) return;

		if (e.key !== 'Tab') return;

		const focusables = Array.from(
			modalContent.querySelectorAll<HTMLElement>(
				'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
			)
		);

		const currentIndex = focusables.findIndex((el) => el === document.activeElement);

		if (currentIndex === -1) return;

		if (e.shiftKey && currentIndex === 0) {
			e.preventDefault();
			focusables[focusables.length - 1]?.focus();
		} else if (!e.shiftKey && currentIndex === focusables.length - 1) {
			e.preventDefault();
			focusables[0]?.focus();
		}
	}

	function handleEscape(e: KeyboardEvent) {
		if (e.key === 'Escape') showModal = false;
	}

	$effect(() => {
		if (showModal && modalContent) {
			const focusables = Array.from(
				modalContent.querySelectorAll<HTMLElement>(
					'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
				)
			);

			firstFocusableElement = focusables[0] ?? null;

			setTimeout(() => {
				if (firstFocusableElement) {
					firstFocusableElement.focus();
				} else {
					modalContent?.focus();
				}
			}, 0);
		}
	});
</script>

{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center overflow-auto bg-black/30"
		onclick={() => (showModal = false)}
		role="dialog"
		tabindex="0"
		aria-modal="true"
		aria-labelledby="modal-title"
		onkeydown={handleEscape}
		bind:this={modalContainer}
	>
		<dialog
			open
			{...rest}
			class="mx-auto w-full border border-white/10 bg-white p-4 font-mono text-sm text-black dark:bg-stone-900 dark:text-white"
			style:max-width={sideContentEl?.innerText ? '48rem' : '28rem'}
			class:transition-[max-width]={animate_width}
			class:duration-300={animate_width}
			onclick={(event) => event.stopPropagation()}
			onkeydown={trapFocus}
		>
			<div id="modal-title" class="flex items-center justify-between">
				{@render header?.()}
				<button
					aria-label={m.close_modal()}
					onclick={() => (showModal = false)}
					class="hover:text-accent text-xl"
					tabindex="0"
				>
					×
				</button>
			</div>

			<hr class="my-2 border-white/10" />

			<div class="flex flex-col sm:flex-row">
				<div bind:this={modalContent} class="flex-1">
					{@render children?.()}
				</div>

				<!-- Vertical Divider -->
				{#if sideContentEl?.innerText}
					<div class="mix-h-full mx-4 border-l border-white/10 sm:block" aria-hidden="true"></div>
				{/if}

				<div bind:this={sideContentEl} class:hidden={!sideContentEl?.innerText} class="flex-1">
					{@render sideContent()}
				</div>
			</div>
			<hr class="my-2 border-white/10" />

			<div class="mt-4 text-right">
				<button
					onclick={() => (showModal = false)}
					class="border border-white px-3 py-1 hover:bg-white hover:text-black dark:hover:bg-white dark:hover:text-black"
					tabindex="0"
				>
					{m.close()}
				</button>
			</div>
		</dialog>
	</div>
{/if}
