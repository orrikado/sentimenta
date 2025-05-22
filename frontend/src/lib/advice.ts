import { advice } from './stores/advice';
import { userId } from './stores/user';
import { refreshUserId } from './user';

export type AdviceEntry = {
	uid: number | undefined;
	date: Date;
	text: string;
	generated_by_websocket: boolean;
};

export async function updateAdvice() {
	const res = await fetch('/api/advice');
	if (res.ok) {
		const data = await res.json();
		data.map((a: { generated_by_websocket: boolean }) => {
			a.generated_by_websocket = false;
		});
		advice.set(data);
	} else {
		console.error('Failed to fetch advice');
		refreshUserId();
		if (!userId) throw new Error('not logged in');
	}
}
