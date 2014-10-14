define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/comments-view')
    var CommentView = require('app/views/comment-view').CommentView;

    var CommentsView = marionette.CompositeView.extend({
        childView: CommentView,
        childViewContainer: '#messages',
        template: template,

        ui: {

        },

        events: {

        },

    });

    exports.CommentsView = CommentsView;
})