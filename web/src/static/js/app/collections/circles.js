define(function(require, exports, module) {

    var backbone = require('backbone');

    var CreateCircle = require('app/models/create-circle').CreateCircle;

    var Circles = Backbone.Collection.extend({
        url: 'api/circles',
        model: CreateCircle,

        parse: function(response) {
            return response.results;
        },

        initialize: function() {

        }
    });

    exports.Circles = Circles;

});
