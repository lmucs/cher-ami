define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/message-view')
    var Message = require('app/models/message').Message

    var MessageView = marionette.ItemView.extend({
        model: Message,
        template: template,

        ui: {

        },

        events: {
        },

        initialize: function(options) {
            this.model = options.model
        }

    });

    exports.MessageView = MessageView;
})
