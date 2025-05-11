import { refreshUserId, type User } from '$lib/user';
import { writable, type Writable } from 'svelte/store';

export const userId: Writable<string | undefined> = writable(undefined);
export const user: Writable<User | undefined> = writable(undefined);

// run once on load
refreshUserId();
