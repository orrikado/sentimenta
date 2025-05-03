<script lang="ts">
	import { userId } from '$lib/stores/user';
	import { setLocale } from '$lib/paraglide/runtime';
	import { dev } from '$app/environment';
	import { m } from '$lib/paraglide/messages';
</script>

<nav
	class="flex items-center justify-between bg-stone-200 p-4 text-stone-900 dark:bg-stone-900 dark:text-white"
>
	<a href="/" class="flex items-center space-x-2">
		<img src="/favicon.svg" alt="Logo" class="h-8 w-auto" />
		<span class="text-xl font-bold text-yellow-700 dark:text-yellow-300">Sentimenta</span>
	</a>
	{#if dev}
		<div>
			<button onclick={() => setLocale('en')} class="text-gray-300 hover:text-yellow-300">en</button
			>
			<button onclick={() => setLocale('ru')} class="text-gray-300 hover:text-yellow-300">ru</button
			>
		</div>
	{/if}
	<ul class="flex space-x-4">
		{#if !$userId}
			<li><b><a href="/login" class="hover:text-gray-400">{m.m_login()}</a></b></li>
			<li><b><a href="/signup" class="hover:text-gray-400">{m.m_signup()}</a></b></li>
		{:else}
			<li><b><a href="/profile" class="hover:text-gray-400">{m.profile()}</a></b></li>
		{/if}
	</ul>
</nav>

<style>
	nav {
		position: sticky;
		top: 0;
		z-index: 1000;
	}
</style>
