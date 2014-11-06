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

        ui: {
            showSingup: '#showSignup',
            signup: '#signup',
            login: '#login'
        },

        events: {
            'click #showSignup': 'showForm'
        },

        showForm: function(options) {
            var newSignup = new SignupView();
            var temp = $(signup).html()
            // newSignup.signup.add(newSignup.el);
            this.ui.login.html(newSignup.el);
            console.log(temp);
            console.log(newSignup.el);
            newSignup.render();
        },

        initialize: function(options) {

        }

    });

    exports.LandingLayout = LandingLayout;
})
