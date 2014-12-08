define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/profile-view')
    var Profile = require('app/models/profile').Profile;

    var ProfileView = marionette.ItemView.extend({
        model: Profile,
        template: template,

        ui: {
            // Search box
            // Profile button
        },

        events: {
            // Search event
        },

        initialize: function(options) {
            this.model = options.model
            this.session = options.session;
        },

    });

    exports.ProfileView = ProfileView;
})
