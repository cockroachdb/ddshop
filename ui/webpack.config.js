var webpack = require("webpack");
var path = require("path");
var CopyWebpackPlugin = require("copy-webpack-plugin");
var CleanWebpackPlugin = require("clean-webpack-plugin");

var production = process.env.NODE_ENV === "production";

var BUILD_DIR = path.resolve(__dirname, "build");
var APP_DIR = path.resolve(__dirname, "src");

var config = {
  entry: path.join(APP_DIR + "/index.jsx"),
  output: {
    path: BUILD_DIR,
    filename: "bundle.js"
  },
  contentBase: "build/",
  publicPath: "/",
  plugins: [
    new CopyWebpackPlugin([{ from: "./index.html" }]),
    new CopyWebpackPlugin([{ from: "./node_modules/todomvc-app-css/index.css" }]),
    new CleanWebpackPlugin([BUILD_DIR])
  ],
  target: "web",
  module: {
    loaders: [
      {
        test: /\jsx?/,
        include: APP_DIR,
        loader: "babel",
        query: {
          presets: ["react", "es2015"]
        }
      },
      {
        test: /\less$/,
        loader: "style!css!less"
      }
    ]
  },
  devtool: production ? false : "source-map",
  debug: !production,
  devServer: {
    proxy: {
      "/api**": {
        target: "http://localhost:26256",
        secure: false
      }
    },
    host: "0.0.0.0",
    port: "3000",
    watchOptions: {
      poll: true
    }
  }
};

module.exports = config;
