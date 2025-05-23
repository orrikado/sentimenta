import { server_status } from './stores/server_status';

export async function refreshServerStatus() {
	const response = await fetch(`/api/status`);
	if (!response.ok) server_status.set(false);
	else server_status.set(true);
}
