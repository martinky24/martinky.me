
// load the things we need
var path = require('path');
var express = require('express');
var app = express();

var bodyParser = require('body-parser');
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

// set the view engine to ejs
app.set('view engine', 'ejs');

app.use(express.static(__dirname + '/public'));
app.use('/public',  express.static(__dirname + '/public'));



/***********************************
 * 
 * Page routes
 *
 **********************************/

app.get('/', function(req, res) {
    var content = [];

    res.render('pages/resume', content);
});

app.get('/resume', function(req, res) {
    var content = [];

    res.render('pages/resume', content);
});

/******************
 * Error pages
 ******************/

app.use(function (req, res) {
    res.status(404);
    res.render('404');
});

app.use(function (err, req, res, next) {
    console.error(err.stack);
    res.status(500);
    res.render('500');
});



/******************
 * Launch communication
 ******************/

// app.listen(8080);
// console.log('http://localhost:8080/');

const port = process.env.PORT || 3000;
app.listen(port);
console.log('http://localhost:3000/');

module.exports = app;