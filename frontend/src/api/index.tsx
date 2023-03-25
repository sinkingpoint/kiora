import { Alert, AlertFilter } from "./models";

interface IAPI {
	getAlerts(): Promise<Alert[]>;
}

class APIV1Impl implements IAPI {
	apiBase: string;

	constructor(apiBase: string) {
		this.apiBase = apiBase;
	}

	url(path: string): string {
		return this.apiBase + path;
	}

	getAlerts(query?: AlertFilter): Promise<Alert[]> {
		const params = new URLSearchParams();

		if (query !== undefined) {
			if (query.id !== undefined) {
				params.append("id", query.id);
			}
		}

		return fetch(this.url("/api/v1/alerts?" + params)).then((response) => response.json());
	}
}

export default new APIV1Impl("http://localhost:4278");
