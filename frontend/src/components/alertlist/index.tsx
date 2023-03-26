import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import { Alert } from "src/api/models";
import api from "../../api";
import Single from "../alertcard";

interface AlertViewState {
	alerts: Alert[];
	error?: string;
	loading: boolean;
}

interface ErrorViewProps {
	error: string;
}

const ErrorView = (props: ErrorViewProps) => {
	return <div>{props.error}</div>;
};

interface SuccessViewProps {
	alerts: Alert[];
}

const SuccessView = (props: SuccessViewProps) => {
	return (
		<>
			{(props.alerts.length > 0 &&
				props.alerts.map((alert) => {
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
		await api
			.getAlerts({}, { order: "DESC", orderBy: ["__starts_at__"], limit: 100, offset: 0 })
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
