define(function(require, exports, module) {

    var backbone = require('backbone');

    var Signup = Backbone.Model.extend({
        /*validate: function(attrs, options) {
            if (attrs.password !== attrs.confirmpassword) {
                return "Must have same password/password confirmation";
            }
        },*/
        url: '/api/signup',
        defaults: {
            handle: null,
            email: null,
            password: null,
            confirmpassword: null
        },

        initialize: function() {

        }
    });

exports.Signup = Signup;

});
