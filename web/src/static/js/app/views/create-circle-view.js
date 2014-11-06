define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/create-circle-view')

    var CreateCircleView = marionette.ItemView.extend({
        template: template,

        ui: {
            dropdown: "#visibility"
        },

        events: {
            'click #visibility': 'onDropdownClick'
        },

        onDropdownClick: function () {
            console.log("I got clicked!!");
            // if($(this).hasClass('dropdown-menu-form')) {
            //     e.stopPropagation();
            // }
        }
        

    });

    exports.CreateCircleView = CreateCircleView;
})