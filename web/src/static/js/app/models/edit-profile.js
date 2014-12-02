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

        initialize: function() {
            var handle = this.get('session').getHandle();
            this.unset('session');
            this.url = "/api/users/" + handle;
        }
    });

    exports.EditProfile = EditProfile;

});
