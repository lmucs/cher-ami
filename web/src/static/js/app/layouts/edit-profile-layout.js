define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/profile-layout')
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView
    

    var EditProfileLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            editProfile: '#edit-profile-container'
        },

        ui: {
        },

        events: {
        },

        initialize: function(options) {
        },
        
        onRender: function() {
            var profile = new ProfileView();
            this.editProfile.show(editProfile);
        },
    });

    exports.EditProfileLayout = EditProfileLayout;
})