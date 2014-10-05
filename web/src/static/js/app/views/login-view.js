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
            this.model = new Login();
        },

        onRememberConfirm: function() {
            // Session-request method goes here
        },

        onLogin: function(event) {
            event.preventDefault();
            var data = {
                handle: this.ui.handle.val(),
                password: this.ui.pass.val()
            }
            var callbacks = {
                success: function() {
                    alert("It worked!");
                },
                error: function() {
                    alert(":(");
                }
            }
            this.model.save(data, callbacks);
        }

    });

    exports.LoginView = LoginView;
})
