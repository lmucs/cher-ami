define(function(require, exports, module) {

    var backbone = require('backbone');
    var marionette = require('marionette');
    var app = require('app/app');
    var Session = require('app/models/session').Session;
    //var session = require('backbone/sessions');
    var $ = require('jquery');

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
    var CircleView = require('app/views/circle-view').CircleView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;

    var LandingLayout = require('app/layouts/landing-layout').LandingLayout;
    var HomeLayout = require('app/layouts/home-layout').HomeLayout;
    var CircleLayout = require('app/layouts/circle-layout').CircleLayout;
    var ProfileLayout = require('app/layouts/profile-layout').ProfileLayout;

    var AppController = marionette.Controller.extend({

        initialize: function(options) {
            this.app = app;
            this.app.session = new Session({}, {
                remote: false
            });

            // Initialize header view.
            this.app.headerRegion.show(new HeaderView());

            // // Logic for auth check.
            if (this.app.session.has('sessionid')) {
                console.log("User logged in.");
                $.ajaxSetup({
                    headers: {'Authorization' : this.app.session.get('sessionid')}
                })
                // user is authed, redirect home
                this.app.mainRegion.show(new HomeLayout({
                    session: this.app.session
                }));
            } else {
                this.app.mainRegion.show(new LandingLayout({
                    session: this.app.session
                }));
            }
            
            // this.app.mainRegion.show(new ProfileLayout())
        },

        // Needed for AppRouter to initialize index route.
        index: function() {

        }

    });

    exports.AppController = AppController;

});
