define(function(require, exports, module) {

    var backbone = require('backbone');

<<<<<<< HEAD
    var Cricle = require('app/models/circle').Cricle;
=======
    var Circle = require('app/models/circle').Circle;
>>>>>>> master

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
