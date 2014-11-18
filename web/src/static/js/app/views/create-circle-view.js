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
            dropdown: "#visibility",
            submitCircle: "#circleAddMemberButton"
        },

        events: {
            'click #circleAddMemberButton': 'onSubmitCircle'
        },

        onSubmitCircle: function(options) {
            event.preventDefault();
            var req = new CreateCircle({
                circleName: this.ui.name.val(),
                description: this.ui.description.val(),
                visibility: this.ui.visibility.val(),
            });
            console.log(this.ui.name.val()),
            console.log(req)
            req.save();
        }
        

    });

    exports.CreateCircleView = CreateCircleView;
})