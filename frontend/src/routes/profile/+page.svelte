<script lang="ts">
	import { onMount, tick } from 'svelte';
	import { goto } from '$app/navigation';
	import { user, userId } from '$lib/stores/user';
	import { logout, refreshUser } from '$lib/user';
	import { m } from '$lib/paraglide/messages';
	import Settings from '$lib/components/Settings.svelte';
	import { refreshServerStatus } from '$lib/status';
	import { server_status } from '$lib/stores/server_status';
	import { env } from '$env/dynamic/public';

	let editMode = $state(false);
	let passwordEditMode = $state(false);
	let error = $state<string | null>(null);
	let success = $state<string | null>(null);
	let tempUser = $state({ username: '', email: '' });
	let verifyPassword = $state('');
	let passwords = $state({ current: '', new: '', confirm: '' });
	let showPassword = $derived(tempUser.email != $user?.email);

	onMount(async () => {
		await refreshServerStatus();
		if (!$userId && $server_status) {
			goto('/login');
			return;
		}
		if (!$user && $server_status) {
			try {
				await refreshUser();
			} catch (err) {
				console.error(err);
				goto('/login');
			}
		}
		if (!$user && $server_status) {
			logout();
			goto('/');
		} else {
			tempUser = { username: $user?.username || '', email: $user?.email || '' };
		}
	});

	async function changePassword() {
		// Reset previous messages
		error = null;
		success = null;

		// Validate input fields
		if (!passwords.current || !passwords.new || !passwords.confirm) {
			error = m.profile_password_all_required();
			return;
		}

		if (hasSpace(passwords.new.toString())) {
			error = m.password_spaces_error();
			return;
		}

		if (passwords.new !== passwords.confirm) {
			error = m.profile_passwords_do_not_match();
			return;
		}

		// Send request to update password
		const response = await fetch('/api/user/update/password', {
			method: 'PUT',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				password: passwords.current,
				new_password: passwords.new
			})
		});

		// Handle response
		if (!response.ok) {
			const data = await response.json().catch(() => ({}));
			error = data.error || m.profile_password_update_failed();
		} else {
			success = m.profile_password_updated_successfully();
			error = null;
			passwords = { current: '', new: '', confirm: '' };
			passwordEditMode = false;
		}
	}

	// Format date to something readable
	function formatDate(str: string | number | Date) {
		const d = new Date(str);
		return d.toLocaleString();
	}

	function hasSpace(str: string) {
		return str.indexOf(' ') >= 0;
	}
</script>

<svelte:head>
	<title>Sentimenta | Profile</title>
</svelte:head>

{#if $user?.email}
	<div class="flex min-h-screen flex-col bg-stone-100 px-4 md:flex-row dark:bg-stone-950">
		<main class="mx-auto w-full max-w-2xl space-y-6 py-8">
			<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
				<!-- Left Column -->
				<div class=" flex flex-col space-y-6">
					<!-- Profile Card -->
					<div
						class="flex-grow border border-stone-300 bg-white p-6 dark:border-stone-700 dark:bg-stone-900"
					>
						<h1 class="mb-4 text-center text-2xl font-bold">{m.profile()}</h1>
						<!-- Status -->
						{#if error}
							<div class="mb-4 bg-red-100 p-3 text-red-700 dark:bg-red-900 dark:text-red-100">
								{error}
							</div>
						{/if}
						{#if success}
							<div
								class="mb-4 bg-green-100 p-3 text-green-700 dark:bg-green-900 dark:text-green-100"
							>
								{success}
							</div>
						{/if}
						<!-- Profile fields -->
						<div class="space-y-4">
							<!-- Username/email fields -->
							<div>
								<label for="username" class="text-stone-500 uppercase dark:text-stone-400"
									>{m.username()}</label
								>
								{#if editMode}
									<input
										bind:value={tempUser.username}
										id="username"
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
									/>
								{:else}
									<p class="text-lg">{$user.username}</p>
								{/if}
							</div>

							<div>
								<label for="email" class="text-stone-500 uppercase dark:text-stone-400">Email</label
								>
								{#if editMode}
									<input
										id="email"
										bind:value={tempUser.email}
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
									/>
								{:else}
									<p class="text-lg">{$user.email}</p>
								{/if}
							</div>

							{#if showPassword}
								<div>
									<label
										for="current-password"
										class="text-stone-500 uppercase dark:text-stone-400"
									>
										{m.password()}
									</label>
									<input
										type="password"
										id="current-password"
										aria-required="true"
										bind:value={verifyPassword}
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
									/>
								</div>
							{/if}

							<div>
								<span class="text-stone-500 uppercase dark:text-stone-400">{m.created_at()}</span>
								<p class="text-lg">{formatDate($user.created_at)}</p>
							</div>

							<div>
								<span class="text-stone-500 uppercase dark:text-stone-400">{m.last_updated()}</span>
								<p class="text-lg">{formatDate($user.updated_at)}</p>
							</div>
						</div>
						<!-- Edit buttons -->
					</div>

					<!-- Password Card -->
					<form
						onsubmit={(event) => {
							event.preventDefault();
							changePassword();
						}}
						class="border border-stone-300 bg-white p-6 dark:border-stone-700 dark:bg-stone-900"
					>
						<h2 class="mb-4 text-xl font-bold">
							{m.change_password()}
						</h2>

						{#if passwordEditMode}
							<div class="space-y-4">
								<div>
									<label
										for="current-password"
										class="text-stone-500 uppercase dark:text-stone-400"
									>
										{m.current_password()}
									</label>
									<input
										type="password"
										id="current-password"
										name="password"
										required
										minlength={parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')}
										bind:value={passwords.current}
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
										aria-required="true"
									/>
								</div>

								<div>
									<label for="new-password" class="text-stone-500 uppercase dark:text-stone-400">
										{m.new_password()}
									</label>
									<input
										id="new-password"
										aria-required="true"
										type="password"
										name="password"
										required
										minlength={parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')}
										bind:value={passwords.new}
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
									/>
								</div>

								<div>
									<label
										for="confirm-password"
										class="text-stone-500 uppercase dark:text-stone-400"
									>
										{m.confirm_password()}
									</label>
									<input
										id="confirm-password"
										aria-required="true"
										type="password"
										name="password"
										required
										minlength={parseInt(env.PUBLIC_PASSWORD_LENGTH_MIN || '8')}
										bind:value={passwords.confirm}
										class="w-full border-b border-stone-300 bg-transparent p-1 focus:outline-none dark:border-stone-600"
									/>
								</div>

								<div class="flex justify-end gap-2">
									<button
										class="px-4 py-2 text-sm uppercase hover:underline"
										onclick={() => {
											passwordEditMode = false;
											error = null;
											success = null;
										}}
									>
										{m.cancel()}
									</button>
									<button
										type="submit"
										class="bg-black px-4 py-2 text-sm text-white uppercase hover:bg-stone-800 dark:bg-white dark:text-black dark:hover:bg-stone-200"
									>
										{m.save_password()}
									</button>
								</div>
							</div>
						{:else}
							<button
								class="w-full border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
								onclick={() => (passwordEditMode = true)}
							>
								{m.change_password()}
							</button>
						{/if}
					</form>
				</div>

				<!-- Right Column -->
				<div class=" flex flex-col space-y-6">
					<!-- Settings Card -->
					<div
						class=" flex-grow border border-stone-300 bg-white p-6 dark:border-stone-700 dark:bg-stone-900"
					>
						<h2 class="mb-4 text-xl font-bold">{m.settings()}</h2>
						<Settings />
					</div>

					<!-- Account Actions -->
					<div
						class="gap-2 border border-stone-300 bg-white p-6 **:w-full dark:border-stone-700 dark:bg-stone-900"
					>
						<h2 class="mb-4 text-xl font-bold">{m.account()}</h2>
						<!-- Edit Profile -->
						<div class="mb-2 flex justify-end gap-2">
							{#if editMode}
								<button
									class="px-4 py-2 text-sm uppercase hover:underline"
									onclick={() => {
										editMode = false;
										if ($user !== undefined) {
											tempUser = { username: $user.username, email: $user.email };
										}
									}}
								>
									{m.cancel()}
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

										if (tempUser.username !== $user.username) {
											body.username = tempUser.username;
										}
										if (tempUser.email !== $user.email) {
											if (verifyPassword.length == 0) {
												error = 'Password is required';
												return;
											}
											body.email = tempUser.email;
											body.password = verifyPassword;
										}

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
									{m.save()}
								</button>
							{:else}
								<button
									class="border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
									onclick={async () => {
										editMode = true;

										if ($user) {
											tempUser = { username: $user.username, email: $user.email };
										}
										// Focus the first input after the DOM updates
										await tick();
										document.getElementById('username')?.focus();
									}}
								>
									{m.edit_profile()}
								</button>
							{/if}
						</div>
						<!-- logout -->
						<div class="flex justify-end">
							<button
								class="w-full border border-stone-700 px-4 py-2 text-sm uppercase hover:bg-black hover:text-white dark:border-stone-300 dark:hover:bg-white dark:hover:text-black"
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
			</div>
		</main>
	</div>
{/if}
