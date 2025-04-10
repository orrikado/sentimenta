<script lang="ts">
	import { onMount } from 'svelte';
	import { userId } from '$lib/stores/user';

	let showLogin = true;

	// Check for the cookie on component mount
	onMount(() => {
		if ($userId) {
			showLogin = false;
		}
	});
</script>

<nav class="flex items-center justify-between bg-yellow-800 p-4 text-white">
	<a href="/" class="text-xl font-bold">Logo</a>
	<ul class="flex space-x-4">
		<li><a href="/track" class="hover:text-gray-400">Track</a></li>
		{#if showLogin}
			<li><b><a href="/login" class="hover:text-gray-400">Login</a></b></li>
			<li><b><a href="/signup" class="hover:text-gray-400">Signup</a></b></li>
		{:else}
			<li><b><a href="/profile" class="hover:text-gray-400">Profile</a></b></li>
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
