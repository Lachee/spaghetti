const webpack               = require('webpack');
const path                  = require('path');



// webpack.config.js
module.exports = (env, options) => {
  return {
    entry: './src/js/spaghetti.js',
    output: {
      filename:       'spaghetti.js',
      path:           path.resolve(__dirname, './resources/bin'),
      library:        'Spaghetti',
      libraryTarget:  'umd',
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
          use: {
            loader: options.mode === 'production' ? 'arraybuffer-loader' : 'file-loader',
            options: {
              name: '[name].[ext]',
              emitFile: false,
            }
          }
        },
        {
          test: /\.css$/i,
          use: ["style-loader", "css-loader"],
        }
      ],
    },
  };
};