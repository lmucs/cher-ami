define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = Backbone.Model.extend({
        defaults: {
            Content: null,
            Id: null
        },

        initialize: function() {

        }
    });

    exports.Message = Message;

});
