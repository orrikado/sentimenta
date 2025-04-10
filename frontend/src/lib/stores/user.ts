import { jwtDecode } from 'jwt-decode';
import { writable, type Writable } from 'svelte/store';

function getCookie(name: string) {
	const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
	return match ? match[2] : null;
}

export const userId: Writable<string | undefined> = writable(undefined);

if (typeof window !== 'undefined') {
	const jwtToken = getCookie('access_token');

	if (jwtToken) {
		try {
			userId.set(jwtDecode(jwtToken).sub);
		} catch (error) {
			console.error('Invalid JWT:', error);
			userId.set(undefined);
		}
	}
}
