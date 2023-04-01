/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type { Alert } from "../models/Alert";
import type { StatsResult } from "../models/StatsResult";

import type { CancelablePromise } from "../core/CancelablePromise";
import { OpenAPI } from "../core/OpenAPI";
import { request as __request } from "../core/request";

export class DefaultService {
	/**
	 * Get alerts details
	 * Takes an optional filter, limit, ordering, and fields and returns alerts based on those
	 *
	 * @param limit The maximum number of results to return
	 * @param offset The offset into the results to return. Used for pagination
	 * @param sort The fields to sort the results by
	 * @param order The order of the results. Only valid if `sort` is also specified.
	 * @param id Get only the given alert by ID
	 * @returns Alert Got alerts
	 * @throws ApiError
	 */
	public static getAlerts(
		limit?: number,
		offset?: number,
		sort?: Array<string>,
		order?: "ASC" | "DESC",
		id?: string
	): CancelablePromise<Array<Alert>> {
		return __request(OpenAPI, {
			method: "GET",
			url: "/alerts",
			query: {
				limit: limit,
				offset: offset,
				sort: sort,
				order: order,
				id: id,
			},
			errors: {
				400: `Invalid query parameters`,
				500: `Backing DB failed`,
			},
		});
	}

	/**
	 * Add, or update alerts
	 * @param requestBody Alerts to add, or update in the system.
	 * @returns any Alerts accepted for addition, or updating
	 * @throws ApiError
	 */
	public static postAlerts(
		requestBody?: Array<{
			labels: Record<string, string>;
			annotations?: Record<string, string>;
			startsAt?: string;
			endsAt?: string;
		}>
	): CancelablePromise<any> {
		return __request(OpenAPI, {
			method: "POST",
			url: "/alerts",
			body: requestBody,
			mediaType: "application/json",
			errors: {
				400: `Alerts are invalid`,
				500: `Sending the alerts to the cluster failed`,
			},
		});
	}

	/**
	 * Query aggregated stats about alerts in the system
	 * @param type
	 * @param args The arguments to the query, depending on the query type.
	 * @returns StatsResult Sucessfully queried stats
	 * @throws ApiError
	 */
	public static getAlertsStats(type: string, args: any): CancelablePromise<Array<StatsResult>> {
		return __request(OpenAPI, {
			method: "GET",
			url: "/alerts/stats",
			query: {
				type: type,
				args: args,
			},
			errors: {
				400: `The arguments provided were invalid for the query type`,
				500: `The underlying database failed when querying alerts`,
			},
		});
	}

	/**
	 * Acknowledge an alert
	 * @param requestBody Metadata when acknowledging an alert
	 * @returns any The alert was sucessfully acknowledged
	 * @throws ApiError
	 */
	public static postAlertsAck(requestBody?: {
		alertID?: string;
		creator?: string;
		comment?: string;
	}): CancelablePromise<any> {
		return __request(OpenAPI, {
			method: "POST",
			url: "/alerts/ack",
			body: requestBody,
			mediaType: "application/json",
			errors: {
				400: `Some data was missing from the acknowledgment`,
				500: `Broadcasting the acknowledgment failed`,
			},
		});
	}
}
