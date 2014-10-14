define(function(require, exports, module) {

    var backbone = require('backbone');

    var Comment = require('app/models/comment').Comment;

    var Comments = Backbone.Collection.extend({
        model: Comment,

        initialize: function() {

        }
    });

    exports.Comments = Comments;

});
