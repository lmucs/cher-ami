define(function(require, exports, module) {

    var marionette = require('marionette');
    var handlebars = require('handlebars');

    // https://github.com/marionettejs/backbone.marionette/blob/master/docs/marionette.templatecache.md
    marionette.TemplateCache.prototype.compileTemplate = function(rawTemplate) {
        return handlebars.compile(rawTemplate);
    }

});