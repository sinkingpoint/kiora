import { h } from "preact";
import { Alert } from "../../api";
import Labels from "./labels";
import style from "./styles.css";

interface CardProps {
	alert: Alert;
}

const AlertCard = ({ alert }: CardProps) => {
	const startTime = new Date(alert.startsAt).toLocaleString();

	return (
		<a href={`/alerts/${alert.id}`} class={style["alert-link"]}>
			<div class={style.single}>
				<div>
					<span class={style["single-top"]}>{startTime}</span>
					<span class={style["single-top"]}>
						{(alert.labels["alertname"] && `alertname="${alert.labels["alertname"]}"' `) || (
							<i>No Alert Name</i>
						)}
					</span>

					<div class={style["single-top"]}>
						<Labels alert={alert} />
					</div>
				</div>
			</div>
		</a>
	);
};

export default AlertCard;
