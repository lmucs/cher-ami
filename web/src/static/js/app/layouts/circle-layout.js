define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')

    var CircleLayout = marionette.LayoutView.extend({
        template: template,

        regions: {

        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
