import { h } from "preact";
import { ChangeEvent } from "preact/compat";
import { useState } from "preact/hooks";
import style from "./styles.css";

type MatcherOperation = "=" | "!=" | "=~" | "!~";

// Matcher is a struct that represents a label matcher.
interface Matcher {
	// name is the name of the label.
	name: string;

	// operator is the operator to use for the matcher.
	operator: MatcherOperation;

	// value is the value to match against.
	value: string;
}

// parseMatcher takes a matcher string and returns a Matcher object if the string is valid.
const parseMatcher = (matcher: string): Matcher | null => {
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
	if(operator.includes("~")) {
		try {
			new RegExp(matcherMatches[3]);
		}
		catch {
			return null;
		}
	}

	return {
		name: matcherMatches[1],
		operator: matcherMatches[2] as MatcherOperation,
		value: matcherMatches[3],
	};
};

// LabelMatcher takes a matcher string and returns a span element that displays the matcher.
const LabelMatcher = (matcher: string, onDelete: () => void) => {
	const { name: labelName, operator, value: labelValue } = parseMatcher(matcher);

	return (
		<span className={style["label-matcher"]}>
			{labelName} {operator} {labelValue}
			<button type="button" onClick={onDelete} class={style["delete-label-button"]}>
				🞫
			</button>
		</span>
	);
};

// getSilenceEnd takes a raw duration string and returns a Date object
// representing the end of the silence, if the duration is valid. Otherwise, it returns null.
const getSilenceEnd = (rawDuration: string): Date => {
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

// setFilterInURL sets the filter query parameter in the URL to the given filters.
const setFilterInURL = (filters: string[]) => {
	const params = new URLSearchParams(window.location.search);
	params.delete("filter");
	filters.forEach((filter) => {
		params.append("filter", filter);
	});
	window.history.replaceState({}, "", `${window.location.pathname}?${params.toString()}`);
};

const NewSilence = () => {
	const params = new URLSearchParams(window.location.search);
	const [duration, setDuration] = useState<string>("1h");
	const [filters, setFilter] = useState<string[]>(params.getAll("filter"));

	const handleSetDuration = (e: ChangeEvent<HTMLInputElement>) => {
		setDuration(e.currentTarget.value);
	};

	const endDate = getSilenceEnd(duration);
	const end =
		endDate !== null ? (
			<span>
				Ends at{" "}
				{endDate.toLocaleString([], {
					day: "numeric",
					month: "short",
					year: "numeric",
					hour: "2-digit",
					minute: "2-digit",
				})}
			</span>
		) : (
			<span>Invalid duration</span>
		);

	// filterSpans is an array of spans that display the label filters.
	const filterSpans = filters.map((filter, i) => {
		return LabelMatcher(filter, () => {
			const newFilters = [...filters];
			newFilters.splice(i, 1);
			setFilterInURL(newFilters);
			setFilter(newFilters);
		});
	});

	return (
		<div class={style["silence-form"]}>
			<h1>New Silence</h1>
			<div>
				<label>Duration</label>
			</div>

			<div style={{ justifyContent: "space-between" }}>
				<input
					type="text"
					title="Duration in Go format, e.g. 1h"
					pattern="[0-9]+[mhdw]"
					value={duration}
					onInput={handleSetDuration}
				/>

				<label>{end}</label>
			</div>

			<div>
				<label>Label Filters</label>
			</div>

			<div>
				<input
					type="text"
					title='Label filter, e.g foo="bar"'
					pattern='[a-zA-Z_]+(=~|!=|!~|=)".+"'
					onKeyPress={(e) => {
						if (e.key === "Enter") {
							if (parseMatcher(e.currentTarget.value) === null) {
								return;
							}

							const newFilters = [...filters];
							newFilters.push(e.currentTarget.value);
							setFilterInURL(newFilters);
							setFilter(newFilters);

							e.currentTarget.value = "";
						}
					}}
				/>
			</div>

			<div style={{flexWrap: "wrap"}}>{filterSpans}</div>

			<div>
				<label>Creator</label>
			</div>

			<div>
				<input type="text" required />
			</div>

			<div>
				<label>Comment</label>
			</div>

			<div>
				<input type="text" required />
			</div>

			<div>
				<button type="button">Preview</button>
			</div>
		</div>
	);
};

export default NewSilence;
