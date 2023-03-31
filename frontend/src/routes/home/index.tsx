import { h } from "preact";
import AlertList from "../../components/alertlist";
import SingleStatPanel from "../../components/stats/single_stat_panel";
import styles from "./styles.css";

const Home = () => {
	return (
		<div>
			<div class={styles.row}>
				<SingleStatPanel title="active alerts" value="5000"/>
				<SingleStatPanel title="silenced alerts" value="5000"/>
				<SingleStatPanel title="inhibited alerts" value="5000"/>
				<SingleStatPanel title="resolved alerts" value="5000"/>
			</div>
			<AlertList />
		</div>
	);
};

export default Home;
