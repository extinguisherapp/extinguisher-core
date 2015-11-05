var webpack = require("webpack");
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var ManifestPlugin = require('webpack-manifest-plugin');

var contextPath = __dirname + '/assets/javascripts';

module.exports = {
    context: contextPath,
    entry: {
        home: './home.js',
        dashboard: './dashboard.js',

        //用來確保所有的圖片跟 vendors js 都會被打上 digest
        static_resource: './static_resources.js',

    },
    output: {
        filename: '[name].[hash].js',
        path: __dirname + '/public/assets',
        publicPath: "/assets/"
    },
    module: {
        loaders: [
            { test: /vendors/,   exclude: /node_modules/, loader: 'file-loader?name=[path][name].[hash].[ext]'},
            { test: /\.jsx?$/, exclude: /(node_modules|vendors)/, loaders: ['react-hot', 'babel'] },
            { test: /\.css$/,  loader: ExtractTextPlugin.extract("style-loader", "css-loader") },
            { test: /\.png$/,  loader: "url-loader?limit=100000&name=[name].[hash].[ext]" },
            { test: /\.jpg$/,  loader: "file-loader?name=[name].[hash].[ext]" },
            { test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                loader: "url-loader?limit=10000&minetype=application/font-woff&name=[name].[hash].[ext]",
            },
            { test: /\.(otf|ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
                loader: "file-loader?name=[name].[hash].[ext]",
            },
            {   test: /\.scss$/,   loader: ExtractTextPlugin.extract(
                    // activate source maps via loader query
                    'css?sourceMap!' +
                    'sass?sourceMap'
                ) }
        ],

        noParse: [
          /[\/\\]vendors[\/\\].*\.js$/
        ]
    },
    resolve: {
        extensions: ['', '.js', '.jsx'],
        modulesDirectories: ["node_modules", "javascripts"],
    },
    externals: {
        "react": 'React',
        "react/addons": "React",
        "jquery": 'window.jQuery',
        semantic: 'semantic'
    },
    devtool: "sourcemap",
    plugins: [
        new ExtractTextPlugin("[name].[hash].css", {
            disable: false,
            allChunks: true
        }),
        new webpack.ProvidePlugin({
            $: "jquery",
            jQuery: "jquery",
            "window.jQuery": "jquery",
            React: 'react'
        }),

        // 生成 manifest.json, imageExtensions 這邊用來處理 static resource
        new ManifestPlugin({
            imageExtensions: /^(css|jpe?g|png|gif|svg|woff|woff2|otf|ttf|eot|svg|js)(\.|$)/i
        })

    ]
};
