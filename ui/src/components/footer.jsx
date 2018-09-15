import React from "react";
// import FilterLink from "../containers/filterlink.jsx";

const Footer = ({ count }) => (
  <footer className="footer">
    <span className="todo-count">{count} {count === 1 ? "item" : "items"} left</span>
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
