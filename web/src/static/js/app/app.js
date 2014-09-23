define(function(require, exports, module) {
    var marionette = require('marionette');
    var app = new marionette.Application();

    // Regions defined in index.html
    app.addRegions({
        headerRegion: '#header',
        mainRegion: '#main',
        footerRegion: '#footer'
    })

    app.addInitializer(function() {
        //http://backbonejs.org/#History
        Backbone.history.start({
            pushState: false
        });
    });

    return app;
})
