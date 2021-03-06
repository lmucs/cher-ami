define(function(require, exports, module) {
    var backbone = require('backbone');
    var marionette = require('marionette');
    var app = require('app/app');
    var $ = require('jquery');

    /** Views **/
    var HeaderView = require('app/views/header-view').HeaderView;
    var SignupView = require('app/views/signup-view').SignupView;
    var FooterView = require('app/views/footer-view').FooterView;
    var LoginView = require('app/views/login-view').LoginView;
    /** Messages views **/    
    var MessagesView = require('app/views/messages-view').MessagesView;
    var MessageView = require('app/views/message-view').MessageView;
    /** Profile views **/    
    var EditProfileView = require('app/views/edit-profile-view').EditProfileView;
    var ProfileView = require('app/views/profile-view').ProfileView;
    /** Comments views **/
    var CommentsView = require('app/views/comments-view').CommentsView;
    var CommentView = require('app/views/comment-view').CommentView;
    /** Circle views **/
    var CircleView = require('app/views/circle-view').CircleView;
    var CreateCircleView = require('app/views/create-circle-view').CreateCircleView;
    /** Settings views **/
    var SettingsView = require('app/views/settings-view').SettingsView;
    /** Sidebar views **/
    var SidebarView = require('app/views/sidebar-view').SidebarView;

    /** Models **/
    var Session = require('app/models/session').Session;
    var Message = require('app/models/message').Message;
    var Comment = require('app/models/comment').Comment;
    var EditProfile = require('app/models/edit-profile').EditProfile;

    /** Collections **/
    var Messages = require('app/collections/messages').Messages;
    var Comments = require('app/collections/comments').Comments;

    /** Layouts **/
    var LandingLayout = require('app/layouts/landing-layout').LandingLayout;
    var HomeLayout = require('app/layouts/home-layout').HomeLayout;
    var CircleLayout = require('app/layouts/circle-layout').CircleLayout;
    var ProfileLayout = require('app/layouts/profile-layout').ProfileLayout;

    var AppController = marionette.Controller.extend({
        initialize: function() {
            this.app = app;
            this.app.session = new Session({}, {
                remote: false
            });
            // Logic for auth check.
            this.app.headerRegion.show(new HeaderView({
                session: this.app.session
            }));
            if (this.app.session.has('token')) {
                console.log("User logged in");
                $.ajaxSetup({
                    headers: {'Authorization' : this.app.session.get('token')}
                })
                // user is authed, redirect home
                this.app.mainRegion.show(new HomeLayout({
                    session: this.app.session
                }));
                // Initialize header view after logged in
            this.app.session = app.session;
            };
        },
        // Needed for AppRouter to initialize index route.
        index: function() {
            console.log(this.app.session)
            if (this.app.session.hasAuth()) {
                this.showHomeLayout();
            } else {
                this.app.mainRegion.show(new LandingLayout({
                    session: this.app.session
                }));
            }
        },

        showHomeLayout: function(options) {
            if (!this.app.session.hasAuth()) {
                this.index();
            } else {
                this.app.sidebarRegion.show(new SidebarView())
                this.app.headerRegion.show(new HeaderView({
                    session: this.app.session
                }));
                this.app.mainRegion.show(new HomeLayout({
                    session: this.app.session
                }));
            }
        },

        showCircle: function(options) {
            if (!this.app.session.hasAuth()) {
                this.index();
            } else {
                this.app.headerRegion.show(new HeaderView({
                    session: this.app.session
                }));
                this.app.sidebarRegion.show(new SidebarView())
                this.app.mainRegion.show(new CircleLayout({
                    session: this.app.session
                }))
            }
        },

        showProfile: function(options) {
            this.app.headerRegion.show(new HeaderView({
                session: this.app.session
            }));
            this.app.sidebarRegion.show(new SidebarView())
            this.app.mainRegion.show(new ProfileLayout({
                session: this.app.session
            }))
        },

        showSettings: function(options) {
            this.app.headerRegion.show(new HeaderView({
                session: this.app.session
            }));
            this.app.sidebarRegion.show(new SidebarView())
            this.app.mainRegion.show(new SettingsView({
                session: this.app.session
            })) 
        }
    });
    exports.AppController = AppController;
});