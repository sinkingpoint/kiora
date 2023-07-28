import { h } from "preact";
import { ChangeEvent } from "preact/compat";
import { useState } from "preact/hooks";
import style from "./styles.css";

const LabelMatcher = (matcher: string, onDelete: () => void) => {
	const matcherMatches = matcher.match(/([a-zA-Z0-9_]+)(!=|!~|=|~)"(.*)"/);
	if (matcherMatches === null) {
		console.log("Invalid matcher", matcher);
		return null;
	}

	const labelName = matcherMatches[1];
	const operator = matcherMatches[2];
	const labelValue = matcherMatches[3];

	return (
		<span>
			<span className={style["label-matcher"]}>
				{labelName} {operator} {labelValue}
			</span>
			<button type="button" onClick={onDelete} class={style["delete-label-button"]}>
				x
			</button>
		</span>
	);
};

// getSilenceEnd takes a raw duration string and returns a Date object
// representing the end of the silence, if the duration is valid. Otherwise, it returns null.
const getSilenceEnd = (rawDuration: string): Date => {
	const durationMatches = rawDuration.match(/([0-9]+)([mhdw])/);
	if (durationMatches === null) {
		return null;
	}

	const durationAmt = parseInt(durationMatches[1]);
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

	const filterSpans = filters.map((filter, i) => {
		return LabelMatcher(filter, () => {
			const newFilters = [...filters];
			newFilters.splice(i, 1);
			setFilter(newFilters);
		});
	});

	return (
		<div>
			<h1>New Silence</h1>
			<form>
				<table>
					<tr>
						<td>
							<label>Duration</label>
						</td>
						<td>
							<input
								type="text"
								title="Duration in Go format, e.g. 1h"
								pattern="[0-9]+[mhdw]"
								value={duration}
								onInput={handleSetDuration}
							/>
						</td>
						<td>{end}</td>
					</tr>

					<tr>{filterSpans}</tr>
				</table>
			</form>
		</div>
	);
};

export default NewSilence;
