define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/login-view');

    var LoginView = marionette.ItemView.extend({

        template: template,
        tagName: "div",
        className: "loginContainer",
        ui: {
            email: '#input-email',
            pass: '#pass1',
            rememberMe: '#remember-me',
            login: '#login'
        },

        events: {
            'click #remember-me': 'onRememberConfirm',
            'click #login': 'onLogin'
        },

        initialize: function(options) {

        },

        onRememberConfirm: function() {
            // Session-request method goes here
        },

        onLogin: function(argument) {
           // needs to be implemented later.
        }

    });

    exports.LoginView = LoginView;
})