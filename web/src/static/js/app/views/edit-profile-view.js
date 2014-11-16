define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/edit-profile-view')

    var EditProfileView = marionette.ItemView.extend({

        template: template,

        ui: {
            firstname: '#first-name',
            lastname: '#last-name',
            gender: '#gender-selector',
            birthday: '#birthday',
            location: '#location',
            bio: '#bio',
            interests: '#interests',
            languages: '#languages'
        },

        events: {
            // Search event
        },

        initialize: function(options) {

        }

    });

    exports.EditProfileView = EditProfileView;
})
