define(function(require, exports, module) {

    var backbone = require('backbone');

    var EditProfile = Backbone.Model.extend({

        // url: '/api/signup',
        defaults: {
            firstname: null,
            lastname: null,
            gender: null,
            birthday: null,
            location: null,
            bio: null,
            interests: null,
            languages: null
        },

        initialize: function() {

        }
    });

exports.EditProfile = EditProfile;

});
