import { h } from "preact";
import Router, { Route } from "preact-router";
import Home from "../routes/home";
import Header from "./header";

const App = () => {
  return (
    <div id="app">
      <Header />
      <main>
        <Router>
          <Route path="/" component={Home} />
        </Router>
      </main>
    </div>
  );
};

export default App;
