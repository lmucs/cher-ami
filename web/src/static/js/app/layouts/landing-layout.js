define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/landing-layout')
    var SignupView = require('app/views/signup-view').SignupView;
    var LoginView = require('app/views/login-view').LoginView;
    
    var LandingLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            signup: '#signup',
            login: '#login'
        },

        initialize: function(options) {

        }

    });

    exports.LandingLayout = LandingLayout;
})
