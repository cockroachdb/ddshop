import React from "react";
import { render } from "react-dom";
import { Provider } from "react-redux";
import { createStore, compose, applyMiddleware } from "redux";
import todoApp from "./reducers/index.jsx";
import App from "./components/app.jsx";
import thunk from "redux-thunk";

// allows us to return promises from actions
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose;
let store = composeEnhancers(applyMiddleware(thunk))(createStore)(todoApp);

render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById("react-app")
);
