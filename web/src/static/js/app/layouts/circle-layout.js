define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/circle-layout')
    
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    var CircleView = require('app/views/circle-view').CircleView;

    var CircleLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
           
            circle: '#circle-container'
        },

        ui: {

        },

        events: {

        },

        onRender: function(options) { 
           
            var CreateCircle = new CreateCircleView();
            var circle = new CircleView();           
            this.circle.show(CreateCircle);
        },

        initialize: function(options) {

        }

    });

    exports.CircleLayout = CircleLayout;
})
