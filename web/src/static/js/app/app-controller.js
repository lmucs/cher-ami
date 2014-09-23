define(function(require, exports, module) {

    var backbone = require('backbone');
    var marionette = require('marionette');
    var app = require('app/app');

    var AppController = marionette.Controller.extend({

        initialize: function(options) {
            this.app = app;
        
            // Initialization of views will go here.
        },

        // Needed for AppRouter to initialize index route.
        index: function() {

        }

    });

    exports.AppController = AppController;

});