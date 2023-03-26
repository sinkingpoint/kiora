import { h } from "preact";
import { Alert } from "src/api/models";
import Label from "../label";
import style from "./styles.css";

interface LabelViewProps {
    alert: Alert;
}

export default ({alert}: LabelViewProps) => {
    return <div class={style.labels}>
    {Object.keys(alert.labels).map((key) => {
        if (key === "alertname") {
            return;
        }
        return <Label labelName={key} labelValue={alert.labels[key]} />;
    })}
    </div>
}