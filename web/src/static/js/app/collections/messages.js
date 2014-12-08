define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = require('app/models/message').Message;

    var Messages = Backbone.Collection.extend({
        url: 'api/messages',
        model: Message,

        parse: function(response) {
            return response;
        },

        initialize: function() {

        }
    });

    exports.Messages = Messages;

});
