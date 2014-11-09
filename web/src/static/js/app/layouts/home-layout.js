define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/home-layout')
    var SidebarView = require('app/views/sidebar-view').SidebarView;
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;

    var HomeLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            sidebar: '#sidebar-container',
            feed: '#feed-container'
        },

        ui: {
            sidebarContainer: '#sidebar-container',
            feedContainer: '#feed-container',
            showContent: '#showContent'
        },

        events: {
        },

        initialize: function(options) {
        },

        onRender: function() {
            var sidebar = new SidebarView();
            var messages = new Messages();
            var feed = new MessagesView({
                collection: messages
            });
            this.sidebar.show(sidebar);
            this.feed.show(feed);
        }

    });

    exports.HomeLayout = HomeLayout;
})
