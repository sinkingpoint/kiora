import { h } from "preact";
import style from "./style.css";

interface LabelProps {
	labelName: string;
	labelValue: string;
}

const LabelCard = ({ labelName, labelValue }: LabelProps) => {
	return (
		<span class={style.label}>
			{labelName}="{labelValue}"
		</span>
	);
};

export default LabelCard;
