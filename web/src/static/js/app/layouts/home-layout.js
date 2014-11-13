define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/layouts/home-layout')
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    var ProfileView = require('app/views/profile-view').ProfileView;
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView;

    var HomeLayout = marionette.LayoutView.extend({
        template: template,

        regions: {
            feed: '#feed-container',
            circle: '#create-circle-view',
            profile: '#profile-container',
            editProfile:'#edit-profile-container'
        },

        ui: {
            feedContainer: '#feed-container',
            showContent: '#showContent',
            createCircle: '#goToCreateCircle',
            displayProfile: '#goToProfile',
            editProfile: '#editProfile'
        },

        events: {
            'click #goToCreateCircle': 'showCreateCircle',
            'click #goToProfile': 'showProfile',
            'click #editProfile': 'showEditProfile',
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

        showProfile: function(options) {
            var showProfile = new ProfileView();
            this.profile.show(showProfile);
        },

        showCreateCircle: function(options) {
            var createCircle = new CreateCircleView();
            this.profile.show(createCircle);
        },

        showEditProfile: function(options) {
            var editProfile = new EditProfileView();
            this.profile.show(editProfile);
        }

    });
    exports.HomeLayout = HomeLayout;
})
