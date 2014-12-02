define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circle-view')
    var CreateCircle = require('app/models/create-circle').CreateCircle;

    var CircleView = marionette.ItemView.extend({
        model: CreateCircle,
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