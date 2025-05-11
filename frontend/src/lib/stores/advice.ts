import type { AdviceEntry } from '$lib/advice';
import { writable, type Writable } from 'svelte/store';

export const advice: Writable<AdviceEntry[]> = writable([]);
