const webpack               = require('webpack');
const path                  = require('path');

// webpack.config.js
module.exports = {
  entry: './src/js/spaghetti.js',
  output: {
    filename:       'spaghetti.js',
    path:           path.resolve(__dirname, './bin'),
    library:        'Spaghetti'
  },
  resolve: { 
    fallback: { 
      os: false,
      fs: false, 
      crypto: false, 
      util: false 
    } 
  },
  module: {
    rules: [
      //{
      //  test: /(\.|_)exec\.js$/i,
      //  use: 'raw-loader',
      //},
      {
        test: /\.m?js$/,
        exclude: /node_modules/,
        use:  {
          loader: 'babel-loader', 
          options: {                
            presets: ['@babel/preset-env'],
            plugins: [
              "@babel/plugin-proposal-class-properties",
              "@babel/plugin-proposal-private-methods",
              '@babel/plugin-transform-runtime'
            ]
          },
        }
      },
      {
        test: /\.wasm$/i,
        use: 'arraybuffer-loader',
      },
      {
        test: /\.css$/i,
        use: ["style-loader", "css-loader"],
      }
    ],
  },
};