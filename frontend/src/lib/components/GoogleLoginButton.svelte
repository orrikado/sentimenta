<script lang="ts">
	import { PUBLIC_GOOGLE_CLIENT_ID } from '$env/static/public';
	import { m } from '$lib/paraglide/messages';

	// Handle Google Login flow
	async function handleGoogleLogin() {
		try {
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
		} catch (error) {
			console.error('Google login error:', error);
			alert(m.google_login_failed());
		}
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

<button
	type="button"
	class="flex w-full cursor-pointer items-center justify-center gap-2 border border-black bg-stone-100 px-4 py-2 text-black transition-none hover:bg-black hover:text-white dark:border-white dark:bg-stone-900 dark:text-white dark:hover:bg-white dark:hover:text-stone-900"
	on:click={handleGoogleLogin}
>
	<svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 30 30">
		<defs>
			<style>
				.color {
					fill: currentColor;
				}
			</style>
		</defs>
		<path
			class="color"
			d="M 15.003906 3 C 8.3749062 3 3 8.373 3 15 C 3 21.627 8.3749062 27 15.003906 27 C 25.013906 27 27.269078 17.707 26.330078 13 L 25 13 L 22.732422 13 L 15 13 L 15 17 L 22.738281 17 C 21.848702 20.448251 18.725955 23 15 23 C 10.582 23 7 19.418 7 15 C 7 10.582 10.582 7 15 7 C 17.009 7 18.839141 7.74575 20.244141 8.96875 L 23.085938 6.1289062 C 20.951937 4.1849063 18.116906 3 15.003906 3 z"
		/>
	</svg>
	<span class="text-sm">{m.login_with_google()}</span>
</button>
