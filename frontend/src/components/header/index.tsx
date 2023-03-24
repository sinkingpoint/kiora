import { h } from "preact";
import style from "./styles.css";

export default () => {
  return (
    <header>
      <a href="/" class={style.logo}>
        <h1>Kiora</h1>
      </a>
    </header>
  );
};