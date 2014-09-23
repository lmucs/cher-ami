define(function(require, exports, module) {
    var marionette = require('marionette')
    var template = require('hbs!../templates/signup-view')

    var SignupView = marionette.ItemView.extend({
        template: template,

        ui: {
            handle: '#handle',
            pass: '#pass1',
            confirmPass: '#pass2',
            rememberMe: '#remember-me',
            signup: '#signup'
        },

        events: {
            'click signup': 'onFormConfirm'
        },

        initialize: function(options) {

        },

        onFormConfirm: function(event) {
            console.log("Test");
        }
    });

    exports.SignupView = SignupView;
})
