define(function(require, exports, module) {

    var backbone = require('backbone');

    var CreateCircle = Backbone.Model.extend({

        url: '/api/circles',
        defaults: {
            name: null,
            description: null,
            visibility: null
        },

        initialize: function() {

        }
    });

    exports.CreateCircle = CreateCircle;

});
