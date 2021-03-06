import React, { PropTypes } from "react";

const Todo = ({ onClick, deleteTodo, updateTodo, todo }) => (
  <li className={todo.completed ? "completed" : ""}>
    <div className="view">
      <input
        className="toggle"
        type="checkbox"
        checked={todo.completed}
        onChange={() => {
          const newTodo = Object.assign({}, todo, {
            completed: !todo.completed,
          });
          updateTodo(newTodo);
        }}
      />
      <label>{todo.title}</label>
      <button
        className={`destroy ${todo.deleting ? "spin" : ""}`}
        onClick={() => {
          deleteTodo();
        }}
      />
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
