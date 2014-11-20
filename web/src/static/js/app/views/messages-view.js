define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/messages-view')
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;

    var MessagesView = marionette.CompositeView.extend({
        childView: MessageView,
        childViewContainer: '#messages',
        template: template,

        ui: {
            submit: '#submitButton',
            postArea: '#postArea',
            messageArea: '#messages',
            messageBox: '#message-box',
        },

        events: {
            'click #submitButton': 'onSubmit',
        },

        onSubmit: function() {
            //alert(this.ui.postArea.val());
            if(this.ui.postArea.val()) {
                var message = new Message({
                    content: this.ui.postArea.val()
                })
                message.save();
                this.collection.add(message);

                this.ui.postArea.val('');
                console.log("Added");
            } else {
                console.log("Unable to add");
            }

        },

        initialize: function(options) {
            this.collection = options.collection;
            this.session = options.session;
            this.collection.fetch();
        }

    });

    exports.MessagesView = MessagesView;
})
