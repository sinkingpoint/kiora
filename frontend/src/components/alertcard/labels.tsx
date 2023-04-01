import { h } from "preact";
import Label from "../labelcard";
import style from "./styles.css";
import { Alert } from "../../api";

interface LabelViewProps {
	alert: Alert;
}

const Labels = ({ alert }: LabelViewProps) => {
	return (
		<div class={style.labels}>
			{Object.keys(alert.labels).map((key) => {
				if (key === "alertname") {
					return;
				}
				return <Label key={key} labelName={key} labelValue={alert.labels[key]} />;
			})}
		</div>
	);
};

export default Labels;
