import { userId } from '$lib/stores/user';
import { refreshUserId } from '$lib/user';
import { moods } from './stores/moods';

export type MoodEntry = {
	uid: number | undefined;
	date: Date;
	score: number;
	description: string;
	emotions: string;
};

export async function updateMoods() {
	const res = await fetch('/api/moods/get');
	if (res.ok) {
		const data = await res.json();
		const parsed = data.map((m: { date: string | number | Date }) => ({
			...m,
			date: new Date(m.date)
		}));
		moods.set(parsed);
	} else {
		refreshUserId();
		if (!userId) throw new Error('not logged in');
	}
}
