export function getMonthDays(year: number, month: number, firstDayOfWeek: number = 0): Date[] {
	const firstDayOfMonth = new Date(year, month, 1);
	const lastDayOfMonth = new Date(year, month + 1, 0).getDate();

	const days: Date[] = [];

	// Step 1: Calculate how many placeholder days (NaN) to insert before the 1st
	const firstDayIndex = firstDayOfMonth.getDay(); // 0 = Sunday
	const adjustedStartDay = (firstDayIndex - firstDayOfWeek + 7) % 7;

	for (let i = 0; i < adjustedStartDay; i++) {
		days.push(new Date(NaN)); // Placeholder before the 1st
	}

	for (let i = 1; i <= lastDayOfMonth; i++) {
		days.push(new Date(year, month, i));
	}

	return days;
}
