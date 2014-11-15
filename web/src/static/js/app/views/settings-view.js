define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/settings-view')

    var SettingsView = marionette.ItemView.extend({

        template: template,

        ui: {
        },

        events: {
        },

        initialize: function(options) {
        }

    });

    exports.SettingsView = SettingsView;
})
