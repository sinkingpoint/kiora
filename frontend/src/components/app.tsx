import { h } from "preact";
import Router, { Route } from "preact-router";
import Home from "../routes/home";
import Header from "./header";
import Alert from "../routes/alert";
import NewSilence from "../routes/new-silence";
import ViewSilence from "../routes/view-silence";

const App = () => {
	return (
		<div id="app">
			<Header />
			<main>
				<Router>
					<Route path="/" component={Home} />
					<Route path="/alerts/:id" component={Alert} />
					<Route path="/silences/new" component={NewSilence} />
					<Route path="/silences/:id" component={ViewSilence} />
				</Router>
			</main>
		</div>
	);
};

export default App;
