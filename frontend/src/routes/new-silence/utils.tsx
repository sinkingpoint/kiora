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

// parseMatcher takes a matcher string and returns a Matcher object if the string is valid.
export const parseMatcher = (matcher: string): Matcher | null => {
	const matcherMatches = matcher.match(/([a-zA-Z0-9_]+)(!=|!~|=~|=)"(.*)"/);
	if (matcherMatches === null) {
		return null;
	}

	const validOperators = ["=", "!=", "=~", "!~"];
	const operator = matcherMatches[2];
	if (!validOperators.includes(operator)) {
		return null;
	}

	// If the operator is a regex operator, check that the regex is valid.
	if (operator.includes("~")) {
		try {
			new RegExp(matcherMatches[3]);
		} catch {
			return null;
		}
	}

	const isRegex = operator.includes("~");
	const isNegative = operator.includes("!");

	return {
		label: matcherMatches[1],
		value: matcherMatches[3],
		isRegex,
		isNegative,
	};
};
