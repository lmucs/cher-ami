define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/landing-layout')
    var SignupView = require('app/views/signup-view').SignupView;
    var LoginView = require('app/views/login-view').LoginView;
    
    var LandingLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            container: '#containerArea'
        },

        ui: {
            showSignup: '#showSignup',
            showLogin: '#showLogin',
            containerArea: '#containerArea',
        },

        events: {
            'click #showSignup': 'showSignupForm',
            'click #showLogin' : 'showLoginForm'
        },

        showSignupForm: function(options) {
            var signupView = new SignupView();
            $("#showSignup").html('show login');
            $("#showSignup").attr('id', 'showLogin')
            this.ui.containerArea.html(signupView.el);
            console.log(signupView.el);
            signupView.render();
        },

        showLoginForm: function(options) {
            var loginView = new LoginView();
            $("#showLogin").html('show sign up');
            $("#showLogin").attr('id', 'showSignup')
            this.ui.containerArea.html(loginView.el);
            console.log(loginView.el);
            loginView.render();
        },

        initialize: function(options) {

        }

    });

    exports.LandingLayout = LandingLayout;
})
