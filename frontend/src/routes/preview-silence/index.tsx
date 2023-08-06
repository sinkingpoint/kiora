import { h } from "preact";
import style from "./styles.css";
import { DefaultService } from "../../api";

interface PreviewSilenceProps {
	labelMatchers: string[];
	comment: string;
	duration: string;
	creator: string;
}

const PreviewSilence = (props: PreviewSilenceProps) => {
	return <div class={style["silence-form"]}></div>;
};

export default PreviewSilence;
