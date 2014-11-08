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
            sidebarContainer: '#sidebar-container',
            feedContainer: '#feed-container'
        },

        ui: {
            sidebarContainer: '#sidebar-container',
            feedContainer: '#feed-container',
            showContent: '#showContent'
        },

        events: {
            'click #showContent': 'showContent'
        },

        showContent: function(options) { 
            var sidebar = new SidebarView();
            var messages = new Messages();
            var newMessages = new MessagesView();
            var feed = new MessagesView({
                    collection: messages
            });

            this.ui.feedContainer.html(feed.el);
            this.ui.sidebarContainer.html(sidebar.el);
            console.log(feed.el);
            //messages.render();
            sidebar.render();
        },

        initialize: function(options) {

        }

    });

    exports.HomeLayout = HomeLayout;
})