<script lang="ts">
	import { goto } from '$app/navigation';
	import GoogleLoginButton from '$lib/components/GoogleLoginButton.svelte';
	import { m } from '$lib/paraglide/messages';
	import { userId } from '$lib/stores/user';
	import { refreshUserId } from '$lib/user';
	import { onMount } from 'svelte';

	let formError: string | null = null;
	const email_regex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

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
		formError = null;

		const form = event.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());
		const { username, email, password } = formData;

		if (!username || !email || !password) {
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
			const response = await fetch('/api/auth/register', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ username, email, password })
			});

			if (!response.ok) {
				const errorData = await response.json();
				formError = errorData.message || m.signup_fail();
				return;
			}

			refreshUserId();
			goto('/profile');
		} catch {
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
	<main class="flex w-full max-w-md flex-col gap-6 border border-black p-6 dark:border-white">
		<h1 class="text-center text-2xl">{m.m_signup()}</h1>

		<form on:submit={handleSubmit} class="flex flex-col gap-4">
			<!-- username -->
			<div class="flex flex-col gap-1">
				<label for="username" class="text-sm">{m.username()}</label>
				<input
					id="username"
					name="username"
					required
					class="border border-black bg-stone-100 px-3 py-2 text-black placeholder-gray-500 dark:border-white dark:bg-stone-950 dark:text-white dark:placeholder-white"
					placeholder="ex: cursed_dude42"
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
				{m.signup()}
			</button>

			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
		</form>
		<GoogleLoginButton />
	</main>
</div>
