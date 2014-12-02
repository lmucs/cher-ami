define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/create-circle-view')
    var CircleView = require('app/views/circle-view').CircleView;
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
                public: this.ui.visibility.val() === 'public',
                visibility: this.ui.visibility.val(),
            });
            console.log(req)
            req.save();
            console.log("Circle Added");
            this.collection.add(req);
        },

        initialize: function(options) {

        }
        

    });

    exports.CreateCircleView = CreateCircleView;
})