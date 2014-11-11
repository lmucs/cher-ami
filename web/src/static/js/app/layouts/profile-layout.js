define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/profile-layout')
    var SidebarView = require('app/views/sidebar-view').SidebarView;
    var ProfileView = require('app/views/profile-view').ProfileView
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView

    var ProfileLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            sidebar: '#sidebar-container',
            profile: '#profile-container'
        },

        ui: {
            editProfile: '#editProfile',
            profileSaveButton: '#profileSaveButton'
        },

        events: {
            'click #editProfile': 'showEditProfile',
            'click #profileSaveButton': 'onRender'
        },

        initialize: function(options) {

        },
        
        onRender: function() {
            var sidebar = new SidebarView();
            var profile = new ProfileView();
            this.sidebar.show(sidebar);
            this.profile.show(profile);
        },

        showEditProfile: function(options) {
            var editProfile = new EditProfileView();
            this.profile.show(editProfile);
        }

    });

    exports.ProfileLayout = ProfileLayout;
})
