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

        initialize: function(options) {

        },
        
        onRender: function() {
            var sidebar = new SidebarView();
            var profile = new ProfileView();
            var editProfile = new EditProfileView();
            this.sidebar.show(sidebar);
            this.profile.show(editProfile);
        }

    });

    exports.ProfileLayout = ProfileLayout;
})
