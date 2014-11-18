define(function(require, exports, module) {

    var backbone = require('backbone');

    var Circle = Backbone.Model.extend({
        url: 'api/circles',
        defaults: {
            Content: null,
            Id: null
        },

        initialize: function() {

        }
    });

    exports.Circle = Circle;

});
