<script lang="ts">
	import { onMount } from 'svelte';
	function getCookie(name: string) {
		const cookies = document.cookie.split('; ');
		for (let cookie of cookies) {
			const [key, value] = cookie.split('=');
			if (key === name) {
				return decodeURIComponent(value);
			}
		}
		return null; // Return null if the cookie is not found
	}

	let showLogin = true;

	// Check for the cookie on component mount
	onMount(() => {
		const loggedin = getCookie('access_token');
		if (loggedin) {
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
