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
          test: /resources/i,
          exclude: /bin/,
          loader: path.resolve('./webpack.loader.js'),
          options: { 
            embed: function(url, mimeType, context) {
              return options.mode === 'production'; //mimeType.startsWith('image/');
            }
          }
        },
        {
          test: /\.wasm$/i,
          use: {
            loader: options.mode === 'production' ? 'arraybuffer-loader' : 'file-loader',
            options: {
              publicPath:   '.',
              name:         '[path][name].[ext]',
              emitFile:     false,
            }
          }
        },
        {
          test: /\.css$/i,
          use: ["style-loader", "css-loader"],
        },
      ],
    },
  };
};




//        {
//          test: /\.(png|jpg|gif)$/i,
//          exclude: /bin/,
//          loader: options.mode === 'production' ? 'url-loader' : 'file-loader',
//          options: {
//            publicPath:   '.',
//            name:         '[path][name].[ext]',
//            emitFile:     false,
//          }
//        },
//        // Text Loader for text is while a smart idea, not actually too useful
//        // {
//        //   test: /\.(glsl|vert|frag|txt|shader|json)/i,
//        //   exclude: /bin/,
//        //   loader:   options.mode === 'production' ? 'text-loader' : 'file-loader',
//        //   options:  options.mode === 'production' ? {} : {
//        //     publicPath:   '.',
//        //     name:         '[path][name].[ext]',
//        //     emitFile:     false,
//        //   }
//        // },
//        {
//          test: /resources/i,
//          exclude: /bin|png|jpg|gif/,
//           /** This can be either a url-loader or a arraybuffer-loader */
//           loader:   options.mode === 'production' ? 'url-loader' : 'file-loader',  
//           options:  options.mode === 'production' ? { limit: true } : {
//             publicPath:   '.',
//             name:         '[path][name].[ext]',
//             emitFile:     false,
//           }
//         },
         