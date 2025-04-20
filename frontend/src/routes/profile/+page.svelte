<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { userId } from '$lib/stores/user';
	import { logout } from '$lib/user';
	import { m } from '$lib/paraglide/messages';

	let user: {
		username: string;
		email: string;
		created_at: string | number | Date;
		updated_at: string | number | Date;
	} = $state({ username: '', email: '', created_at: '', updated_at: '' });

	let editMode = $state(false);
	let passwordEditMode = $state(false);
	let error = $state<string | null>(null);
	let success = $state<string | null>(null);
	let tempUser = $state({ username: '', email: '' });
	let verifyPassword = $state('');
	let passwords = $state({ current: '', new: '', confirm: '' });
	let showPassword = $derived(tempUser.email != user.email);

	onMount(async () => {
		if (!$userId) return goto('/');
		await loadUser();
	});

	async function loadUser() {
		try {
			const response = await fetch(`/api/user/get`);
			if (!response.ok) throw new Error('Failed to fetch user');
			user = await response.json();
			tempUser = { username: user.username, email: user.email };
		} catch (err) {
			console.error('Error fetching user:', err);
			goto('/login');
		}
	}

	// Format date to something readable
	function formatDate(str: string | number | Date) {
		const d = new Date(str);
		return d.toLocaleString();
	}
</script>

<svelte:head>
	<title>Sentimenta | Profile</title>
</svelte:head>

{#if user.email.length > 0}
	<div class="flex min-h-screen items-center justify-center bg-stone-100 px-4 dark:bg-stone-950">
		<div class="w-full max-w-md space-y-6">
			<!-- Profile Card -->
			<div class="border border-stone-300 bg-white p-6 dark:border-stone-700 dark:bg-stone-900">
				<h1 class="mb-4 text-center text-2xl font-bold">{m.profile()}</h1>

				{#if error}
					<div class="mb-4 bg-red-100 p-3 text-red-700 dark:bg-red-900 dark:text-red-100">
						{error}
					</div>
				{/if}

				{#if success}
					<div class="mb-4 bg-green-100 p-3 text-green-700 dark:bg-green-900 dark:text-green-100">
						{success}
					</div>
				{/if}

				<div class="space-y-4">
					<div>
						<label class="text-stone-500 uppercase dark:text-stone-400">{m.username()}</label>
						{#if editMode}
							<input
								bind:value={tempUser.username}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						{:else}
							<p class="text-lg">{user.username}</p>
						{/if}
					</div>

					<div>
						<label class="text-stone-500 uppercase dark:text-stone-400">Email</label>
						{#if editMode}
							<input
								bind:value={tempUser.email}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						{:else}
							<p class="text-lg">{user.email}</p>
						{/if}
					</div>

					{#if showPassword}
						<div>
							<label class="text-stone-500 uppercase dark:text-stone-400">
								<!-- {m.password()} -->
								password
							</label>
							<input
								type="password"
								bind:value={verifyPassword}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						</div>
					{/if}

					<div>
						<label class="text-stone-500 uppercase dark:text-stone-400">{m.created_at()}</label>
						<p class="text-lg">{formatDate(user.created_at)}</p>
					</div>

					<div>
						<label class="text-stone-500 uppercase dark:text-stone-400">{m.last_updated()}</label>
						<p class="text-lg">{formatDate(user.updated_at)}</p>
					</div>
				</div>

				<div class="mt-6 flex justify-end gap-2">
					{#if editMode}
						<button
							class="px-4 py-2 text-sm uppercase hover:underline"
							onclick={() => {
								editMode = false;
								if (user !== undefined) {
									tempUser = { username: user.username, email: user.email };
								}
							}}
						>
							<!-- {m.cancel()} -->
							cancel
						</button>
						<button
							class="bg-black px-4 py-2 text-sm text-white uppercase hover:bg-stone-800 dark:bg-white dark:text-black dark:hover:bg-stone-200"
							onclick={async () => {
								let body: {
									username: string | undefined;
									email: string | undefined;
									password: string | undefined;
								} = {
									username: undefined,
									email: undefined,
									password: undefined
								};

								if (tempUser.username !== user.username) {
									body.username = tempUser.username;
								}
								if (tempUser.email !== user.email) {
									if (verifyPassword.length == 0) {
										error = 'Password is required';
										return;
									}
									body.email = tempUser.email;
									body.password = verifyPassword;
								}

								console.log(body);
								console.log(JSON.stringify(body));

								let response = await fetch(`/api/user/update`, {
									method: 'PATCH',
									headers: {
										'Content-Type': 'application/json'
									},
									body: JSON.stringify(body)
								});

								if (!response.ok) {
									error = response.statusText;
									return;
								} else {
									success = 'Profile updated successfully';
									error = '';
								}
							}}
						>
							<!-- {m.save()} -->
							save
						</button>
					{:else}
						<button
							class="border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
							onclick={() => {
								editMode = true;

								if (user !== undefined) {
									tempUser = { username: user.username, email: user.email };
								}
							}}
						>
							<!-- {m.edit_profile()} -->
							edir profile
						</button>
					{/if}
				</div>
			</div>

			<!-- Password Update Card -->
			<div class="border border-stone-300 bg-white p-6 dark:border-stone-700 dark:bg-stone-900">
				<h2 class="mb-4 text-xl font-bold">
					<!-- {m.change_password()} -->
					change password
				</h2>

				{#if passwordEditMode}
					<div class="space-y-4">
						<div>
							<label class="text-stone-500 uppercase dark:text-stone-400">
								<!-- {m.current_password()} -->
								current password
							</label>
							<input
								type="password"
								bind:value={passwords.current}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						</div>

						<div>
							<label class="text-stone-500 uppercase dark:text-stone-400">
								<!-- {m.new_password()} -->
								new password
							</label>
							<input
								type="password"
								bind:value={passwords.new}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						</div>

						<div>
							<label class="text-stone-500 uppercase dark:text-stone-400">
								<!-- {m.confirm_password()} -->
								confirm password
							</label>
							<input
								type="password"
								bind:value={passwords.confirm}
								class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
							/>
						</div>

						<div class="flex justify-end gap-2">
							<button
								class="px-4 py-2 text-sm uppercase hover:underline"
								onclick={() => (passwordEditMode = false)}
							>
								<!-- {m.cancel()} -->
								cancel
							</button>
							<button
								class="bg-black px-4 py-2 text-sm text-white uppercase hover:bg-stone-800 dark:bg-white dark:text-black dark:hover:bg-stone-200"
							>
								<!-- {m.save()} -->
								save
							</button>
						</div>
					</div>
				{:else}
					<button
						class="w-full border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
						onclick={() => (passwordEditMode = true)}
					>
						<!-- {m.change_password()} -->
						change password
					</button>
				{/if}
			</div>

			<!-- Logout Button -->
			<div class="text-right">
				<button
					class="border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
					onclick={() => {
						logout();
						goto('/login');
					}}
				>
					{m.logout()}
				</button>
			</div>
		</div>
	</div>
{/if}
