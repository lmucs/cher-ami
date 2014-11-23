define(function(require, exports, module) {

    var backbone = require('backbone');

    var Circle = require('app/models/circle').Circle;

    var Circles = Backbone.Collection.extend({
        url: 'api/circles',
        model: Circle,

        parse: function(response) {
            // return JSON.parse(response.Objects);
        },

        initialize: function() {

        }
    });

    exports.Circles = Circles;

});
