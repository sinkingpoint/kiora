import { h } from "preact";
import { Alert } from "src/api/models";
import Label from "../label";
import style from "./styles.css";

interface SingleProps {
	alert: Alert;
}

export default ({ alert }: SingleProps) => {
	const startTime = new Date(Date.parse(alert.startsAt)).toLocaleString();

	return (
		<a href={`/alerts/${alert.id}`} class={style["alert-link"]}>
			<div class={style.single}>
				<div>
					<span class={style["single-top"]}>{startTime}</span>
					<span class={style["single-top"]}>
						{(alert.labels["alertname"] && 'alertname="' + alert.labels["alertname"] + '"') || (
							<i>No Alert Name</i>
						)}
					</span>
				</div>
				<div class={style.labels}>
					{Object.keys(alert.labels).map((key) => {
						if (key === "alertname") {
							return;
						}
						return <Label labelName={key} labelValue={alert.labels[key]} />;
					})}
				</div>
			</div>
		</a>
	);
};
