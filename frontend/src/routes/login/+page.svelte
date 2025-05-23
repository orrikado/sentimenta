<script lang="ts">
	import { goto } from '$app/navigation';
	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import { refreshUserId } from '$lib/user';
	import { onMount } from 'svelte';

	import GoogleLoginButton from '$lib/components/GoogleLoginButton.svelte';
	import GithubLoginButton from '$lib/components/GithubLoginButton.svelte';
	import { env } from '$env/dynamic/public';

	let submitInProcess = $state(false);

	let formError: string | null = $state(null);

	let email = $state('');
	let password = $state('');

	const canSubmit = $derived(() => {
		if (!email || !password) return false;
		if (hasSpace(password)) return false;
		if (!email_regex.test(email)) return false;
		if (email.length > 255) return false;
		if (password.length > 255) return false;
		if (password.length < parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')) return false;

		return true;
	});

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
		submitInProcess = true;

		const form = event.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());
		const { email, password } = formData;

		if (!email || !password) {
			formError = m.field_error();
			submitInProcess = false;
			return;
		}
		if (hasSpace(password.toString())) {
			formError = m.password_spaces_error();
			submitInProcess = false;
			return;
		}
		if (!email_regex.test(email.toString())) {
			formError = m.invalid_email_error();
			submitInProcess = false;
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
				formError = errorData.error || m.login_fail();
				console.error(errorData);
				submitInProcess = false;
				return;
			}

			submitInProcess = false;
			refreshUserId();
			goto('/track');
		} catch {
			submitInProcess = false;
			formError = m.error_occured();
		}
	}
</script>

<svelte:head>
	<title>Sentimenta | Login</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<main
		class="flex w-full max-w-md flex-col justify-center gap-6 border border-stone-300 bg-stone-100 p-6 dark:border-white/10 dark:bg-stone-900"
	>
		<h1 class="text-center text-2xl">{m.m_login()}</h1>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<!-- email -->
			<div class="flex flex-col gap-1">
				<label for="email" class="text-sm">{m.email()}</label>
				<input
					type="email"
					id="email"
					name="email"
					required
					bind:value={email}
					class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
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
					bind:value={password}
					minlength={parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')}
					class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
					placeholder="••••••••"
				/>
			</div>

			<!-- submit -->
			<button
				type="submit"
				class="relative flex items-center justify-center border border-current px-6 py-2 text-center transition-colors duration-300"
				class:bg-white={canSubmit()}
				class:text-black={canSubmit()}
				class:cursor-pointer={canSubmit()}
				class:opacity-50={!canSubmit()}
				class:pointer-events-none={!canSubmit() || submitInProcess}
				aria-disabled={!canSubmit() || submitInProcess}
				disabled={submitInProcess || !canSubmit()}
			>
				{m.login()}
			</button>

			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
		</form>
		<div class="flex flex-col gap-1">
			<GoogleLoginButton />
			<GithubLoginButton />
		</div>
		<span class="text-sm text-gray-500 dark:text-gray-400">
			{m.privacy_policy_notice()}
			<a href="/privacy" class="text-center underline">{m.privacyPolicy_title()}</a>
		</span>
	</main>
</div>
