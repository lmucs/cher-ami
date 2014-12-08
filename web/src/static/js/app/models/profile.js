define(function(require, exports, module) {

    var backbone = require('backbone');

    var Profile = Backbone.Model.extend({

        url: 'api/users/',
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

    exports.Profile = Profile;

});
