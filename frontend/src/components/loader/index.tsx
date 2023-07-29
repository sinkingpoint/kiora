import { useEffect, useState } from "preact/hooks";
import Spinner from "../spinner";
import { h } from "preact";

interface LoaderProps {
	loader: () => void;
	done: JSX.Element;
}

const Loader = ({ loader, done }: LoaderProps) => {
	const [loaded, setLoaded] = useState(false);

	useEffect(() => {
		if (!loaded) {
			loader();
			setLoaded(true);
		}
	}, [loader, loaded, setLoaded]);

	return loaded ? done : <Spinner />;
};

export default Loader;
