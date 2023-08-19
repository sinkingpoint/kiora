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

	let effectInputs: Inputs = [loader, loaded, setLoaded];
	if(inputs) {
		effectInputs = [...effectInputs, ...inputs];
	}

	useEffect(() => {
		if (!loaded) {
			loader();
			setLoaded(true);
		}
	}, effectInputs);

	return loaded ? children : <Spinner />;
};

export default Loader;
