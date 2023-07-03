import { h } from "preact";
import styles from "./styles.css";

interface SingleStatPanelProps {
	title: string;
	value: string;
	color?: string;
}

const SingleStatPanel = ({ title, value, color }: SingleStatPanelProps) => {
	return (
		<div class={styles.card}>
			<div class={styles.value}>
				{value}
			</div>

			<label class={styles.title}>{title}</label>
		</div>
	);
};

export default SingleStatPanel;
