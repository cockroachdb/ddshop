import React from "react";
import FilterLink from "../containers/filterlink.jsx";

const Footer = () => (
  <footer className="footer">
    <span className="todo-count">{5} items left</span>
    <ul className="filters">
      <li>
        <a href="#">All</a>
      </li>
      <li>
        <a href="#">Active</a>
      </li>
      <li>
        <a href="#">Completed</a>
      </li>
    </ul>
  </footer>
);

export default Footer;
