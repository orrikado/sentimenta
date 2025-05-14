<script lang="ts">
	import { PUBLIC_GITHUB_CLIENT_ID } from '$env/static/public';
	import { generateCodeVerifier, generateCodeChallenge, generateRandomString } from '$lib/auth';
	import { m } from '$lib/paraglide/messages';

	// Handle Google Login flow
	async function handleGitHubLogin() {
		const codeVerifier = generateCodeVerifier();
		const codeChallenge = await generateCodeChallenge(codeVerifier);

		sessionStorage.setItem('github_code_verifier', codeVerifier);

		const state = generateRandomString();
		sessionStorage.setItem('github_state', state);

		const clientId = PUBLIC_GITHUB_CLIENT_ID;
		const redirectUri = encodeURIComponent(window.location.origin + '/auth/github/callback');
		const scope = encodeURIComponent('read:user user:email'); 

		window.location.href =
			`https://github.com/login/oauth/authorize?` +
			`client_id=${clientId}` +
			`&redirect_uri=${redirectUri}` +
			`&scope=${scope}` +
			`&state=${state}` +
			`&code_challenge=${codeChallenge}` +
			`&code_challenge_method=S256`;
	}
</script>

<button
	type="button"
	class="flex w-full cursor-pointer items-center justify-center gap-2 border border-black bg-stone-100 px-4 py-2 text-black transition-none hover:bg-black hover:text-white dark:border-white dark:bg-stone-900 dark:text-white dark:hover:bg-white dark:hover:text-stone-900"
	onclick={handleGitHubLogin}
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
			d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"
		></path>
	</svg>
	<span class="text-sm">{m.login_with_github()}</span>
</button>
