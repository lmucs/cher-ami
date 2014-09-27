define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/header-view')

    var HeaderView = marionette.ItemView.extend({

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

    exports.HeaderView = HeaderView;
})
