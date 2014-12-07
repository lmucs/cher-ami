define(function(require, exports, module) {
    var marionette = require('marionette');
    var Session = require('app/models/session').Session;

    var CherAmi = marionette.Application.extend({
        initialize: function(options) {
            this.session = new Session({}, {
                remote: false
            });
            if (this.session.hasAuth()) {
                $.ajaxSetup({
                    headers: {'Authorization' : this.session.getTokenValue()}
                })
            }
        }
    })

    var app = new CherAmi();

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
        $.ajaxSetup({
            statusCode: {
                401: function() {
                    window.location.reload()
                }
            }
        })
    });

    return app;
})
