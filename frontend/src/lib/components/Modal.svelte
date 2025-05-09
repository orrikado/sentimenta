<script lang="ts">
	import { m } from '$lib/paraglide/messages';

	let { showModal = $bindable(), header, children, sideContent, ...rest } = $props();

	// State with proper null handling
	let modalContent = $state<HTMLElement | null>(null);
	let sideContentEl = $state<HTMLElement | null>(null);
	let firstFocusableElement = $state<HTMLElement | null>(null);

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

	function repositionSideContent() {
		if (!modalContent || !sideContentEl) return;

		const rect = modalContent.getBoundingClientRect();
		const top = rect.bottom + 16; // 16px gap below modal

		sideContentEl.style.position = 'absolute';
		sideContentEl.style.top = `${top + 6}px`;
		if (window.innerWidth > 640) {
			sideContentEl.style.left = '50%';
			sideContentEl.style.transform = 'translateX(-50%)';
		} else {
			sideContentEl.style.left = '5%';
			sideContentEl.style.transform = 'translateX(-2.5%)';
		}
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

			// Position side content
			repositionSideContent();
		}
	});

	// Update position on resize
	$effect(() => {
		if (!showModal) return;

		const handleResize = () => repositionSideContent();
		window.addEventListener('resize', handleResize);
		return () => window.removeEventListener('resize', handleResize);
	});
</script>

{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/30"
		onclick={() => (showModal = false)}
		role="dialog"
		tabindex="0"
		aria-modal="true"
		aria-labelledby="modal-title"
		onkeydown={handleEscape}
	>
		<dialog
			open
			{...rest}
			class="mx-auto w-full max-w-md border border-white/10 bg-white p-4 font-mono text-sm text-black dark:bg-stone-900 dark:text-white"
			onclick={(event) => event.stopPropagation()}
			onkeydown={trapFocus}
		>
			<div bind:this={modalContent}>
				<div id="modal-title" class="flex items-center justify-between">
					{@render header?.()}
					<button
						aria-label={m.close_modal()}
						onclick={() => (showModal = false)}
						class="text-xl hover:text-yellow-500"
						tabindex="0"
					>
						×
					</button>
				</div>
				<hr class="my-2 border-white/10" />
				{@render children?.()}
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
			</div>
		</dialog>

		<div
			bind:this={sideContentEl}
			class="z-50 border border-white/20 bg-stone-200 pt-3 pb-4 text-black dark:bg-stone-800 dark:text-white"
			class:hidden={!sideContentEl?.innerText}
		>
			{@render sideContent()}
		</div>
	</div>
{/if}
