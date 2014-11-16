define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/create-circle-view')
    var CreateCircle = require('app/models/create-circle').CreateCircle;

    var CreateCircleView = marionette.ItemView.extend({
        template: template,

        ui: {
            name: "#circle-name",
            description: "#description",
            visibility: "#visibilitySelector",
            dropdown: "#visibility"
            submitCircle: "#circleAddMemberButton"
        },

        events: {
            'click #visibility': 'onDropdownClick',
            'click #circleAddMemberButton': 'onSubmitCircle'
        },
        onSubmitCircle: function(options) {

        },
        
        onDropdownClick: function (e) {
            console.log("I got clicked!!");
            console.log("this: ", document.getElementById("dropdown-menu-form"));

            if(document.getElementById("dropdown-menu-form")) {
                console.log("I got here");
                e.stopPropagation();
            }
        }
        

    });

    exports.CreateCircleView = CreateCircleView;
})