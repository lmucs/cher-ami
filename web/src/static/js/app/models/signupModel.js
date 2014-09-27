define(function(require, exports, module) {

    var backbone = require('backbone');

    var SignupModel = Backbone.Model.extend({
        defaults: {

        },

        initialize: function() {
            // https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/now
            this.set('date', Date.now());
        }
    });

exports.SignupModel;

)};