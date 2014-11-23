define(function(require, exports, module) {
    var marionette = require('marionette');
    var Session = require('app/models/session');
    var app = new marionette.Application();

    // Regions defined in index.html
    app.addRegions({
        headerRegion: '#header',
        mainRegion: '#main',
        sidebarRegion:'#sidebarMain'
    })

    app.addInitializer(function() {
        //http://backbonejs.org/#History
        Backbone.history.start({
            pushState: false
        });
    });

    return app;
})
