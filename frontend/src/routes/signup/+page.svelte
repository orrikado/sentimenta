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
		} catch (error) {
			formError = 'An error occurred. Please try again.';
		}
	}
</script>

<div
	class="flex min-h-screen items-center justify-center bg-stone-100 text-black dark:bg-stone-950"
>
	<div class="w-full max-w-md rounded-lg bg-white p-6 shadow-md dark:bg-stone-900">
		<h1 class="mb-6 text-center text-2xl font-bold text-gray-800 dark:text-gray-100">Sign Up</h1>

		<form on:submit={handleSubmit} class="space-y-4">
			<div>
				<label for="username" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
					>Username</label
				>
				<input
					type="text"
					id="username"
					name="username"
					required
					class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-800 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
				/>
			</div>

			<div>
				<label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
					>Email</label
				>
				<input
					type="email"
					id="email"
					name="email"
					required
					class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-800 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
				/>
			</div>

			<div>
				<label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
					>Password</label
				>
				<input
					type="password"
					id="password"
					name="password"
					required
					minlength="8"
					class="mt-1 block w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-gray-800 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm dark:border-gray-600 dark:bg-gray-700 dark:text-gray-100"
				/>
			</div>

			<button
				type="submit"
				class="w-full rounded-md bg-yellow-300 px-4 py-2 text-black hover:bg-yellow-200 focus:ring focus:ring-blue-300 focus:outline-none"
			>
				Sign Up
			</button>

			{#if formError}
				<p class="mt-2 text-sm text-red-600 dark:text-red-400">{formError}</p>
			{/if}
		</form>
	</div>
</div>
