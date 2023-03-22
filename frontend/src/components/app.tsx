import { h } from "preact";
import Router, { Route } from "preact-router";
import Home from "../routes/home";

const App = () => {
  return (
    <div id="app">
      <Router>
        <Route path="/" component={Home} />
      </Router>
    </div>
  );
};

export default App;
