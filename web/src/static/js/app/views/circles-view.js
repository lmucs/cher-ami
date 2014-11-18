define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/circle-view')


    var CirclesView = marionette.CollectionView.extend({
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
