import { h } from "preact";
import style from "./styles.css";

const Header = () => {
	return (
		<header>
			<a href="/" class={style.logo}>
				<h1>Kiora</h1>
			</a>
		</header>
	);
};

export default Header;
