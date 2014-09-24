define(function(require, exports, module) {

    var backbone = require('backbone');
    var marionette = require('marionette');

    //http://www.safaribooksonline.com/library/view/getting-started-with/9781783284252/ch06s02.html

    var vent = new backbone.Wreqr.EventAggregator();

    exports.vent = vent;

});