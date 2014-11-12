define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/profile-layout')
    var ProfileView = require('app/views/profile-view').ProfileView
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView
    

    var ProfileLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            profile: '#profile-container'
        },

        ui: {
            editProfile: '#editProfile',
            profileSaveButton: '#profileSaveButton',
            
        },

        events: {
            'click #editProfile': 'showEditProfile',
            'click #profileSaveButton': 'onRender',
            
        },

        initialize: function(options) {

        },
        
        onRender: function() {
            var profile = new ProfileView();
            this.profile.show(profile);
        },

        showEditProfile: function(options) {
            var editProfile = new EditProfileView();
            this.profile.show(editProfile);
        }

    });

    exports.ProfileLayout = ProfileLayout;
})
