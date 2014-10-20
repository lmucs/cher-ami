define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/comments-view')
    var CommentView = require('app/views/comment-view').CommentView;
    // var postValidator = require('app/utils/post-validator').PostValidator;


    var CommentsView = marionette.CompositeView.extend({
        childView: CommentView,
        childViewContainer: '#comments',
        template: template,

        ui: {
            submit: "#submitButton",
            commentInput: "#commentInput",

        },

        events: {
            'click #submitButton' : 'onSubmit',
            'keydown #commentInput': 'onConfirm',

        },

        onSubmit: function() {
            console.log(this.ui.commentInput.val());
            if(this.ui.commentInput.val()) {
                this.collection.add({
                    commentData: this.ui.commentInput.val()
                })
            }
            this.ui.commentInput.val('');
            console.log("Added");
        },

        initialize: function(options) {
            this.collection = options.collection;
        }

    });

    exports.CommentsView = CommentsView;
})