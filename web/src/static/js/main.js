define(function(require, exports, module){

    var marionette = require('marionette');
    var renderer = require('app/renderer');
    var app = require('app/app');
    var AppRouter = require('app/app-router').AppRouter;

    app.appRouter = new AppRouter();
    app.start();

});