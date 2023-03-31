import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Alert, DefaultService } from "../../api";
import Single from "../alertcard";

interface AlertViewState {
	alerts: Alert[];
	error?: string;
	loading: boolean;
}

interface ErrorViewProps {
	error: string;
}

const ErrorView = ({error}: ErrorViewProps) => {
	return <div>{error}</div>;
};

interface SuccessViewProps {
	alerts: Alert[];
}

const SuccessView = ({alerts}: SuccessViewProps) => {
	return (
		<>
			{(alerts.length > 0 &&
				alerts.map((alert) => {
					return <Single alert={alert} />;
				})) || <div>No alerts</div>}
		</>
	);
};

const AlertList = () => {
	const [alerts, setAlerts] = useState<AlertViewState>({
		alerts: [],
		loading: true,
	});

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

	useEffect(() => {
		if (alerts.loading) {
			fetchAlerts();
		}
	}, [alerts]);

	if (alerts.loading) {
		return <div>Loading...</div>;
	} else if (alerts.error) {
		return <ErrorView error={alerts.error} />;
	} else {
		return <SuccessView alerts={alerts.alerts} />;
	}
};

export default AlertList;
