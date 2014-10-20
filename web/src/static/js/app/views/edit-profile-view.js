define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/edit-profile-view')

    var EditProfileView = marionette.ItemView.extend({

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

    exports.EditProfileView = EditProfileView;
})
