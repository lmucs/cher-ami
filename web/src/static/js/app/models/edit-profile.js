define(function(require, exports, module) {

    var backbone = require('backbone');

    var EditProfile = Backbone.Model.extend({

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

        initialize: function(options) {
            // this.session = options.session;
            console.log(this.get('firstname'));
            this.url = "/api/users/" + this.get('firstname');
        }
    });

    exports.EditProfile = EditProfile;

});
