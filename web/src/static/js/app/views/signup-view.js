define(function(require, exports, module) {
    var marionette = require('marionette');
    var template = require('hbs!../templates/signup-view');

    var Signup = require('app/models/signup').Signup;

    var SignupView = marionette.ItemView.extend({
        template: template,

        //takes the div marionette creates and give it a class named mainContainer.
        tagName: "div",
        className: "mainContainer",
        ui: {
            handle: '#handle',
            email: '#input-email',
            pass: '#pass1',
            confirmPass: '#pass2',
            rememberMe: '#remember-me',
            signup: '#signup'
        },

        events: {
            'click #remember-me': 'onRememberConfirm',
            'keyup #pass2': 'passwordMatch',
            'click #signup': 'onFormConfirm'
        },

        initialize: function(options) {

        },

        onRememberConfirm: function(options) {
            // Session-request method goes here
        },

        onFormConfirm: function(event) {
            event.preventDefault();
            var req = new Signup({
                handle: this.ui.handle.val(),
                email: this.ui.email.val(),
                password: this.ui.pass.val(),
                confirmpassword: this.ui.confirmPass.val()
            });
            console.log(req)
            req.save();
        },

        passwordMatch: function(event) {
            //Store the password field objects into variables ...
            var pass1 = $('#pass1');
            var pass2 = $('#pass2');
            //Store the Confimation Message Object ...
            var message = $('#confirmMessage');
            //Set the colors we will be using ...
            var goodColor = "#66cc66";
            var badColor = "#ff6666";
            //measure password strength
            var getStrength = function(password) {
               var strength = 0;
               if (pass1.length > 7) {strength += 1};
               if (pass1.val().match(/([a-zA-Z])/) && pass1.val().match(/([0-9])/))  {strength += 1};
               if (pass1.val().match(/([!,%,&,@,#,$,^,*,?,_,~])/))  {strength += 1};
               if (pass1.val().match(/(.*[!,%,&,@,#,$,^,*,?,_,~].*[!,%,&,@,#,$,^,*,?,_,~])/)) {strength += 1};
               if(strength < 2) {
                return 'Weak';
               } else if(strength == 2) {
                return 'Strong';
               } else {
                return 'Very Strong';
               }
            }

            if(pass1.val() == pass2.val() && pass1.val().length >= 8) {
               pass1.css( "background-color", goodColor);
               pass2.css( "background-color", goodColor);
               message.css("color", goodColor);
               message.text("Passwords Match!" + " " + getStrength(pass1));
               return true;
            } else if(pass1.val().length < 8) {
               pass1.css( "background-color", badColor);
               pass2.css( "background-color", badColor);
               message.css("color", badColor);
               message.text("Password has to be more than 8 charecters");
            } else {
               //The passwords do not match.
               //Set the color to the bad color and
               //notify the user.
               pass1.css( "background-color", badColor);
               pass2.css( "background-color", badColor);
               message.text("Passwords Do Not Match!");
               message.innerHTML = "Passwords Do Not Match!"
            }
        }

    });

    exports.SignupView = SignupView;
})
