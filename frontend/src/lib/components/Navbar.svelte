<script lang="ts">
	import { userId } from '$lib/stores/user';
	import { setLocale } from '$lib/paraglide/runtime';
	import { dev } from '$app/environment';
	import { m } from '$lib/paraglide/messages';
	import { refreshServerStatus } from '$lib/status';
	import { onMount } from 'svelte';
	import { server_status } from '$lib/stores/server_status';

	let refreshed = $state(false);
	onMount(async () => {
		await refreshServerStatus();
		refreshed = true;
		if (navigator.languages.length > 0) {
			switch (navigator.language) {
				case 'ru-RU':
				case 'ru-UA':
				case 'ru-KZ':
				case 'ru-BY':
				case 'ru':
					setLocale('ru');
					break;
				case 'en':
					setLocale('en');
					break;
			}
		}
	});
</script>

<nav
	class="flex items-center justify-between bg-stone-200 p-4 text-stone-900 dark:bg-stone-900 dark:text-white"
>
	<a href="/" class="flex items-center space-x-2">
		<img src="/favicon.svg" width="96" height="96" alt="Logo" class="h-8 w-auto" />
		<span class="text-accent text-xl font-bold">Sentimenta</span>
	</a>
	{#if dev}
		<div>
			<button onclick={() => setLocale('en')} class="hover:text-accent text-gray-300">en</button>
			<button onclick={() => setLocale('ru')} class="hover:text-accent text-gray-300">ru</button>
		</div>
	{/if}
	{#if !$server_status && refreshed}
		<div class="text-red-500 dark:text-red-400">
			<b>{m.server_status_error()}</b>
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
		z-index: 20;
	}
</style>
