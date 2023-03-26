import { h } from "preact";
import { Alert } from "src/api/models";
import Label from "../label";
import Labels from "./labels";
import style from "./styles.css";

interface CardProps {
	alert: Alert;
}

export default ({ alert }: CardProps) => {
	const startTime = new Date(alert.startsAt).toLocaleString();

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
				
				<Labels alert={alert} />
			</div>
		</a>
	);
};
