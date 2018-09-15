import React from "react";
import { connect } from "react-redux";
import Footer from "../components/footer.jsx";

function mapStateToProps(state) {
  return {
    count: state.todos.filter((todo) => !todo.completed).length,
  };
}

export default connect(mapStateToProps)(Footer);
