<script>
	import { m } from '$lib/paraglide/messages';

	let { showModal = $bindable(), header, children } = $props();
	// let dialog = $state(); // HTMLDialogElement

	// $effect(() => {
	// 	if (showModal) dialog.showModal();
	// });
</script>

{#if showModal}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/30"
		onclick={() => (showModal = false)}
		role="presentation"
		aria-hidden="true"
		id="backdrop"
	>
		<dialog
			open
			class="mx-auto w-full max-w-md border border-white/10 bg-white p-4 font-mono text-sm text-black dark:bg-stone-900 dark:text-white"
			onclick={(event) => event.stopPropagation()}
		>
			<div>
				{@render header?.()}
				<hr class="my-2 border-white/10" />
				{@render children?.()}
				<hr class="my-2 border-white/10" />
				<div class="mt-4 text-right">
					<button
						onclick={() => (showModal = false)}
						class="border border-white px-3 py-1 hover:bg-white hover:text-black dark:hover:bg-white dark:hover:text-black"
					>
						{m.close()}
					</button>
				</div>
			</div>
		</dialog>
	</div>
{/if}
