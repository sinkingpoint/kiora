import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Alert } from "../..//api/models";
import LabelList from "../../components/alertcard/labels";
import api from "../../api";
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

const SuccessView = (props: SuccessViewProps) => {
	const startTime = new Date(props.alert.startsAt);
	const endTime = new Date(props.alert.endsAt);

	return (
		<div class={style["alert-view"]}>
			<span class={style["alert-row"]}>
				<h1>{props.alert.labels["alertname"] || <i>No Alert Name</i>}</h1>
			</span>

			{props.alert.acknowledgement !== undefined && (
				<span>
					<span class={style["alert-row"]}>
						Acknowledged by {props.alert.acknowledgement.creator}
					</span>
				</span>
			)}

			<span class={style["alert-row"]}>
				<LabelList alert={props.alert} />
			</span>

			<span class={style["alert-row"]}>
				<label>Status:</label> {props.alert.status}
			</span>

			<span class={style["alert-row"]}>
				<label>ID:</label> {props.alert.id}
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
			{Object.keys(props.alert.annotations).map((key) => {
				return (
					<span>
						<label>{key}:</label> {props.alert.annotations[key]}
					</span>
				);
			})}
		</div>
	);
};

export default ({ id }: AlertProps) => {
	const [state, setState] = useState<AlertState>({
		loading: true,
	});

	useEffect(() => {
		if (!state.loading) {
			return;
		}

		api
			.getAlerts({ id: id })
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
	}, [state]);

	if (state.loading) {
		return <div>Loading...</div>;
	} else if (state.error) {
		return <div>{state.error}</div>;
	} else if (state.alert) {
		return <SuccessView alert={state.alert} />;
	}

	return <div></div>;
};
