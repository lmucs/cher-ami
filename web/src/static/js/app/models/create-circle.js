define(function(require, exports, module) {

    var backbone = require('backbone');

    var CreateCircle = Backbone.Model.extend({

        url: '/api/circles',
        defaults: {
            circleName: null,
            description: null,
            visibility: null
        },

        initialize: function() {

        }
    });

    exports.CreateCircle = CreateCircle;

});
