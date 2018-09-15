const todo = (state = {}, action) => {
  switch (action.type) {
    case "ADD_TODO":
      return {
        id: action.id,
        title: action.text,
        completed: false
      };
    // Commented out so that the UI doesn't show a change
    // until it's been acknowledged by the server.
    // case "TOGGLE_TODO":
    //   if (state.id !== action.todo.id) {
    //     return state;
    //   }
    //   return Object.assign({}, state, {
    //     completed: !state.completed
    //   });

    default:
      return state;
  }
};

const todos = (state = [], action) => {
  switch (action.type) {
    case "ADD_TODO":
      return [...state, todo(undefined, action)];
    // Commented out so that the UI doesn't show a change
    // until it's been acknowledged by the server.
    // case "TOGGLE_TODO":
    //   console.log("toggle todo");
    //   return state.map(t => todo(t, action));
    case "ATTEMPT_DELETE_TODO":
      return state.map((todo) => (
        todo.id === action.todoID
          ? Object.assign({}, todo, { deleting: true })
          : todo
      ));
    case "DELETE_TODO_ERR":
      return state.map((todo) => (
        todo.id === action.todoID
          ? Object.assign({}, todo, { deleting: false })
          : todo
      ));
    case "RECEIVE_TODO":
      return [...state, action.todo];
    case "RECEIVE_TODOS":
      return action.todos;
    case "RECEIVE_DELETE_TODO":
      return state.filter(function(item) {
        return item.id != action.id;
      });
    case "RECEIVE_UPDATE_TODO": {
      const newState = state.filter((item) => {
        if (item.id === action.todo.id) {
          item.completed = action.todo.completed;
        }
        return true;
      });
      return newState;
    }
    default:
      return state;
  }
};

export default todos;
