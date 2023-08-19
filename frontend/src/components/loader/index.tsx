import { Inputs, useEffect, useState } from "preact/hooks";
import Spinner from "../spinner";
import { h } from "preact";

interface LoaderProps {
	loader: () => void;
	inputs?: Inputs;
	children?: JSX.Element;
}

const Loader = ({ loader, inputs, children }: LoaderProps) => {
	const [loaded, setLoaded] = useState(false);

	if(!inputs) {
		inputs = [];
	}

	useEffect(() => {
		if (!loaded) {
			loader();
			setLoaded(true);
		}
	}, [loader, loaded, setLoaded, ...inputs]);

	return loaded ? children : <Spinner />;
};

export default Loader;
