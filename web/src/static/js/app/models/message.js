define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = Backbone.Model.extend({
        //url: '/api/Message',
        defaults: {
            messageData: null,
        },

        initialize: function() {

        }
    });

    exports.Message = Message;

});
