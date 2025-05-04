<script lang="ts">
	import { goto } from '$app/navigation';
	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import { refreshUserId } from '$lib/user';
	import { onMount } from 'svelte';

	import { PUBLIC_GOOGLE_CLIENT_ID } from '$env/static/public';

	let formError: string | null = null;
	const email_regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

	function hasSpace(str: string) {
		return str.indexOf(' ') >= 0;
	}

	onMount(() => {
		if ($userId) {
			goto('/profile');
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();
		formError = null;

		const form = event.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());
		const { email, password } = formData;

		if (!email || !password) {
			formError = m.field_error();
			return;
		}
		if (hasSpace(password.toString())) {
			formError = m.password_spaces_error();
			return;
		}
		if (!email_regex.test(email.toString())) {
			formError = m.invalid_email_error();
			return;
		}

		try {
			const response = await fetch('/api/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password })
			});

			if (!response.ok) {
				const errorData = await response.json();
				formError = errorData.message || m.login_fail();
				console.error(errorData);
				return;
			}

			refreshUserId();
			goto('/track');
		} catch {
			formError = m.error_occured();
		}
	}

	// Google OAuth Handler
	async function handleGoogleLogin() {
		// Generate PKCE verifier and challenge
		const codeVerifier = generateCodeVerifier();
		const codeChallenge = await generateCodeChallenge(codeVerifier);
		sessionStorage.setItem('google_code_verifier', codeVerifier);

		// Generate random state
		const state = generateRandomString();
		sessionStorage.setItem('google_state', state);

		// Redirect to Google's OAuth endpoint
		const redirectUri = encodeURIComponent(window.location.origin + '/auth/google/callback');
		const scope = encodeURIComponent('email profile');

		window.location.href =
			`https://accounts.google.com/o/oauth2/v2/auth?` +
			`client_id=${PUBLIC_GOOGLE_CLIENT_ID}` +
			`&redirect_uri=${redirectUri}` +
			`&scope=${scope}` +
			`&response_type=code` +
			`&state=${state}` +
			`&code_challenge=${codeChallenge}` +
			`&code_challenge_method=S256`;
	}

	// PKCE Helpers
	function generateRandomString() {
		return window.crypto.randomUUID();
	}

	function generateCodeVerifier() {
		const array = new Uint8Array(32);
		window.crypto.getRandomValues(array);
		return btoa(String.fromCharCode(...array)).replace(/[^a-zA-Z0-9]/g, '');
	}

	async function generateCodeChallenge(codeVerifier: string) {
		const encoder = new TextEncoder();
		const digestBuffer = await crypto.subtle.digest('SHA-256', encoder.encode(codeVerifier));

		// Convert buffer to Uint8Array, then to Base64 string
		const digestArray = Array.from(new Uint8Array(digestBuffer));
		const base64Digest = btoa(String.fromCharCode(...digestArray));

		// Replace + -> -, / -> _, and remove padding (=)
		const urlSafeBase64 = base64Digest.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');

		return urlSafeBase64;
	}
</script>

<svelte:head>
	<title>Sentimenta | Login</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<main class="flex w-full max-w-md flex-col gap-6 border border-black p-6 dark:border-white">
		<h1 class="text-center text-2xl">{m.m_login()}</h1>

		<form on:submit={handleSubmit} class="flex flex-col gap-4">
			<!-- email -->
			<div class="flex flex-col gap-1">
				<label for="email" class="text-sm">email</label>
				<input
					type="email"
					id="email"
					name="email"
					required
					class="border border-black bg-stone-100 px-3 py-2 text-black placeholder-gray-500 dark:border-white dark:bg-stone-950 dark:text-white dark:placeholder-white"
					placeholder="email@proton.me"
				/>
			</div>

			<!-- password -->
			<div class="flex flex-col gap-1">
				<label for="password" class="text-sm">{m.password()}</label>
				<input
					type="password"
					id="password"
					name="password"
					required
					minlength="8"
					class="border border-black bg-stone-100 px-3 py-2 text-black placeholder-gray-500 dark:border-white dark:bg-stone-950 dark:text-white dark:placeholder-white"
					placeholder="••••••••"
				/>
			</div>

			<!-- submit -->
			<button
				type="submit"
				class="border border-black px-4 py-2 text-black transition-none hover:bg-black hover:text-white dark:border-white dark:text-white dark:hover:bg-white dark:hover:text-black"
			>
				{m.login()}
			</button>

			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
		</form>

		<button
			type="button"
			class="mt-4 border border-black px-4 py-2 text-black hover:bg-black hover:text-white dark:border-white dark:text-white dark:hover:bg-white dark:hover:text-black"
			on:click={handleGoogleLogin}
		>
			Login with Google
		</button>
	</main>
</div>
