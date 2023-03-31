import { h } from "preact";
import Label from "../labelcard";
import style from "./styles.css";
import { Alert } from "../../api";

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