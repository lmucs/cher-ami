define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/messages-view')
    var MessageView = require('app/views/message-view').MessageView;
    
    var postValidator = require('app/utils/post-validator').PostValidator;

    var MessagesView = marionette.CompositeView.extend({
        childView: MessageView,
        childViewContainer: '#messages',
        template: template,

        ui: {
            submit: '#submitButton',
            postArea: '#postArea',
<<<<<<< HEAD
            messageArea: '#messages',
            messageBox: '#message-box',
            postContainer: '#postContainer'
=======
>>>>>>> FETCH_HEAD
        },

        events: {
            'keyup #postContainer': 'PostValidat',
            'click #submitButton': 'onSubmit',
            'keydown #postArea': 'onConfirm'
        },

        onSubmit: function() {
            if(this.ui.postArea.val() !== '') {
                this.collection.add({
                    messageData: this.ui.postArea.val()
                })
            }

            this.ui.postArea.val('');

            console.log("Added");
        },

        PostValidat: function(event) {
            postValidator(this.ui.postArea)
        },

        onConfirm: function(event) {
            var enterkey = 13;
            if (event.which == enterkey) {
                this.collection.add({
                    messageData: this.ui.postArea.val()
                })
                this.ui.postArea.val('');
                console.log(this.ui.postArea.val().length);
            }
            console.log("Enter key was pressed");
        },

        initialize: function(options) {
            this.collection = options.collection;
        }

    });

    exports.MessagesView = MessagesView;
})