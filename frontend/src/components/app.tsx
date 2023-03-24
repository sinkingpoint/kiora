import { h } from "preact";
import Router, { Route } from "preact-router";
import Home from "../routes/home";
import Header from "./header";
import Alert from "../routes/alert";

const App = () => {
	return (
		<div id="app">
			<Header />
			<main>
				<Router>
					<Route path="/" component={Home} />
					<Route path="/alerts/:id" component={Alert} />
				</Router>
			</main>
		</div>
	);
};

export default App;
