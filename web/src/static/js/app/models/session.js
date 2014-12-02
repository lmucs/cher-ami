define(function(require, exports, module) {

    var backbone = require('backbone');
    var BackboneSession = require('backbone/sessions');

    var Session = BackboneSession.extend({
        defaults: {
            handle: null,
            token: null
        },

        options: {
            local: true,
            remote: false,
            persist: false
        },

        hasAuth: function() {
            return this.has("token")
        },

        getTokenValue: function() {
            return this.get("token")
        },

        getHandle: function() {
            return this.get("handle")
        }
    });

    exports.Session = Session;

});
