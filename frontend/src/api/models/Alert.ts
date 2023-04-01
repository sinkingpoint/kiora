/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Acknowledgement } from "./Acknowledgement";

export type Alert = {
	id: string;
	labels: Record<string, string>;
	annotations: Record<string, string>;
	status: Alert.status;
	acknowledgement?: Acknowledgement;
	startsAt: string;
	endsAt?: string;
	timeoutDeadline: string;
};

export namespace Alert {
	export enum status {
		FIRING = "firing",
		ACKED = "acked",
		RESOLVED = "resolved",
		TIMED_OUT = "timed out",
		SILENCED = "silenced",
	}
}
