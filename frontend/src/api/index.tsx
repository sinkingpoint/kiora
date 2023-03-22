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

  getAlerts(): Promise<Alert[]> {
    return fetch(this.url("/api/v1/alerts")).then((response) =>
      response.json()
    );
  }
}

export default new APIV1Impl("localhost:4278");
