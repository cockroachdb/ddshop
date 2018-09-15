import React, { PropTypes } from "react";

const Todo = ({ onClick, deleteTodo, updateTodo, todo }) => (
  <li>
    <div className="view">
      <input
        className="toggle"
        type="checkbox"
        checked={todo.completed}
        onClick={evt => {
          evt.preventDefault();
          todo.completed = !todo.completed;
          updateTodo();
        }}
      />
      <label>{todo.title}</label>
      <button
        className="destroy"
        onClick={evt => {
          evt.preventDefault();
          deleteTodo();
        }}
      />

      {/* <span
        style={{
          textDecoration: todo.completed? 'line-through': 'none'
        }}
        onClick={onClick}
      >
        {todo.title}
      </span>

      <a
        href="#"
        style={{
          marginLeft: "10px",
          textDecoration: todo.completed? 'line-through': 'none'
        }}
        onClick={e => {
          e.preventDefault()
          todo.completed = !todo.completed
          updateTodo()
        }}
      >
          (complete)
      </a>

      <a
        href="#"
        style={{
          marginLeft: "10px"
        }}
        onClick={e => {
            e.preventDefault()
            deleteTodo()
        }}
      >
        (delete)
      </a> */}
    </div>
  </li>
);

Todo.propTypes = {
  onClick: PropTypes.func.isRequired,
  deleteTodo: PropTypes.func.isRequired,
  updateTodo: PropTypes.func.isRequired,
  todo: PropTypes.object.isRequired
};

export default Todo;
