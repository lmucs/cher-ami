define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = Backbone.Model.extend({
        defaults: {
            content: null,
            id: null
        },

        initialize: function() {

        }
    });

    exports.Message = Message;

});
