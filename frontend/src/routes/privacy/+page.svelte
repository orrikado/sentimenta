<script>
	import { m } from '$lib/paraglide/messages';

	// Dynamic date (replace [Date] in "Last Updated")
	const lastUpdatedDate = new Date('2025-05-23').toLocaleDateString(undefined, {
		year: 'numeric',
		month: 'long',
		day: 'numeric'
	});
</script>

<div class="mx-auto max-w-3xl px-4 py-10">
	<h1 class="mb-2 text-3xl font-bold">{m.privacyPolicy_title()}</h1>
	<p class="mb-6 text-sm text-gray-500">
		{m.privacyPolicy_lastUpdated({ date: lastUpdatedDate })}
	</p>

	<p class="mb-6">{m.privacyPolicy_intro()}</p>

	<!-- Section 1 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_1_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_1_item_1()}</li>
			<li>{m.privacyPolicy_section_1_item_2()}</li>
			<li>{m.privacyPolicy_section_1_item_3()}</li>
			<li>{m.privacyPolicy_section_1_item_4()}</li>
		</ul>
	</section>

	<!-- Section 2 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_2_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_2_item_1()}</li>
			<li>{m.privacyPolicy_section_2_item_2()}</li>
			<li>{m.privacyPolicy_section_2_item_3()}</li>
			<li>{m.privacyPolicy_section_2_item_4()}</li>
		</ul>
	</section>

	<!-- Section 3 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_3_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_3_item_1()}</li>
			<li>{m.privacyPolicy_section_3_item_2()}</li>
			<li>{m.privacyPolicy_section_3_item_3()}</li>
			<li>{m.privacyPolicy_section_3_item_4()}</li>
		</ul>
	</section>

	<!-- Section 4 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_4_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_4_item_1()}</li>
			<li>{m.privacyPolicy_section_4_item_2()}</li>
			<li>{m.privacyPolicy_section_4_item_3()}</li>
		</ul>
	</section>

	<!-- Section 5 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_5_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_5_item_1()}</li>
			<li>{m.privacyPolicy_section_5_item_2()}</li>
			<li>{m.privacyPolicy_section_5_item_3()}</li>
			<li>{m.privacyPolicy_section_5_item_4()}</li>
		</ul>
	</section>

	<!-- Section 6 -->
	<section class="mb-8">
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_6_title()}</h2>
		<ul class="list-disc space-y-2 pl-5">
			<li>{m.privacyPolicy_section_6_item_1()}</li>
			<li>{m.privacyPolicy_section_6_item_2()}</li>
			<li>{m.privacyPolicy_section_6_item_3()}</li>
		</ul>
	</section>

	<!-- Section 7 -->
	<section>
		<h2 class="mb-3 text-xl font-semibold">{m.privacyPolicy_section_7_title()}</h2>
		<p>{m.privacyPolicy_section_7_text()}</p>
	</section>
</div>
