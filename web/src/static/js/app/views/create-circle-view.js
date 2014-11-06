define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/create-circle-view')

    var CreateCircleView = marionette.ItemView.extend({
        template: template,

        ui: {

        },

        events: {
        },

        

    });

    exports.CreateCircleView = CreateCircleView;
})