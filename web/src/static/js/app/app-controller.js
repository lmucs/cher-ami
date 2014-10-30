define(function(require, exports, module) {

    var backbone = require('backbone');
    var marionette = require('marionette');
    var app = require('app/app');
    var Session = require('app/models/session').Session;

    var HeaderView = require('app/views/header-view').HeaderView;
    var SignupView = require('app/views/signup-view').SignupView;
    var FooterView = require('app/views/footer-view').FooterView;
    var LoginView = require('app/views/login-view').LoginView;
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView;
    var ProfileView = require('app/views/profile-view').ProfileView;
    var CommentsView = require('app/views/comments-view').CommentsView;
    var CommentView = require('app/views/comment-view').CommentView;
    var SidebarView = require('app/views/sidebar-view').SidebarView;
    var Comment = require('app/models/comment').Comment;
    var Comments = require('app/collections/comments').Comments;
    var AppController = marionette.Controller.extend({

        initialize: function(options) {
            this.app = app;
            this.app.session = new Session();
            // if (this.app.session.authenticated()) {
            //     // user is authed, redirect home
            //     this.app.mainRegion.show(new ProfileView());
            // } else {
            //     this.app.mainRegion.show(new LoginView({
            //         session: this.app.session
            //     }));
            // }

            //var test = new Messages();
            //var testComment = new Comments();

            // Initialization of views will go here.
            this.app.headerRegion.show(new HeaderView());
            //this.app.mainRegion.show(new SidebarView());
            // this.app.mainRegion.show(new MessagesView({
            //     collection: test
            // }));
            // this.app.mainRegion.show(new CommentsView({
            //     collection: testComment
            // }));
            this.app.mainRegion.show(new ProfileView());
            //this.app.mainRegion.show(new SignupView());
            // this.app.footerRegion.show(new FooterView());
        },

        // Needed for AppRouter to initialize index route.
        index: function() {

        }

    });

    exports.AppController = AppController;

});
