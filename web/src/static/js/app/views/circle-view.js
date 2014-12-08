define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circle-view')
    var Circle = require('app/models/circle').Circle;

    var CircleView = marionette.ItemView.extend({
        model: Circle,
        template: template,

        ui: {

        },

        events: {

        },

        onRender: function(options) {
        },

        initialize: function(options) {
            this.model = options.model
        }

    });

    exports.CircleView = CircleView;
})