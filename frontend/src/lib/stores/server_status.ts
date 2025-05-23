import { writable, type Writable } from 'svelte/store';

export const server_status: Writable<boolean> = writable(false);
