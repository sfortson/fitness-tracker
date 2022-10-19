/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-env node, es6 */

const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

const extensions = ['.ts', '.tsx', '.js'];
// const tsConfigFile = path.join(__dirname, 'tsconfig.json');

module.exports = {
  entry: path.resolve(__dirname, './src/index.tsx'),
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        use: ['babel-loader'],
      },
    ],
  },
  resolve: {
    extensions,
  },
  output: {
    path: path.resolve(__dirname, './dist'),
    filename: 'bundle.js',
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new HtmlWebpackPlugin({
      favicon: path.resolve(__dirname, 'static/favicon.ico'),
      title: 'Fitness Tracker',
      template: path.resolve(__dirname, 'src/index.html'),
      hash: true,
    }),
  ],
  devServer: {
    static: path.resolve(__dirname, './dist'),
    hot: true,
    historyApiFallback: true,
  },
};
