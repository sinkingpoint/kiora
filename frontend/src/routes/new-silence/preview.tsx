import { h, Fragment } from "preact";
import Loader from "../../components/loader";
import { Alert, DefaultService } from "../../api";
import { useState } from "react";
import AlertCard from "../../components/alertcard";
import Button from "../../components/button";
import { getSilenceEnd, parseMatcher } from "./utils";

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
	}).then(() => {
		window.location.href = "/";
	});
};

const PreviewPage = ({ duration, creator, comment, matchers }: PreviewPageProps) => {
	const [alerts, setAlerts] = useState<Alert[]>([]);
	const fetchAffectedAlerts = () => {
		DefaultService.getAlerts({ matchers }).then((alerts) => {
			setAlerts(alerts);
		});
	};

	return (
		<>
			<table>
				<tr>
					<td>Duration:</td>
					<td>{duration}</td>
				</tr>

				<tr>
					<td>Creator:</td>
					<td>{creator}</td>
				</tr>

				<tr>
					<td>Comment:</td>
					<td>{comment}</td>
				</tr>
			</table>

			<div>
				<Button
					label="Create"
					onClick={() => {
						CreateSilence({ duration, creator, comment, matchers });
					}}
				/>
			</div>

			<Loader
				loader={fetchAffectedAlerts}
				done={
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
				}
			/>
		</>
	);
};

export default PreviewPage;
