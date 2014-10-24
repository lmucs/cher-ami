define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/profile-view')

    var ProfileView = marionette.ItemView.extend({

        template: template,

        ui: {
            // Search box
            // Profile button
        },

        events: {
            // Search event
        },

        initialize: function(options) {

        }

    });

    exports.ProfileView = ProfileView;
})
