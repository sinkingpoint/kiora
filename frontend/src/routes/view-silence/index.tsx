import { h, Fragment } from "preact";
import { DefaultService, Silence } from "../../api";
import { useState } from "preact/hooks";
import Loader from "../../components/loader";
import { formatDate, formatDuration } from "../../utils/date";

export interface ViewSilenceProps {
	id: string;
}

export interface ViewSilenceState {
	silence?: Silence;
	error?: string;
}

const silenceView = (silence: Silence) => {
	const startTime = new Date(silence.startsAt);
	let secondsSinceStart = Math.floor((Date.now() - startTime.getTime()) / 1000);
	const hasntStartedYet = secondsSinceStart < 0;
	if (hasntStartedYet) {
		secondsSinceStart *= -1;
	}

	const endTime = new Date(silence.endsAt);
	let secondsTillEnd = Math.floor((endTime.getTime() - Date.now()) / 1000);
	const hasEnded = secondsTillEnd < 0;
	if (hasEnded) {
		secondsTillEnd *= -1;
	}

	const startString = hasntStartedYet
		? `In ${formatDuration(secondsSinceStart)}`
		: `${formatDuration(secondsSinceStart)} Ago`;
	const endString = hasEnded
		? `${formatDuration(secondsTillEnd)} Ago`
		: `In ${formatDuration(secondsTillEnd)}`;

	return (
		<div>
			<h1>Silence {silence.id}</h1>
			<div>
				<h2>Creator</h2>
				<p>{silence.creator}</p>
			</div>

			<div>
				<h2>Comment</h2>
				<p>{silence.comment}</p>
			</div>

			<div>
				<h2>Started at</h2>
				<p>
					{formatDate(startTime)} ({startString})
				</p>
			</div>

			<div>
				<h2>Ends at</h2>
				<p>
					{formatDate(endTime)} ({endString})
				</p>
			</div>
		</div>
	);
};

const errorView = (error: string) => {
	return (
		<div>
			<h1>Error</h1>
			<p>{error}</p>
		</div>
	);
};

const ViewSilence = ({ id }: ViewSilenceProps) => {
	const [silence, setSilence] = useState<ViewSilenceState>({});

	const fetchSilence = () => {
		DefaultService.getSilences({ matchers: [`__id__=${id}`] }).then((response) => {
			if (response.length == 0) {
				setSilence({
					error: "Silence not found",
				});

				return;
			}

			setSilence({
				silence: response[0],
			});
		});
	};

	return (
		<div>
			<Loader loader={fetchSilence} inputs={[setSilence]}>
				<>{silence?.silence ? silenceView(silence.silence) : errorView(silence.error)}</>
			</Loader>
		</div>
	);
};

export default ViewSilence;
