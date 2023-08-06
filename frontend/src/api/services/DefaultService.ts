/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Alert } from '../models/Alert';
import type { AlertAcknowledgement } from '../models/AlertAcknowledgement';
import type { Silence } from '../models/Silence';
import type { StatsResult } from '../models/StatsResult';

import type { CancelablePromise } from '../core/CancelablePromise';
import { OpenAPI } from '../core/OpenAPI';
import { request as __request } from '../core/request';

export class DefaultService {

    /**
     * Get alerts details
     * Takes an optional filter, limit, ordering, and fields and returns alerts based on those
     *
     * @returns Alert Got alerts
     * @throws ApiError
     */
    public static getAlerts({
        limit,
        offset,
        sort,
        order,
        id,
    }: {
        /**
         * The maximum number of results to return
         */
        limit?: number,
        /**
         * The offset into the results to return. Used for pagination
         */
        offset?: number,
        /**
         * The fields to sort the results by
         */
        sort?: Array<string>,
        /**
         * The order of the results. Only valid if `sort` is also specified.
         */
        order?: 'ASC' | 'DESC',
        /**
         * Get only the given alert by ID
         */
        id?: string,
    }): CancelablePromise<Array<Alert>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/alerts',
            query: {
                'limit': limit,
                'offset': offset,
                'sort': sort,
                'order': order,
                'id': id,
            },
            errors: {
                400: `Invalid query parameters`,
                500: `Backing DB failed`,
            },
        });
    }

    /**
     * Add, or update alerts
     * @returns any Alerts accepted for addition, or updating
     * @throws ApiError
     */
    public static postAlerts({
        requestBody,
    }: {
        /**
         * Alerts to add, or update in the system.
         */
        requestBody?: Array<Alert>,
    }): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/alerts',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Alerts are invalid`,
                500: `Sending the alerts to the cluster failed`,
            },
        });
    }

    /**
     * Query aggregated stats about alerts in the system
     * @returns StatsResult Sucessfully queried stats
     * @throws ApiError
     */
    public static getAlertsStats({
        type,
        args,
    }: {
        type: string,
        /**
         * The arguments to the query, depending on the query type.
         */
        args?: Record<string, string>,
    }): CancelablePromise<Array<StatsResult>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/alerts/stats',
            query: {
                'type': type,
                'args': args,
            },
            errors: {
                400: `The arguments provided were invalid for the query type`,
                500: `The underlying database failed when querying alerts`,
            },
        });
    }

    /**
     * Acknowledge an alert
     * @returns any The alert was sucessfully acknowledged
     * @throws ApiError
     */
    public static postAlertsAck({
        requestBody,
    }: {
        /**
         * Metadata when acknowledging an alert
         */
        requestBody?: AlertAcknowledgement,
    }): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/alerts/ack',
            body: requestBody,
            mediaType: 'application/json',
            errors: {
                400: `Some data was missing from the acknowledgment`,
                500: `Broadcasting the acknowledgment failed`,
            },
        });
    }

    /**
     * Get silences
     * @returns Silence Returns all the silences
     * @throws ApiError
     */
    public static getSilences(): CancelablePromise<Array<Silence>> {
        return __request(OpenAPI, {
            method: 'GET',
            url: '/silences',
        });
    }

    /**
     * Silence alerts
     * @returns Silence The silence was created
     * @throws ApiError
     */
    public static postSilences({
        requestBody,
    }: {
        /**
         * A silence to add
         */
        requestBody?: Silence,
    }): CancelablePromise<Silence> {
        return __request(OpenAPI, {
            method: 'POST',
            url: '/silences',
            body: requestBody,
            mediaType: 'application/json',
        });
    }

}
