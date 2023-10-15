import { h } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Alert, DefaultService } from "../../api";
import Single from "../alertcard";
import styles from "./styles.css";
import Loader from "../loader";

interface AlertViewState {
	alerts: Alert[];
	error?: string;
}

interface ErrorViewProps {
	error: string;
}

const ErrorView = ({ error }: ErrorViewProps) => {
	return <div>{error}</div>;
};

interface SuccessViewProps {
	alerts: Alert[];
}

const SuccessView = ({ alerts }: SuccessViewProps) => {
	return (
		<div class={styles.success}>
			{(alerts.length > 0 &&
				alerts.map((alert) => {
					return <Single key={alert.id} alert={alert} />;
				})) || <div>No alerts</div>}
		</div>
	);
};

const AlertList = () => {
	const [alerts, setAlerts] = useState<AlertViewState>({
		alerts: [],
	});

	const fetchAlerts = async () => {
		await DefaultService.getAlerts({ sort: ["__starts_at__"], order: "DESC", limit: 100 })
			.then((newAlerts) => {
				setAlerts({
					alerts: newAlerts,
					error: "",
				});
			})
			.catch((error) => {
				setAlerts({
					alerts: [],
					error: error.toString(),
				});
			});
	};

	let contents: JSX.Element;

	if (alerts.error) {
		contents = <ErrorView error={alerts.error} />;
	} else {
		contents = <SuccessView alerts={alerts.alerts} />;
	}

	return (
		<Loader loader={fetchAlerts} inputs={[alerts, setAlerts]}>
			{contents}
		</Loader>
	);
};

export default AlertList;
