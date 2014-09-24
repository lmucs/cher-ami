define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('templates/footer-view')

    var FooterView = marionette.ItemView.extend({

        template: template,

        ui: {
        },

        events: {
        },

        initialize: function(options) {

        }

    });

    exports.FooterView = FooterView;
})
