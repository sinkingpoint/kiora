import { h, Fragment } from "preact";
import { useEffect, useState } from "preact/hooks";
import { DefaultService } from "../../api";
import AlertList from "../../components/alertlist";
import SingleStatPanel from "../../components/stats/single_stat_panel";
import styles from "./styles.css";

interface StatsState {
	firingAlerts: number;
	silencedAlerts: number;
	ackedAlerts: number;
	resolvedAlerts: number;
	timedOutAlerts: number;
	loading: boolean;
	error?: string;
}

// StatsRow is a component that displays a row of stats about the alerts in the system, breaking down alerts by their state.
const StatsRow = () => {
	const [stats, setStats] = useState<StatsState>({
		firingAlerts: 0,
		silencedAlerts: 0,
		ackedAlerts: 0,
		resolvedAlerts: 0,
		timedOutAlerts: 0,
		loading: true,
	});

	const fetchStats = async () => {
		await DefaultService.getAlertsStats("status_count", {})
			.then((result) => {
				const newStats = {
					...stats,
					loading: false,
				};

				result.forEach((stat) => {
					if (stat.labels.status === "firing") {
						newStats.firingAlerts = stat.frames[0][0];
					} else if (stat.labels.status === "silenced") {
						newStats.silencedAlerts = stat.frames[0][0];
					} else if (stat.labels.status === "acked") {
						newStats.ackedAlerts = stat.frames[0][0];
					} else if (stat.labels.status === "resolved") {
						newStats.resolvedAlerts = stat.frames[0][0];
					} else if (stat.labels.status === "timed out") {
						newStats.timedOutAlerts = stat.frames[0][0];
					}
				});

				setStats(newStats);
			})
			.catch((error) => {
				setStats({
					...stats,
					error: error.toString(),
					loading: false,
				});
			});
	};

	useEffect(() => {
		if (stats.loading) {
			fetchStats();
		}
	});

	if (stats.loading) {
		return <div>Loading...</div>;
	}

	if (stats.error) {
		return <div>{stats.error}</div>;
	}

	return (
		<div class={styles.row}>
			<SingleStatPanel title="Firing Alerts" value={stats.firingAlerts.toString()} />
			<SingleStatPanel title="Silenced Alerts" value={stats.silencedAlerts.toString()} />
			<SingleStatPanel title="Acked Alerts" value={stats.ackedAlerts.toString()} />
			<SingleStatPanel title="Resolved Alerts" value={stats.resolvedAlerts.toString()} />
			<SingleStatPanel title="Timed Out Alerts" value={stats.timedOutAlerts.toString()} />
		</div>
	);
};

const Home = () => {
	return (
		<>
			<StatsRow />
			<div class={styles.row}><AlertList /></div>
		</>
	);
};

export default Home;
