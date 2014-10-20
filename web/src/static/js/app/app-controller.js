define(function(require, exports, module) {

    var backbone = require('backbone');
    var marionette = require('marionette');
    var app = require('app/app');
    var Sessions = require('backbone/sessions')

    var HeaderView = require('app/views/header-view').HeaderView;
    var SignupView = require('app/views/signup-view').SignupView;
    var FooterView = require('app/views/footer-view').FooterView;
    var LoginView = require('app/views/login-view').LoginView;
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    var Message = require('app/models/message').Message;
    var Messages = require('app/collections/messages').Messages;
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView;
    var CommentsView = require('app/views/comments-view').CommentsView;
    var CommentView = require('app/views/comment-view').CommentView;
    var Comment = require('app/models/comment').Comment;
    var Comments = require('app/collections/comments').Comments;
    var AppController = marionette.Controller.extend({

        initialize: function(options) {
            this.app = app;
            /*var sessions = new Sessions()
            this.app.session = sessions.APP.Session({
                handle: "test",
                password: "12345678"
            }, {});*/

            var test = new Messages();
            var testComment = new Comments();

            // Initialization of views will go here.
            this.app.headerRegion.show(new HeaderView());
            // this.app.mainRegion.show(new MessagesView({
            //     collection: test
            // }));
            // this.app.mainRegion.show(new CommentsView({
            //     collection: testComment
            // }));
            this.app.mainRegion.show(new EditProfileView());
            // this.app.footerRegion.show(new FooterView());
        },

        // Needed for AppRouter to initialize index route.
        index: function() {

        }

    });

    exports.AppController = AppController;

});
