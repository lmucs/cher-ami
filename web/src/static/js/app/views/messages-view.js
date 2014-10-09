define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/messages-view')
    var MessageView = require('app/views/message-view').MessageView;

    var MessagesView = marionette.CompositeView.extend({
        childView: MessageView,
        childViewContainer: '#messages',
        template: template,

        ui: {
            submit: '#submitButton',
            postArea: '#postArea'
        },

        events: {
            'click #submitButton': 'onSubmit',
            'keypress postArea': 'onConfirm'
        },

        onSubmit: function() {
            this.collection.add({
                messageData: this.ui.postArea.val()
            })
            console.log("Added");
        },

        onConfirm: function(event) {
            var enterkey = 13;
            if (event.which === enterkey) {
                this.collection.add({
                    messageData: ui.postArea.val()
                })
            }
        },

        initialize: function(options) {
            this.collection = options.collection;
        }

    });

    exports.MessagesView = MessagesView;
})