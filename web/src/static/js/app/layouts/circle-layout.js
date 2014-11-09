define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')
    var SidebarView = require('app/views/sidebar-view').SidebarView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;

    var CircleLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            sidebarContainer: '#sidebar-container',
            circleContainer: '#circle-container'
        },

        ui: {
            sidebarContainer: '#sidebar-container',
            circleContainer: '#circle-container',
            showContent: '#showContent'
        },

        events: {
            'click #showContent': 'showContent'
        },

        showContent: function(options) { 
            var sidebar = new SidebarView();
            var circle = new CreateCircleView();
            this.ui.sidebarContainer.html(sidebar.el);
            this.ui.circleContainer.html(circle.el);
            sidebar.render();
            circle.render();
        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
