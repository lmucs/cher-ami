define(function(require, exports, module) {

    var backbone = require('backbone');

    var Circle = Backbone.Model.extend({
        url: 'api/circles',
        defaults: {
            name: null,
            url: null,
            description: null,
            owner: null,
            visibility: null,
            members: null,
            created: null
        },

        initialize: function() {

        }
    });

    exports.Circle = Circle;

});
