define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/profile-layout')
    var ProfileView = require('app/views/profile-view').ProfileView    

    var ProfileLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            profile: '#profile-container',           
        },

        ui: {            
            profileSaveButton: '#profileSaveButton',           
        },

        events: {
            
            'click #profileSaveButton': 'onRender'
        },

        initialize: function(options) {
        },
        
        onRender: function() {
            var profile = new ProfileView();
            this.profile.show(profile);
        },
    });
    exports.ProfileLayout = ProfileLayout;
})
