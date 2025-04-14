export function getMonthDays(year: number, month: number): Date[] {
	const days = [];
	const firstDay = new Date(year, month, 1);
	const startDay = firstDay.getDay(); // 0 = Sunday

	// Push empty placeholders before the 1st
	for (let i = 0; i < startDay; i++) {
		days.push(new Date(NaN));
	}

	const lastDay = new Date(year, month + 1, 0).getDate();

	for (let i = 1; i <= lastDay; i++) {
		days.push(new Date(year, month, i));
	}

	return days;
}
