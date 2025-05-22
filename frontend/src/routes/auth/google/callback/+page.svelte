<script lang="ts">
	import { goto } from '$app/navigation';
	import { refreshUserId } from '$lib/user';
	import { onMount } from 'svelte';

	onMount(async () => {
		// Extract query params from URL
		const url = new URL(window.location.href);
		const code = url.searchParams.get('code');
		const state = url.searchParams.get('state');

		// Validate state to prevent CSRF
		if (state !== sessionStorage.getItem('google_state')) {
			alert('Invalid state');
			window.location.href = '/login';
			return;
		}

		const codeVerifier = sessionStorage.getItem('google_code_verifier');

		if (!code || !codeVerifier) {
			window.location.href = '/login';
			return;
		}

		// Send code and code_verifier to backend
		try {
			const response = await fetch('/api/auth/google/callback', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					code,
					codeVerifier,
					timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
				}),
				credentials: 'include' // Required for cookies
			});

			if (response.ok) {
				// Clear PKCE and state
				sessionStorage.removeItem('google_code_verifier');
				sessionStorage.removeItem('google_state');

				// Redirect user to dashboard
				refreshUserId();
				localStorage.setItem('justRegistered', 'true');
				goto('/track');
			} else {
				console.error('Google auth failed');
				window.location.href = '/login';
			}
		} catch (error) {
			console.error('Error during Google auth:', error);
			window.location.href = '/login';
		}
	});
</script>

<!-- Optional: Show loading indicator -->
<div class="p-4 text-center">
	<p>Logging in with Google...</p>
</div>
