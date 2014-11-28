define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/sidebar')

    var SidebarView = marionette.ItemView.extend({
        template: template,

        ui: {
            home: '#goToHome',
            profile: '#goToProfile',
            circles: '#goToCircles',
            createCircle: '#goToCreateCircle',
            settings: '#goToSettings',
        },

        events: {

        },

        initialize: function(options) {
            // this.model = options.model
        }

    });

    exports.SidebarView = SidebarView;
})