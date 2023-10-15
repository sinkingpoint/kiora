// formatDate formats the given date into a string consistent across the UI.
export const formatDate = (d: Date): string => {
	return d.toLocaleString([], {
		day: "numeric",
		month: "short",
		year: "numeric",
		hour: "2-digit",
		minute: "2-digit",
	});
};

// formatDuration formats the given duration in seconds into a human readable string.
export const formatDuration = (seconds: number): string => {
	if (seconds < 60) {
		return `${seconds}s`;
	}

	const minutes = Math.floor(seconds / 60);

	if (minutes < 60) {
		return `${minutes}m`;
	}

	const hours = Math.floor(minutes / 60);
	const remainingMinutes = minutes % 60;

	if (hours < 24) {
		return `${hours}h ${remainingMinutes}m`;
	}

	const days = Math.floor(hours / 24);
	const remainingHours = hours % 24;

	return `${days}d ${remainingHours}h ${remainingMinutes}m`;
};
