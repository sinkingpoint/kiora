/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { AlertAcknowledgement } from './AlertAcknowledgement';

export type Alert = {
    readonly id: string;
    labels: Record<string, string>;
    annotations: Record<string, string>;
    status: Alert.status;
    readonly acknowledgement?: AlertAcknowledgement;
    startsAt: string;
    endsAt?: string;
    readonly timeoutDeadline: string;
};

export namespace Alert {

    export enum status {
        FIRING = 'firing',
        ACKED = 'acked',
        RESOLVED = 'resolved',
        TIMED_OUT = 'timed out',
        SILENCED = 'silenced',
    }


}

