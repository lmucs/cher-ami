define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')
    var SidebarView = require('app/views/sidebar-view').SidebarView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    var CircleView = require('app/views/circle-view').CircleView;

    var CircleLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            sidebar: '#sidebar-container',
            circle: '#circle-container'
        },

        ui: {

        },

        events: {

        },

        onRender: function(options) { 
            var sidebar = new SidebarView();
            var CreateCircle = new CreateCircleView();
            var circle = new CircleView();
            this.sidebar.show(sidebar);
            this.circle.show(CreateCircle);
            // this.circle.show(circle);
        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
