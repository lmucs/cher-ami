define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/messages-view')
    var MessageView = require('app/views/message-view').MessageView;
    // var postValidator = require('app/utils/post-validator').PostValidator;

    var MessagesView = marionette.CompositeView.extend({
        childView: MessageView,
        childViewContainer: '#messages',
        template: template,

        ui: {
            submit: '#submitButton',
            postArea: '#postArea',
            messageArea: '#messages',
            messageBox: '#message-box',
            // postContainer: '#postContainer'
        },

        events: {
            'click #submitButton': 'onSubmit',
            // 'keyup #postContainer': 'PostValidat'
        },

        onSubmit: function() {
            //alert(this.ui.postArea.val());
            if(this.ui.postArea.val()) {
                this.collection.create({
                    // API requires handle for now,
                    // needs to be removed
                    handle: '',
                    content: this.ui.postArea.val()
                }, {wait: true})
            }

            this.ui.postArea.val('');

            console.log("Added");
        },

        // PostValidat: function(event) {
        //     postValidator(this.ui.postArea)
        // },

        initialize: function(options) {
            this.collection = options.collection;
            this.session = options.session;
            this.collection.fetch();
        }

    });

    exports.MessagesView = MessagesView;
})
