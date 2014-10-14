define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/comment-view')
    var Comment = require('app/models/comment').Comment

    var CommentView = marionette.ItemView.extend({
        model: Comment,
        template: template,

        ui: {

        },

        events: {
        },

        initialize: function(options) {
            this.model = options.model
        }

    });

    exports.CommentView = CommentView;
})
