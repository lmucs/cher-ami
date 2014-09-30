define(function(require, exports, module) {
    var marionette = require('marionette')
    var template = require('hbs!../templates/signup-view')

    var Signup = require('app/models/signup').Signup;

    var SignupView = marionette.ItemView.extend({
        template: template,

        //take the div marionette creates and give it a class named mainContainer.
        tagName: "div",
        className: "mainContainer",
        ui: {
            handle: '#handle',
            email: '#input-email',
            pass: '#pass1',
            confirmPass: '#pass2',
            rememberMe: '#remember-me',
            signup: '#signup'
        },

        events: {
            'click #remember-me': 'onRememberConfirm',
            'click #signup': 'onFormConfirm'
        },

        initialize: function(options) {

        },

        onRememberConfirm: function(options) {
            // Session-request method goes here
        },

        onFormConfirm: function(event) {
            event.preventDefault();
            var req = new Signup({
                handle: this.ui.handle.val(),
                email: this.ui.email.val(),
                password: this.ui.pass.val(),
                confirmpassword: this.ui.confirmPass.val()
            });
            console.log(req)
            req.save();
        }
    });

    exports.SignupView = SignupView;
})
