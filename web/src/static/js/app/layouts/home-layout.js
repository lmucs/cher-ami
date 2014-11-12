define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/home-layout')
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;

    var HomeLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            feed: '#feed-container',
            circle: '#create-circle-view'
        },

        ui: {
            feedContainer: '#feed-container',
            showContent: '#showContent',
            createCircle: '#goToCreateCircle'
        },

        events: {
            'click #goToCreateCircle': 'showCreateCircle'
        },

        initialize: function(options) {
        },

        

        onRender: function() {
            var messages = new Messages();
            var feed = new MessagesView({
                collection: messages
            });
            this.feed.show(feed);
        },
        showCreateCircle: function(options) {
            var createCircle = new CreateCircleView();
            this.circle.show(createCircle);
        }

    });
    exports.HomeLayout = HomeLayout;
})
