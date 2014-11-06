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
            containerArea: '#containerArea',
        },

        events: {
            'click #showSignup': 'showForm'
        },

        showForm: function(options) {
            var signupView = new SignupView();
            this.ui.containerArea.html(signupView.el);
            console.log(signupView.el);
            signupView.render();
        },

        initialize: function(options) {

        }

    });

    exports.LandingLayout = LandingLayout;
})
