import { h } from "preact";
import { useState } from "preact/hooks";
import style from "./styles.css";
import CreatePage from "./create";
import PreviewPage, { PreviewPageProps } from "./preview";

const NewSilence = () => {
	const [preview, setPreview] = useState<PreviewPageProps>(null);

	let page: JSX.Element;
	if (preview === null) {
		page = <CreatePage onPreview={setPreview} />;
	} else {
		page = <PreviewPage {...preview} />;
	}

	return (
		<div class={style["silence-form"]}>
			<h1>New Silence</h1>

			{page}
		</div>
	);
};

export default NewSilence;
