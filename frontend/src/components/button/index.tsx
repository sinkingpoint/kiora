import { h } from "preact";
import style from "./styles.css";

interface ButtonProps {
	label: string;
	onClick?: () => void;
}

const Button = ({ label, onClick }: ButtonProps) => {
	return (
		<button onClick={onClick} class={style["button"]}>
			{label}
		</button>
	);
};

export default Button;
