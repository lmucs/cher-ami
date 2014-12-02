define(function(require, exports, module){

    var marionette = require('marionette');
    var AppController = require('app/app-controller').AppController;
    var AppRouter = marionette.AppRouter.extend({
        controller: new AppController(),
        appRoutes: {
            '': 'index',
            'Home': 'showHomeLayout',
            'Circles': 'showCircle',
            // 'CreateCircle': 'showCreateCircle',
            'Profile': 'showProfile',
            'Settings': 'showSettings'
        },

    });

    exports.AppRouter = AppRouter;

});