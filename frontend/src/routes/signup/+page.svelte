<script lang="ts">
	let formError: string | null = null;
	const email_regex =
		/^(?:[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$/;

	function hasSpace(str: string) {
		return str.indexOf(' ') >= 0;
	}

	async function handleSubmit(event: Event) {
		event.preventDefault();
		formError = null;

		const form = event.target as HTMLFormElement;
		const formData = Object.fromEntries(new FormData(form).entries());
		const { username, email, password } = formData;

		if (!username || !email || !password) {
			formError = 'Please fill out all fields.';
			return;
		}
		if (hasSpace(password.toString())) {
			formError = "Your password can't contain spaces";
			return;
		}
		if (!email_regex.test(email.toString())) {
			formError = 'Invalid email';
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
				formError = errorData.message || 'Signup failed.';
				return;
			}

			alert('Signup successful!');
		} catch {
			formError = 'An error occurred. Please try again.';
		}
	}
</script>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 font-mono text-black dark:bg-stone-950 dark:text-white"
>
	<div class="flex w-full max-w-md flex-col gap-6 border border-black p-6 dark:border-white">
		<h1 class="text-center text-2xl">sign up</h1>

		<form on:submit={handleSubmit} class="flex flex-col gap-4">
			<!-- username -->
			<div class="flex flex-col gap-1">
				<label for="username" class="text-sm">username</label>
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
				<label for="password" class="text-sm">password</label>
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
				sign up
			</button>

			{#if formError}
				<p class="text-sm text-red-500 dark:text-red-400">{formError}</p>
			{/if}
		</form>
	</div>
</div>
