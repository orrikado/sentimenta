<!-- src/lib/components/RegistrationModal.svelte -->
<script lang="ts">
	import Settings from './Settings.svelte';

	// Make showModal bindable
	let { showModal = $bindable(), onClose } = $props();
</script>

{#if showModal}
	<div class="bg-opacity-50 fixed inset-0 z-50 flex items-center justify-center bg-black">
		<div class="w-full max-w-md bg-white p-6 dark:bg-stone-900">
			<h2 class="mb-4 text-xl font-bold">Complete Your Setup</h2>

			<!-- Settings in Modal -->
			<Settings />

			<div class="mt-6 flex justify-end">
				<button
					class="border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
					onclick={() => {
						showModal = false;
						onClose?.();
					}}
				>
					Finish
				</button>
			</div>
		</div>
	</div>
{/if}
