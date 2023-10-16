import { h } from "preact";
import style from "./styles.css";
import { Silence } from "../../api";
import LabelMatcherCard from "../labelmatchercard";

interface LabelViewProps {
	silence: Silence;
}

const Labels = ({ silence }: LabelViewProps) => {
	return (
		<div class={style.labels}>
			{silence.matchers.map((matcher) => {
				return <LabelMatcherCard matcher={matcher} />;
			})}
		</div>
	);
};

export default Labels;
