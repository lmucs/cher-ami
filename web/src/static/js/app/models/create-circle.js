define(function(require, exports, module) {

    var backbone = require('backbone');

    var CreateCircle = Backbone.Model.extend({

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

    exports.CreateCircle = CreateCircle;

});
