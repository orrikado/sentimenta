<script>
	let name = '';
	let responseMessage = '';

	async function handleSubmit() {
		console.log('Submitted:', name);
		const res = await fetch('https://api.example.com/submit', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ name })
		});

		if (res.ok) {
			const data = await res.json();
			responseMessage = data.message;
		} else {
			responseMessage = 'Something went wrong.';
		}
	}
</script>

<div class="flex justify-center">
	<form on:submit|preventDefault={handleSubmit} class="rounded-lg">
		<input bind:value={name} class="rounded-4xl" placeholder="Enter your name" />
		<button type="submit">Submit</button>
	</form>

	{#if responseMessage}
		<p>{responseMessage}</p>
	{/if}
</div>
