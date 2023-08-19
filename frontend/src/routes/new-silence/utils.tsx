// getSilenceEnd takes a raw duration string and returns a Date object

import { Matcher } from "../../api";

// representing the end of the silence, if the duration is valid. Otherwise, it returns null.
export const getSilenceEnd = (rawDuration: string): Date => {
	const durationMatches = rawDuration.match(/^([0-9]+)([mhdw])$/);
	if (durationMatches === null) {
		return null;
	}

	const durationAmt = parseInt(durationMatches[1], 10);

	// durationUnit is the unit of the duration, e.g. m for minutes, h for hours, d for days.
	const durationUnit = durationMatches[2];

	const now = new Date();
	switch (durationUnit) {
		case "m":
			now.setMinutes(now.getMinutes() + durationAmt);
			break;
		case "h":
			now.setHours(now.getHours() + durationAmt);
			break;
		case "d":
			now.setDate(now.getDate() + durationAmt);
			break;
		case "w":
			now.setDate(now.getDate() + durationAmt * 7);
			break;
		default:
			return null;
	}

	return now;
};
