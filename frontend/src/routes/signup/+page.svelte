<script lang="ts">
	import { goto } from '$app/navigation';
	import GithubLoginButton from '$lib/components/GithubLoginButton.svelte';
	import GoogleLoginButton from '$lib/components/GoogleLoginButton.svelte';
	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import { refreshUserId, setCookie } from '$lib/user';
	import { onMount } from 'svelte';
	import { env } from '$env/dynamic/public';

	let formError: string | null = $state(null);
	const email_regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

	let email = $state('');
	let password = $state('');
	let username = $state('');

	let submitInProcess = $state(false);

	let registrationDisabled = env.PUBLIC_REGISTRATION_ENABLED !== 'true';

	$effect(() => {
		if (registrationDisabled) {
			formError = '';
			return;
		}
		if (canSubmit()) {
			formError = '';
			return;
		}
		if (!email || !password || !username) {
			formError = '';
			return;
		}

		if (hasSpace(username)) {
			formError = m.username_spaces_error();
			return;
		} else if (!email_regex.test(email)) {
			formError = m.invalid_email_error();
			return;
		} else if (password.length < parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')) {
			formError = m.password_too_short_error();
			return;
		}
	});

	const canSubmit = $derived(() => {
		if (!email || !password || !username) return false;
		if (hasSpace(password)) return false;
		if (hasSpace(username)) return false;
		if (hasSpace(email)) return false;
		if (!email_regex.test(email)) return false;
		if (password.length < parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')) return false;

		return true;
	});

	function hasSpace(str: string) {
		return str.indexOf(' ') >= 0;
	}

	onMount(() => {
		if ($userId) {
			goto('/track');
		}
	});

	async function handleSubmit(event: Event) {
		event.preventDefault();

		// 🔒 Prevent submission if registration is disabled
		if (registrationDisabled) {
			formError = m.registration_disabled();
			return;
		}

		formError = null;
		submitInProcess = true;

		const form = event.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());
		const { username, email, password } = formData;

		if (!username || !email || !password) {
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
			const response = await fetch('/api/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					username,
					email,
					password,
					timezone: Intl.DateTimeFormat().resolvedOptions().timeZone
				})
			});

			if (!response.ok) {
				const errorData = await response.json();
				formError = errorData.error || m.signup_fail();
				submitInProcess = false;
				return;
			}

			const data = await response.json();
			if (data.token) {
				setCookie('access_token', data.token, 30);
			}
			submitInProcess = false;
			refreshUserId();
			localStorage.setItem('justRegistered', 'true');
			goto('/track');
		} catch {
			submitInProcess = false;
			formError = m.error_occured();
		}
	}
</script>

<svelte:head>
	<title>Sentimenta | Signup</title>
</svelte:head>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<main
		class="flex w-full max-w-md flex-col justify-center gap-6 border border-stone-300 bg-stone-100 p-6 dark:border-white/10 dark:bg-stone-900"
	>
		<h1 class="text-center text-2xl">{m.m_signup()}</h1>

		<form onsubmit={handleSubmit} class="flex flex-col gap-4">
			<!-- username -->
			<div class="flex flex-col gap-1">
				<label for="username" class="text-sm">{m.username()}</label>
				<input
					id="username"
					name="username"
					required
					bind:value={username}
					class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
					placeholder="ex: cursed_dude42"
					disabled={registrationDisabled || submitInProcess}
				/>
			</div>

			<!-- email -->
			<div class="flex flex-col gap-1">
				<label for="email" class="text-sm">email</label>
				<input
					type="email"
					id="email"
					name="email"
					required
					bind:value={email}
					class="w-full border border-current bg-transparent px-3 py-2 placeholder-current"
					placeholder="email@proton.me"
					disabled={registrationDisabled || submitInProcess}
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
					disabled={registrationDisabled || submitInProcess}
				/>
			</div>

			<!-- submit -->
			<button
				type="submit"
				class="relative flex items-center justify-center border border-current px-6 py-2 text-center transition-colors duration-300"
				class:bg-white={canSubmit() && !registrationDisabled}
				class:text-black={canSubmit() && !registrationDisabled}
				class:opacity-50={registrationDisabled || !canSubmit() || submitInProcess}
				class:pointer-events-none={registrationDisabled || !canSubmit() || submitInProcess}
				aria-disabled={registrationDisabled || !canSubmit() || submitInProcess}
				disabled={registrationDisabled || !canSubmit() || submitInProcess}
			>
				{m.signup()}
			</button>

			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
		</form>
		<div class="flex flex-col gap-1">
			<GoogleLoginButton />
			<GithubLoginButton />
		</div>
		<a href="/login" class="text-center text-sm text-gray-500 underline dark:text-gray-400">
			{m.already_have_account()}
		</a>

		<span class="text-sm text-gray-500 dark:text-gray-400">
			{m.privacy_policy_notice()}
			<a href="/privacy" class="text-center underline">{m.privacyPolicy_title()}</a>
		</span>
	</main>

	<!-- 🔒 Overlay when registration is disabled -->
	{#if registrationDisabled}
		<div class="bg-opacity-60 fixed inset-0 z-10 flex items-center justify-center backdrop-blur-sm">
			<p class="text-center text-xl text-white">{m.registration_disabled()}</p>
		</div>
	{/if}
</div>
