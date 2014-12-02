define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')

    /** Circle views **/
    var CircleView = require('app/views/circle-view').CircleView;
    var CirclesView = require('app/views/circles-view').CirclesView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    var CreateCircle = require('app/models/create-circle').CreateCircle;
    var Circles = require('app/collections/circles').Circles;

    var CircleLayout = marionette.LayoutView.extend({
        template: template,
        circles: new Circles(),
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
            var circle = new CirclesView({
                collection: this.circles
            });
            this.circle.show(circle);
        },

        showCreateCircle: function(options) {
            var CreateCircle = new CreateCircleView({
                collection: this.circles
            });
            this.circle.show(CreateCircle)
        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
