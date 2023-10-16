import { h } from "preact";
import style from "./styles.css";
import { Link } from "preact-router/match";

const Header = () => {
	return (
		<header class={style.header}>
			<a href="/" class={style.logo}>
				<h1>Kiora</h1>
			</a>

			<nav>
				<Link activeClassName={style.active} href="/silences">
					Silences
				</Link>
			</nav>
		</header>
	);
};

export default Header;
