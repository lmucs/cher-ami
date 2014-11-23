define(function(require, exports, module){

    var marionette = require('marionette');
    var AppController = require('app/app-controller').AppController;
    var AppRouter = marionette.AppRouter.extend({
        controller: new AppController(),
        appRoutes: {
            '': 'index',
            'goToCircles': 'showCircle',
            'goToCreateCircle': 'showCreateCircle',
            'goToProfile': 'showProfile',
            'editProfile': 'showEditProfile',
        },

    });

    exports.AppRouter = AppRouter;

});