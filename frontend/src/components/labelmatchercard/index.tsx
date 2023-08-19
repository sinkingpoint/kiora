import { h } from "preact";
import style from "./style.css";
import { Matcher } from "../../api";

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

export interface LabelMatcherCardProps {
    matcher: string;
    onDelete: () => void;
}

// LabelMatcher takes a matcher string and returns a span element that displays the matcher.
const LabelMatcherCard = ({matcher, onDelete}: LabelMatcherCardProps) => {
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

export default LabelMatcherCard;