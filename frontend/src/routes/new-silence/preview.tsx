import { h, Fragment } from "preact";
import Loader from "../../components/loader";
import { Alert, DefaultService } from "../../api";
import { useState } from "react";
import AlertCard from "../../components/alertcard";
import Button from "../../components/button";
import { getSilenceEnd } from "./utils";
import { formatDate } from "../../utils/date";
import LabelMatcherCard, { parseMatcher } from "../../components/labelmatchercard";

const MaxAlertsToDisplay = 20;

export interface PreviewPageProps {
	duration: string;
	matchers: string[];
	creator: string;
	comment: string;
}

const CreateSilence = ({ duration, creator, comment, matchers }: PreviewPageProps) => {
	const startsAt = new Date().toISOString();
	const endsAt = getSilenceEnd(duration).toISOString();

	const modelMatchers = matchers.map((matcher) => parseMatcher(matcher));

	DefaultService.postSilences({
		requestBody: {
			id: "",
			startsAt,
			endsAt,
			matchers: modelMatchers,
			creator,
			comment,
		},
	}).then((response) => {
		window.location.href = `/silences/${response.id}`;
	});
};

const PreviewPage = ({ duration, creator, comment, matchers }: PreviewPageProps) => {
	const [alerts, setAlerts] = useState<Alert[]>([]);
	const fetchAffectedAlerts = () => {
		DefaultService.getAlerts({ matchers }).then((alerts) => {
			setAlerts(alerts);
		});
	};

	const endDate = getSilenceEnd(duration);
	const end =
		endDate !== null ? (
			<span>
				Ends at{" "}
				{formatDate(endDate)}
			</span>
		) : (
			<span>Invalid duration</span>
		);
	
	const filterSpans = matchers.map((filter, i) => {
		return <LabelMatcherCard matcher={filter} />
	});

	return (
		<>
			<div>
				<label>Duration</label>
			</div>

			<div style={{ justifyContent: "space-between", flexWrap: "wrap" }}>
				{duration} {end}
			</div>

			<div>
				<label>Matchers</label>
			</div>
			<div style={{ flexWrap: "wrap" }}>{filterSpans}</div>

			<div>
				<label>Creator</label>
			</div>
			<div>{creator}</div>

			<div>
				<label>Comment</label>
			</div>
			<div>{creator}</div>

			<div>
				<Button
					label="Create"
					onClick={() => {
						CreateSilence({ duration, creator, comment, matchers });
					}}
				/>
			</div>

			<Loader loader={fetchAffectedAlerts}>
					<>
						<div>
							<h2>
								{alerts.length} affected alert{alerts.length != 1 && "s"}
							</h2>
						</div>
						<div style={{ display: "flex", flexDirection: "column" }}>
							{alerts.slice(0, MaxAlertsToDisplay).map((alert) => (
								<AlertCard key={alert.id} alert={alert} />
							))}
							{alerts.length > MaxAlertsToDisplay && (
								<div>...and {alerts.length - MaxAlertsToDisplay} more</div>
							)}
						</div>
					</>
			</Loader>
		</>
	);
};

export default PreviewPage;
