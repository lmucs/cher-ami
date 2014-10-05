define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/login-view');

    var Login = require('app/models/login').Login;

    var LoginView = marionette.ItemView.extend({

        template: template,
        tagName: "div",
        className: "loginContainer",
        ui: {
            handle: '#handle',
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

        onLogin: function(event) {
            event.preventDefault();
            var req = new Login();
            req.save({
                handle: this.ui.handle.val(),
                password: this.ui.pass.val()
            },
            {
                success: function() {
                    alert("It worked!");
                },
                error: function() {
                    alert(":(");
                }
            });
        }

    });

    exports.LoginView = LoginView;
})
