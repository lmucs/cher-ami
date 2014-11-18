define(function(require, exports, module) {

    var backbone = require('backbone');

    var Cricle = require('app/models/cricle').Cricle;

    var Cricles = Backbone.Collection.extend({
        url: 'api/circles',
        model: Cricle,

        parse: function(response) {
            // return JSON.parse(response.Objects);
        },

        initialize: function() {

        }
    });

    exports.Cricles = Cricles;

});
