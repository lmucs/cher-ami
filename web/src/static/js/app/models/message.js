define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = Backbone.Model.extend({
        url: 'api/messages',
        defaults: {
            Content: null,
            Id: null
        },

        initialize: function() {

        }
    });

    exports.Message = Message;

});
