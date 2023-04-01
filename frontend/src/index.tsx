import "./style";
import App from "./components/app";
import { OpenAPI } from "./api";

OpenAPI.VERSION = "v1";
OpenAPI.BASE = `http://localhost:4278/api/${OpenAPI.VERSION}`;

export default App;
