// PKCE Helpers
export function generateRandomString() {
	return window.crypto.randomUUID();
}

export function generateCodeVerifier() {
	const array = new Uint8Array(32);
	window.crypto.getRandomValues(array);
	return btoa(String.fromCharCode(...array)).replace(/[^a-zA-Z0-9]/g, '');
}

export async function generateCodeChallenge(codeVerifier: string) {
	const encoder = new TextEncoder();
	const digestBuffer = await crypto.subtle.digest('SHA-256', encoder.encode(codeVerifier));

	// Convert buffer to Uint8Array, then to Base64 string
	const digestArray = Array.from(new Uint8Array(digestBuffer));
	const base64Digest = btoa(String.fromCharCode(...digestArray));

	// Replace + -> -, / -> _, and remove padding (=)
	const urlSafeBase64 = base64Digest.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');

	return urlSafeBase64;
}
