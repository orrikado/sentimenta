import type { MoodEntry } from '$lib/moods';
import { writable, type Writable } from 'svelte/store';

export const moods: Writable<MoodEntry[]> = writable([]);
