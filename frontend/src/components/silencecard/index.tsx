import { h } from "preact";
import { Silence } from "../../api";
import style from "./styles.css";
import Labels from "./labels";
import { formatDate, formatDuration } from "../../utils/date";

interface SilenceCardProps {
    silence: Silence;
}

const SilenceCard = ({silence}: SilenceCardProps) => {
    const startDate = new Date(silence.startsAt);
    const endDate = new Date(silence.endsAt);

    return <a href={`/silences/${silence.id}`} class={style["alert-link"]}>
        <div class={style.single}>
            <div>
                <div class={style["single-top"]} style={{marginBottom: "5px"}}>
                    {silence.id} created by {silence.creator}
                </div>

                <div class={style["single-top"]}>
                    {startDate > new Date() ? "Starts" : "Started"} at {formatDate(startDate)}
                </div>

                <div class={style["single-top"]}>
                    {endDate > new Date() ? "Ends" : "Ended"} at {formatDate(endDate)}
                </div>

                <div class={style["single-top"]}>
                    <Labels silence={silence} />
                </div>
            </div>
        </div>
    </a>;
};

export default SilenceCard;
