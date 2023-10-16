import { h } from "preact";
import style from "./style.css";
import { Matcher } from "../../api";

// parseMatcher takes a matcher string and returns a Matcher object if the string is valid.
export const parseMatcher = (matcher: string): Matcher | null => {
	const matcherMatches = matcher.match(/([a-zA-Z0-9_]+)(!=|!~|=~|=)"(.*)"/);
	if (matcherMatches === null) {
		console.log("Invalid matcher", matcher);
		return null;
	}

	const validOperators = ["=", "!=", "=~", "!~"];
	const operator = matcherMatches[2];
	if (!validOperators.includes(operator)) {
		console.log("Invalid operator for matcher", matcher);
		return null;
	}

	// If the operator is a regex operator, check that the regex is valid.
	if (operator.includes("~")) {
		try {
			new RegExp(matcherMatches[3]);
		} catch {
			console.log("Invalid regex for matcher", matcher);
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

export interface LabelMatcherCardProps {
	matcher: string | Matcher;
	onDelete?: () => void;
}

// LabelMatcher takes a matcher string and returns a span element that displays the matcher.
const LabelMatcherCard = ({ matcher, onDelete }: LabelMatcherCardProps) => {
	if(typeof matcher === "string") {
		matcher = parseMatcher(matcher);
		if(matcher === null) {
			return <span>Invalid matcher</span>;
		}
	}

	const { label, value, isRegex, isNegative } = matcher;

	let operator = "";

	if (isRegex) {
		operator = isNegative ? "!~" : "=~";
	} else {
		operator = isNegative ? "!=" : "=";
	}

	const canBeEdited = onDelete !== undefined;

	return (
		<span className={style["label-matcher"]}>
			{label} {operator} {value}
			{canBeEdited && (
				<button type="button" onClick={onDelete} class={style["delete-label-button"]}>
					🞫
				</button>
			)}
		</span>
	);
};

export default LabelMatcherCard;
