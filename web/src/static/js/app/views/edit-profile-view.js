define(function(require, exports, module) {

    var marionette = require('marionette');
    var template = require('hbs!../templates/edit-profile-view')
    var EditProfile = require('app/models/edit-profile').EditProfile

    var EditProfileView = marionette.ItemView.extend({

        template: template,

        ui: {
            firstname: '#first-name',
            lastname: '#last-name',
            gender: '#gender-selector',
            birthday: '#birthday',
            location: '#location',
            bio: '#bio',
            interests: '#interests',
            languages: '#languages',
            submitchanges: '#submitChanges'
        },

        events: {
            'click #submitChanges': 'onFormSubmit',
        },

        onFormSubmit: function() {
            event.preventDefault();
            var req = new EditProfile({
                firstname: this.ui.firstname.val(),
                lastname: this.ui.lastname.val(),
                gender: this.ui.gender.val(),
                birthday: this.ui.birthday.val(),
                location: this.ui.location.val(),
                bio: this.ui.bio.val(),
                interests: this.ui.interests.val(),
                languages: this.ui.languages.val(),
                session: this.session
            });
            console.log(req);
            req.save({patch: true}, {
                type: 'patch'
            });
        },

        initialize: function(options) {
            this.session = options.session;
        },
        
    });

    exports.EditProfileView = EditProfileView;
})
