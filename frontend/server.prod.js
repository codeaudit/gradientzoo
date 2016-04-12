var express = require('express');
var compression = require('compression');
var httpProxy = require('http-proxy');
var basicAuth = require('basic-auth-connect');

var app = new (express)()
var port = 3000

app.use(compression())

app.use('/static', express.static('dist'))

var proxy = httpProxy.createProxyServer({})
proxy.on('error', function(err, req, res) {
  console.log('PROXY ERROR: ' + err)
})

app.use('/api', function(req, res) {
  proxy.web(req, res, {
    target: 'http://' + process.env.GRADIENTZOO_API_SVC_SERVICE_HOST + ':' + process.env.GRADIENTZOO_API_SVC_SERVICE_PORT
  })
})

app.use(basicAuth('ericflo', 'poop'))

app.use(function(req, res) {
  res.sendFile(__dirname + '/index.html')
})

app.listen(port, function(error) {
  if (error) {
    console.error(error)
  } else {
    console.info("==> 🌎  Listening on port %s. Open up http://localhost:%s/ in your browser.", port, port)
  }
})