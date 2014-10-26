define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/profile-layout')

    var ProfileLayout = marionette.LayoutView.extend({
        template: template,

        regions: {

        },

        initialize: function(options) {

        }

    });

    exports.ProfileLayout = ProfileLayout;
})
