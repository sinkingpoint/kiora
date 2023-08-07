import "./style";
import App from "./components/app";
import { OpenAPI } from "./api";

let apiHost = "";
// When running in dev, we're running in `preact-cli` which allows us to set
// environment variables because it's a full Javascript runtime. When running
// in a full Kiora deployment, we're running embedded in the Go app, which can't set them.
// The reason we have to do this try-catch fudge is that when preact sets `process.env.PREACT_APP_API_HOST`
// it _doesn't_ set `process` or `process.env`, so checking if those exist will throw an error.
try {
	apiHost = process.env.PREACT_APP_API_HOST;
} catch {}

OpenAPI.VERSION = "v1";
OpenAPI.BASE = `${apiHost}/api/${OpenAPI.VERSION}`;

export default App;
