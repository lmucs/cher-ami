define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/landing-layout')

    var LandingLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            signup: '#signUpContainer',
            login: '#loginContainer'
        },

        initialize: function(options) {

        }

    });

    exports.LandingLayout = LandingLayout;
})
