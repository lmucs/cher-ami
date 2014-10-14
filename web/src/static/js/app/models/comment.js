define(function(require, exports, module) {

    var backbone = require('backbone');

    var Comment = Backbone.Model.extend({
        //url: '/api/Message',
        defaults: {
            commentData: null,
        },

        initialize: function() {

        }
    });

    exports.Comment = Comment;

});
