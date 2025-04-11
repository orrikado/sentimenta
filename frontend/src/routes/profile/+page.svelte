<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { userId } from '$lib/stores/user';
	import { logout } from '$lib/user';
	import { m } from '$lib/paraglide/messages';

	let user: {
		username: any;
		email: any;
		created_at: string | number | Date;
		updated_at: string | number | Date;
	};

	onMount(async () => {
		if ($userId == undefined) {
			goto('/');
		} else {
			try {
				let response = await fetch(`/api/user/${$userId}`);
				user = await response.json();
			} catch (error) {
				console.error('Error fetching user:', error);
				goto('/login');
			}
		}
	});

	// Format date to something readable
	function formatDate(str: string | number | Date) {
		const d = new Date(str);
		return d.toLocaleString();
	}
</script>

{#if user}
	<div class="flex min-h-screen items-center justify-center bg-stone-100 px-4 dark:bg-stone-950">
		<div
			class="w-full max-w-md border border-stone-300 bg-white p-6 font-mono text-sm dark:border-stone-700 dark:bg-stone-900 dark:text-white"
		>
			<h1 class="mb-4 text-center text-2xl font-bold tracking-tight">{m.profile()}</h1>

			<div class="space-y-4">
				<div>
					<p class="text-stone-500 uppercase dark:text-stone-400">{m.username()}</p>
					<p class="text-lg">{user.username}</p>
				</div>

				<div>
					<p class="text-stone-500 uppercase dark:text-stone-400">email</p>
					<p class="text-lg">{user.email}</p>
				</div>

				<div>
					<p class="text-stone-500 uppercase dark:text-stone-400">{m.created_at()}</p>
					<p class="text-lg">{formatDate(user.created_at)}</p>
				</div>

				<div>
					<p class="text-stone-500 uppercase dark:text-stone-400">{m.last_updated()}</p>
					<p class="text-lg">{formatDate(user.updated_at)}</p>
				</div>
			</div>

			<div class="mt-6 flex justify-end">
				<button
					class="border border-stone-700 bg-transparent px-4 py-2 text-sm tracking-wide text-black uppercase hover:bg-black hover:text-white dark:text-white dark:hover:bg-white dark:hover:text-black"
					on:click={() => {
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
