import { jwtDecode } from 'jwt-decode';
import { user, userId } from './stores/user';

export type User = {
	username: string;
	email: string;
	created_at: Date;
	updated_at: Date;
};

function getCookie(name: string) {
	const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
	return match ? match[2] : null;
}

export function logout() {
	deleteCookie('access_token');
	userId.set(undefined);
}

function deleteCookie(name: string) {
	document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
}

export function isTokenExpired(token: string | null) {
	if (!token) return true; // No token means expired or invalid
	try {
		const decodedToken = jwtDecode(token);
		const currentTime = Date.now() / 1000; // Current time in seconds
		if (decodedToken.exp === undefined) {
			return true; // Token has no expiration time
		}
		return decodedToken.exp < currentTime; // Compare expiration time with current time
	} catch (error) {
		console.error('Error decoding token:', error);
		return true; // Assume expired if decoding fails
	}
}

export function refreshUserId() {
	if (typeof window === 'undefined') return;

	const jwtToken = getCookie('access_token');

	if (isTokenExpired(jwtToken)) {
		deleteCookie('access_token');
		userId.set(undefined);
		return;
	}

	if (jwtToken) {
		try {
			const { sub } = jwtDecode(jwtToken);
			userId.set(sub);
		} catch (error) {
			console.error('Invalid JWT:', error);
			userId.set(undefined);
		}
	} else {
		userId.set(undefined);
	}
}

export async function refreshUser() {
	if (!userId) return;
	const response = await fetch(`/api/user/get`);
	if (!response.ok) throw new Error('Failed to fetch user');
	user.set(await response.json());
}
