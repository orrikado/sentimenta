<script lang="ts">
	import { m } from '$lib/paraglide/messages.js';
</script>

<div
	class="flex min-h-screen flex-col items-center justify-center gap-4 px-4 text-center font-mono text-white"
>
	<h1 class="text-2xl">{m.ai_judge()}</h1>
	<p class="max-w-md text-base">
		{m.frontpage_ai_notice()}
	</p>
	<a
		href="/start"
		class="border border-white px-4 py-2 transition-none hover:bg-white hover:text-black"
		>{m.ill_try()}</a
	>
</div>
