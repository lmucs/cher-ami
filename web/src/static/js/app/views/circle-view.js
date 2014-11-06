define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circle-view')

    var CircleView = marionette.ItemView.extend({
        template: template,

        ui: {

        },

        events: {
        },

        

    });

    exports.CircleView = CircleView;
})