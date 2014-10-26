define(function(require, exports, module) {

    var backbone = require('backbone');

    var Session = Backbone.Model.extend({
        url: 'api/sessions',
        defaults: {
            handle: null,
            sessionid: null
        },
        initialize: function() {
            //this.load();
            $.ajaxPrefilter(function(options, originalOptions, jqXHR) {
                options.xhrFields = {
                    withCredentials: true
                };
            });
        },

        /*load: function() {
            this.model.set({
                user_id: $.cookie('handle'),
                sessionid: $.cookie('sessionid')
            })
        },*/

        // Takes in a login model
        login: function(login_model) {
            console.log(login_model);
            var that = this;
            var credentials = login_model.toJSON();
            this.save(credentials, {
                success: function(model, response) {
                    that.unset('password');
                    that.set(response);
                    // Removes response property from response object
                    that.unset('response');
                }
            })
        },

        logout: function() {
            this.destroy({
                success: function(model, resp) {
                    model.clear({silent: true})
                }
            })
        },

        authenticated: function() {
            // Add logic here which fetches auth.
            return false;
        }
    });

exports.Session = Session;

});
