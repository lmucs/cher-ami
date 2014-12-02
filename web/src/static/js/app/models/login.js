define(function(require, exports, module) {

    var backbone = require('backbone');

    var Login = Backbone.Model.extend({
        url: 'api/sessions',
        defaults: {
            handle: null,
            password: null
        },

        initialize: function(options) {
            this.session = options.session;
        },

        authenticate: function() {
            var session = this.session;
            var that = this;
            // http://bit.ly/1wJmFiY
            this.save({}, {
                success: function(model, response) {
                    that.unset('password');
                    console.log(response);
                    that.set(response);
                    // Removes response property from response object
                    that.unset('response');
                    session.set('token', response.token);
                    session.set('handle', response.handle);
                    console.log(session.toJSON())
                }
            })
        }
    });

    exports.Login = Login;

});
