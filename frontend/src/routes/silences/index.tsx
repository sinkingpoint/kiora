import { h } from "preact";
import Loader from "../../components/loader";
import { useState } from "preact/hooks";
import { DefaultService, Silence } from "../../api";
import SilenceCard from "../../components/silencecard";

const silencesView = (silences: Silence[]) => {
	if (silences.length === 0) {
		return (
			<div>
				<p>No silences found</p>
			</div>
		);
	}

	return (
		<div>
			{silences.map((silence) => {return <SilenceCard silence={silence} />})}
		</div>
	);
};

const errorView = (error: string) => {
	return (
		<div>
			<p>{error}</p>
		</div>
	);
};

interface AllSilencesViewState {
	silences?: Silence[];
	error?: string;
}

const AllSilencesView = () => {
	const [silences, setSilences] = useState<AllSilencesViewState>({});

	const fetchSilences = () => {
		DefaultService.getSilences({ limit: 10, sort: ["__starts_at__"], order: "DESC" })
			.then((response) => {
				setSilences({
					silences: response,
					error: "",
				});
			})
			.catch((error) => {
				setSilences({
					silences: [],
					error: error,
				});
			});
	};

	return (
		<div>
			<Loader loader={fetchSilences}>
				<div>{silences.silences ? silencesView(silences.silences) : errorView(silences.error)}</div>
			</Loader>
		</div>
	);
};

export default AllSilencesView;
