import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Alert, DefaultService } from "../../api";
import Single from "../alertcard";
import styles from "./styles.css";

interface AlertViewState {
	alerts: Alert[];
	error?: string;
	loading: boolean;
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
		loading: true,
	});

	useEffect(() => {
		const fetchAlerts = async () => {
			await DefaultService.getAlerts(100, 0, ["__starts_at__"], "DESC")
				.then((newAlerts) => {
					setAlerts({
						...alerts,
						alerts: newAlerts,
						loading: false,
					});
				})
				.catch((error) => {
					setAlerts({
						...alerts,
						error: error.toString(),
						loading: false,
					});
				});
		};

		if (alerts.loading) {
			fetchAlerts();
		}
	}, [alerts]);

	if (alerts.loading) {
		return <div>Loading...</div>;
	} else if (alerts.error) {
		return <ErrorView error={alerts.error} />;
	}

	return <SuccessView alerts={alerts.alerts} />;
};

export default AlertList;
