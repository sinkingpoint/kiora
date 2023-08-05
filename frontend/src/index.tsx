import "./style";
import App from "./components/app";
import { OpenAPI } from "./api";

OpenAPI.VERSION = "v1";
OpenAPI.BASE = `${process.env.PREACT_APP_API_HOST}/api/${OpenAPI.VERSION}`;

export default App;
