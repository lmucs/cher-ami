define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/home-layout')

    var HomeLayout = marionette.LayoutView.extend({
        template: template,

        regions: {

        },

        initialize: function(options) {

        }

    });

    exports.HomeLayout = HomeLayout;
})
