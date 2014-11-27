define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circles-view')
    var CircleView = require('app/views/circle-view').CircleView;
    var Circle = require('app/models/circle').Circle;

    var CirclesView = marionette.CollectionView.extend({
        childView: CircleView,
        childViewContainer: '#circleAreaContainer',
        template: template,

        ui: {

        },

        events: {
        },

        onSubmit: function() {
        },

        initialize: function(options) {

        }

    });

    exports.CirclesView = CirclesView;
})
