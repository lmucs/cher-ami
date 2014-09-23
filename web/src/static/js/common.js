// Require.js Configuration
require.config({
  baseUrl: 'static/js',

  paths : {
    'marionette': 'vendor/backbone/marionette'
  },

  packages: [

        {
            location: 'app',
            name: 'app'
        },

        {
            location: 'vendor/jquery',
            name: 'jquery',
            main:'jquery'
        },

        {
            location: 'vendor/backbone',
            name: 'backbone',
            main:'backbone'
        },

        {
            location: 'vendor/hbs',
            name: 'hbs',
            main:'hbs'
        }
    ],

    map: {
        '*': {
            'underscore': 'vendor/underscore/lodash',
            'handlebars': 'hbs/handlebars',
        },
    },

  hbs: {
        templateExtension : 'html',
        disableI18n : true,
        helperDirectory: 'app/shared/hbs'
  },

  shim : {

    'backbone': {
        'deps': ['jquery', 'underscore'],
        'exports': 'Backbone'
    },

    'marionette': {
        'deps': ['jquery', 'underscore', 'backbone'],
        'exports': 'Marionette'
    }
  },

  wrapShim: true,

});