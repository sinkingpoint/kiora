type AlertStatus = "firing" | "silenced" | "acked" | "resolved" | "timed out";

interface AlertAcknowledgement {
  creator?: string;
  comment?: string;
}

interface Alert {
  id?: string;
  labels: { [key: string]: string };
  annotations: { [key: string]: string };
  status: AlertStatus;
  startsAt: string;
  endsAt: string;
  timeOutDeadline: string;
  acknowledgement?: AlertAcknowledgement;
}
