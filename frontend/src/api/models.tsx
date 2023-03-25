type AlertStatus = "firing" | "silenced" | "acked" | "resolved" | "timed out";

export interface AlertAcknowledgement {
	creator?: string;
	comment?: string;
}

export interface Alert {
	id?: string;
	labels: { [key: string]: string };
	annotations: { [key: string]: string };
	status: AlertStatus;
	startsAt: string;
	endsAt: string;
	timeOutDeadline: string;
	acknowledgement?: AlertAcknowledgement;
}

export interface AlertFilter {
	id?: string;
}
