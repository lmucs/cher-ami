define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')
    /** Circle views **/
    var CircleView = require('app/views/circle-view').CircleView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    var Circle = require('app/models/circle').Circle;

    var CircleLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            circle: '#circle-container'
        },

        ui: {

        },

        events: {
            'click #goToCreateCircles': 'showCreateCircle',
            'click #circleAddMemberButton': 'onRender'
        },

        onRender: function(options) {
            var circle = new CircleView();           
            this.circle.show(circle);
        },

        showCreateCircle: function(options) {
            var CreateCircle = new CreateCircleView();
            this.circle.show(CreateCircle)
        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
