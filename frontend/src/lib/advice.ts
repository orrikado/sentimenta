import { advice } from './stores/advice';
import { userId } from './stores/user';
import { refreshUserId } from './user';

export type AdviceEntry = {
	uid: number | undefined;
	date: Date;
	text: string;
};

export async function updateAdvice() {
	const res = await fetch('/api/advice');
	if (res.ok) {
		const data = await res.json();
		console.log(data);
		advice.set(data);
	} else {
		console.error('Failed to fetch advice');
		refreshUserId();
		if (!userId) throw new Error('not logged in');
	}
}
