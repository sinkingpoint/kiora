import { h, Fragment } from "preact";
import Button from "../../components/button";
import style from "./styles.css";
import { ChangeEvent, useState } from "preact/compat";
import { PreviewPageProps } from "./preview";
import { getSilenceEnd, parseMatcher } from "./utils";

// LabelMatcher takes a matcher string and returns a span element that displays the matcher.
const LabelMatcher = (matcher: string, onDelete: () => void) => {
	const { label, value, isRegex, isNegative } = parseMatcher(matcher);
	let operator = "";

	if (isRegex) {
		operator = isNegative ? "!~" : "=~";
	} else {
		operator = isNegative ? "!=" : "=";
	}

	return (
		<span className={style["label-matcher"]}>
			{label} {operator} {value}
			<button type="button" onClick={onDelete} class={style["delete-label-button"]}>
				🞫
			</button>
		</span>
	);
};

// setFilterInURL sets the filter query parameter in the URL to the given matchers.
const setFilterInURL = (matchers: string[]) => {
	const params = new URLSearchParams(window.location.search);
	params.delete("filter");
	matchers.forEach((filter) => {
		params.append("filter", filter);
	});
	window.history.replaceState({}, "", `${window.location.pathname}?${params.toString()}`);
};

// checkFormValidity checks if the form is valid and displays errors if it is not.
const checkFormValidity = () => {
	const duration = document.getElementById("duration") as HTMLInputElement;
	const creator = document.getElementById("creator") as HTMLInputElement;
	const comment = document.getElementById("comment") as HTMLInputElement;

	if (getSilenceEnd(duration.value) === null) {
		duration.setCustomValidity("Invalid duration");
		duration.reportValidity();
	} else if (creator.value === "") {
		creator.setCustomValidity("Creator cannot be empty");
		creator.reportValidity();
	} else if (comment.value === "") {
		comment.setCustomValidity("Comment cannot be empty");
		comment.reportValidity();
	}

	return true;
};

interface CreatePageProps {
	onPreview: (p: PreviewPageProps) => void;
}

const CreatePage = ({ onPreview }: CreatePageProps) => {
	const params = new URLSearchParams(window.location.search);
	const [duration, setDuration] = useState<string>("1h");
	const [matchers, setMatchers] = useState<string[]>(params.getAll("filter"));
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

	// filterSpans is an array of spans that display the label matchers.
	const filterSpans = matchers.map((filter, i) => {
		return LabelMatcher(filter, () => {
			const newFilters = [...matchers];
			newFilters.splice(i, 1);
			setFilterInURL(newFilters);
			setMatchers(newFilters);
		});
	});

	return (
		<>
			<div>
				<label>Duration</label>
			</div>

			<div style={{ justifyContent: "space-between", flexWrap: "wrap" }}>
				<input
					id="duration"
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
					id="label-filter"
					type="text"
					title='Label filter, e.g foo="bar"'
					pattern='[a-zA-Z_]+(=~|!=|!~|=)".+"'
					onKeyPress={(e) => {
						if (e.key === "Enter") {
							if (parseMatcher(e.currentTarget.value) === null) {
								return;
							}

							const newFilters = [...matchers];
							newFilters.push(e.currentTarget.value);
							setFilterInURL(newFilters);
							setMatchers(newFilters);
						}
					}}
				/>
			</div>

			<div style={{ flexWrap: "wrap" }}>{filterSpans}</div>

			<div>
				<label>Creator</label>
			</div>

			<div>
				<input id="creator" type="text" required />
			</div>

			<div>
				<label>Comment</label>
			</div>

			<div>
				<input id="comment" type="text" required />
			</div>

			<div style={{ marginTop: "20px", flexDirection: "row" }}>
				<Button
					label="Preview"
					onClick={() =>
						checkFormValidity() &&
						onPreview({
							duration,
							creator: (document.getElementById("creator") as HTMLInputElement).value,
							comment: (document.getElementById("comment") as HTMLInputElement).value,
							matchers,
						})
					}
				/>
			</div>
		</>
	);
};

export default CreatePage;
