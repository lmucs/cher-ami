define(function(require, exports, module) {

    var backbone = require('backbone');

    var Message = require('app/models/message').Message;

    var Messages = Backbone.Collection.extend({
        model: Message,

        initialize: function() {

        }
    });

    exports.Messages = Messages;

});
