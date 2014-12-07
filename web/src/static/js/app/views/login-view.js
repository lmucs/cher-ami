define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/login-view');

    var Login = require('app/models/login').Login;
    var inputValidator = require('app/utils/input-validator').InputValidator;

    var LoginView = marionette.ItemView.extend({
        template: template,

        tagName: "div",
        className: "loginContainer",
        ui: {
            handle: '#handle',
            pass: '#pass1',
            rememberMe: '#remember-me',
            login: '#login',
            inputForm: '#loginForm',
            response: '#response'
        },

        events: {
            'keyup #handle': 'inputValidate',
            'click #remember-me': 'onRememberConfirm',
            'click #login': 'onLogin',
        },

        initialize: function(options) {
            this.model = new Login({
                session: options.session
            });
            this.session = options.session;
        },

        onRememberConfirm: function(event) {
            // Session-request method goes here
        },


        onLogin: function(event) {
            if(!this.model.validate()) {
                this.ui.response.html("Your username or password was incorrect.");
            }
            console.log("Logging in");
            event.preventDefault();
            this.model.set("handle", this.ui.handle.val());
            this.model.set("password", this.ui.pass.val());
            this.model.authenticate();
            this.model.clear();
        },

        inputValidate: function(event) {
            inputValidator(this.ui.inputForm)
        }

    });

    exports.LoginView = LoginView;
})
