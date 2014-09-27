define(function(require, exports, module) {
    var marionette = require('marionette')
    var template = require('hbs!../templates/signup-view')

    var SignupView = marionette.ItemView.extend({
        template: template,

        //take the div marionette creates and give it a class named mainContainer.
        tagName: "div",
        className: "mainContainer",
        ui: {
            handle: '#handle',
            pass: '#pass1',
            confirmPass: '#pass2',
            rememberMe: '#remember-me',
            signup: '#signup'
        },

        events: {
            'click #remember-me': 'onFormConfirm'
        },

        initialize: function(options) {

        },

        onFormConfirm: function() {
            alert(this.ui.handle.val() + "  hola la");
        }
    });

    exports.SignupView = SignupView;
})
