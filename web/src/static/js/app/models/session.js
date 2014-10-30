define(function(require, exports, module) {

    var backbone = require('backbone');
    var BackboneSession = require('backbone/sessions');

    var Session = BackboneSession.extend({
        defaults: {
            handle: null,
            sessionid: null
        },

        options: {
            local: true,
            remote: false,
            persist: false
        },

        authenticated: function() {
        }
    });

exports.Session = Session;

});
