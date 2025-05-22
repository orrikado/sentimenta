<script lang="ts">
	import { goto } from '$app/navigation';
	import { refreshUser, refreshUserId } from '$lib/user';
	import { onMount } from 'svelte';

	onMount(async () => {
		const url = new URL(window.location.href);
		const code = url.searchParams.get('code');
		const state = url.searchParams.get('state');

		// Validate state to prevent CSRF
		if (state !== sessionStorage.getItem('github_state')) {
			alert('Invalid state');
			window.location.href = '/login';
			return;
		}

		const codeVerifier = sessionStorage.getItem('github_code_verifier');

		if (!code || !codeVerifier) {
			window.location.href = '/login';
			return;
		}

		// Send code and code_verifier to backend
		try {
			const response = await fetch('/api/auth/github/callback', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ code, codeVerifier }),
				credentials: 'include'
			});

			if (response.ok) {
				sessionStorage.removeItem('github_code_verifier');
				sessionStorage.removeItem('github_state');

				refreshUserId();
				refreshUser();
				if ((await response.json()).justRegistered) {
					localStorage.setItem('justRegistered', 'true');
				}
				goto('/track');
			} else {
				console.error('GitHub auth failed');
				window.location.href = '/login';
			}
		} catch (error) {
			console.error('Error during GitHub auth:', error);
			window.location.href = '/login';
		}
	});
</script>

<!-- Optional loading message -->
<div class="p-4 text-center">
	<p>Logging in with GitHub...</p>
</div>
