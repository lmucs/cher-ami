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

        authenticate: function() {
            this.save({}, {
                success: function(model, response) {
                    window.location.replace('/#Home');
                    window.location.reload()
                }
            })
        },

        initialize: function() {

        }
    });

    exports.Signup = Signup;

});
