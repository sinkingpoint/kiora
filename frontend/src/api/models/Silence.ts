/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */

import type { Matcher } from './Matcher';

export type Silence = {
    readonly id: string;
    creator: string;
    comment: string;
    startsAt: string;
    endsAt: string;
    matchers: Array<Matcher>;
};

