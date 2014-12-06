define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/header-view')
    var Session = require('app/models/session').Session;

    var HeaderView = marionette.ItemView.extend({

        template: template,

        ui: {
            // Search box
            // Profile button
            logoutButton: "#logout"
        },

        events: {
            'click #logout': 'onLogout'
        },

        initialize: function(options) {
            console.log('Initializing header view', options);
            this.session = options.session;
        },

        onLogout: function(event) {
            event.preventDefault();
            this.session.logout({
                reload: true
            });
            window.location.replace('/');
        }

    });

    exports.HeaderView = HeaderView;
})