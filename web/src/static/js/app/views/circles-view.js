define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circles-view')
    var CircleView = require('app/views/circle-view').CircleView;
    var Circle = require('app/models/circle').Circle;
    var CreateCircle = require('app/models/create-circle').CreateCircle;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;

    var CirclesView = marionette.CompositeView.extend({
        childView: CircleView,
        childViewContainer: '#circleAreaContainer',
        template: template,

        ui: {

        },

        events: {
        },

        onRender: function() {
            console.log(this.collection.length)

        },

        initialize: function(options) {            
            this.collection = options.collection;
            this.session = options.session;
            // this.collection.fetch({
            //     success: function(res) {
            //         console.log(res);
            //     }
            // });
        }

    });

    exports.CirclesView = CirclesView;
})
