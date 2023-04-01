import { h } from "preact";
import { useEffect, useState } from "preact/hooks";
import LabelList from "../../components/alertcard/labels";
import { Alert, DefaultService } from "../../api";
import style from "./styles.css";

interface AlertState {
	loading: boolean;
	alert?: Alert;
	error?: string;
}

interface AlertProps {
	id: string;
}

interface SuccessViewProps {
	alert: Alert;
}

const SuccessView = ({ alert }: SuccessViewProps) => {
	const startTime = new Date(alert.startsAt);
	const endTime = new Date(alert.endsAt);

	return (
		<div class={style["alert-view"]}>
			<span class={style["alert-row"]}>
				<h1>{alert.labels["alertname"] || <i>No Alert Name</i>}</h1>
			</span>

			{alert.acknowledgement !== undefined && (
				<span>
					<span class={style["alert-row"]}>Acknowledged by {alert.acknowledgement.creator}</span>
				</span>
			)}

			<span class={style["alert-row"]}>
				<LabelList alert={alert} />
			</span>

			<span class={style["alert-row"]}>
				<label>Status:</label> {alert.status}
			</span>

			<span class={style["alert-row"]}>
				<label>ID:</label> {alert.id}
			</span>

			<span class={style["alert-row"]}>
				<label>Started At:</label> {startTime.toLocaleString()}
			</span>

			{endTime.getTime() > 0 && (
				<span class={style["alert-row"]}>
					<label>Ended At:</label> {endTime.toLocaleString()}
				</span>
			)}

			<span>
				<h3>Annotations:</h3>
			</span>
			{Object.keys(alert.annotations).map((key) => {
				return (
					<span key={key}>
						<label>{key}:</label> {alert.annotations[key]}
					</span>
				);
			})}
		</div>
	);
};

const AlertView = ({ id }: AlertProps) => {
	const [state, setState] = useState<AlertState>({
		loading: true,
	});

	useEffect(() => {
		if (!state.loading) {
			return;
		}

		DefaultService.getAlerts(null, null, null, null, id)
			.then((alerts) => {
				if (alerts.length === 0) {
					return;
				}

				setState({
					loading: false,
					alert: alerts[0],
				});
			})
			.catch((error) => {
				setState({
					loading: false,
					error: error.toString(),
				});
			});
	}, [state, id]);

	if (state.loading) {
		return <div>Loading...</div>;
	} else if (state.error) {
		return <div>{state.error}</div>;
	} else if (state.alert) {
		return <SuccessView alert={state.alert} />;
	}
};

export default AlertView;
