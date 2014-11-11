define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/landing-layout')
    var SignupView = require('app/views/signup-view').SignupView;
    var LoginView = require('app/views/login-view').LoginView;
    var Login = require('app/models/login').Login;

    var LandingLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            container: '#containerArea'
        },

        ui: {
            existingUser: '#goToLogin',
            newUser: '#goToSignup'
        },

        events: {
            'click #goToLogin': 'showLoginForm',
            'click #goToSignup': 'onRender'
        },

        onRender: function(options) {
            var signupView = new SignupView();
            this.container.show(signupView);
        },

        showLoginForm: function(options) {
            var loginView = new LoginView({
                session: this.session
            });

            this.container.show(loginView);
        },

        initialize: function(options) {
            this.session = options.session;
        }

    });

    exports.LandingLayout = LandingLayout;
})
