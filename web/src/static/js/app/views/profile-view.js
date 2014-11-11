define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/profile-view')
    var ProfileView = require('app/views/profile-view').ProfileView
    var EditProfileView = require('app.views/edit-profile-view').EditProfileView

    var ProfileView = marionette.ItemView.extend({

        template: template,
        
        regions: {
            sidebar: '#sidebar-container',
            feed: '#feed-container'
        },

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
