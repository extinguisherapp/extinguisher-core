var postcss = require('postcss');
var fs= require("fs")

//postcss([ require('precss'), require('autoprefixer'), require('cssnano') ])
var scss = require('postcss-scss');

//postcss([ require('precss'),require('autoprefixer'), require('cssnext') ])
postcss([ require('postcss-cssnext'), require('precss'), require('autoprefixer'), require('postcss-browser-reporter')])
.process(fs.readFileSync('css/1.scss'), { from: 'css/1.scss', to: '1.css',  map: { inline: false }})
.then(function (result) {
    console.log(result.css)
    /*
    fs.writeFileSync('1.css', result.css);
    if ( result.map ) fs.writeFileSync('1.css.map', result.map);
    */
},
     function(error){
         console.log(error)
     })

