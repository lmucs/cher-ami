define(function(require, exports, module) {

    var backbone = require('backbone');

    var Session = Backbone.Model.extend({
        url: '',
        initialize: function() {

        },

        login: function(credentials) {

        },

        logout: function() {

        },

        getAuthentication: function(callback) {
            this.fetch({
                success: callback
            });
        }
    });

exports.Session = Session;

});
