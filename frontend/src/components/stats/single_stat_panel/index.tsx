import { h } from "preact";
import styles from "./styles.css";

interface SingleStatPanelProps {
    title: string;
    value: string;
    color?: string;
}

export default ({title, value, color}: SingleStatPanelProps) => {
    if(color === undefined) {
        color = "#fff";
    }

    return <div class={styles.card}>
        <div class={styles.value} style={{color: color}}>
            {value}
        </div>

        <label class={styles.title}>
            {title}
        </label>
    </div>;
}