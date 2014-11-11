define(function(require, exports, module) {
    var marionette = require('marionette');
    var template = require('hbs!../templates/sidebar-view');
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;

    var SidebarView = marionette.LayoutView.extend({
        template: template,

        regions: {
            circle: '#create-circle-view'
        },

        ui: {
            createCircle: '#goToCreateCircle'
        },

        events: {
            'click #goToCreateCircle': 'showCreateCircle'
        },

        initialize: function(options) {
        },

        showCreateCircle: function(options) {
            var createCircle = new CreateCircleView();
            this.circle.show(createCircle);
        }

    });

    exports.SidebarView = SidebarView;
})