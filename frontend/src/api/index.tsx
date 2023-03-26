import { Alert, AlertFilter } from "./models";

interface QueryOpts {
	order?: "ASC" | "DESC";
	orderBy?: string[];
	offset?: number;
	limit?: number;
}

interface IAPI {
	getAlerts(query?: AlertFilter, opts?: QueryOpts): Promise<Alert[]>;
}

class APIV1Impl implements IAPI {
	apiBase: string;

	constructor(apiBase: string) {
		this.apiBase = apiBase;
	}

	url(path: string): string {
		return this.apiBase + path;
	}

	getAlerts(query?: AlertFilter, opts?: QueryOpts): Promise<Alert[]> {
		const params = new URLSearchParams();

		if (query !== undefined) {
			if (query.id !== undefined) {
				params.append("id", query.id);
			}
		}

		if (opts !== undefined) {
			if (opts.orderBy !== undefined) {
				opts.orderBy.forEach((o) => params.append("sort", o));
				if (opts.order !== undefined) {
					params.append("order", opts.order);
				}
			}

			if (opts.offset !== undefined) {
				params.append("offset", opts.offset.toString());
			}

			if (opts.limit !== undefined) {
				params.append("limit", opts.limit.toString());
			}
		}

		return fetch(this.url("/api/v1/alerts?" + params)).then((response) => response.json());
	}
}

export default new APIV1Impl("http://localhost:4278") as IAPI;
