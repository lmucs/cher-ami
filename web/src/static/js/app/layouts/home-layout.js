define(function(require, exports, module) {
    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/home-layout')

    //messages requirements
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;

    //profile requirements
    var ProfileView = require('app/views/profile-view').ProfileView;
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView;
    var SettingsView = require('app/views/settings-view').SettingsView;

    var SidebarView = require('app/views/sidebar-view').SidebarView;

    var HomeLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            sidebar: '#sidebar-view',
            feed: '#feed-container',
            circle: '#create-circle-view',
            profile: '#profile-container',
        },

        ui: {
            feedContainer: '#feed-container',
            showContent: '#showContent',
            showCircles: '#goToCircles',
            createCircle: '#goToCreateCircle',
        },

        events: {
            'click #goToHome': 'showFeed',
        },

        initialize: function(options) {
        },

        onRender: function() {
            // this.sidebar.show(new SidebarView());
            this.showFeed();
        },

        showFeed: function(options) {
            var messages = new Messages();
            var feed = new MessagesView({
                collection: messages
            });
            this.profile.show(feed);
        },

    });
    exports.HomeLayout = HomeLayout;
})
