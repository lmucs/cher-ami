define(function(require, exports, module) {

        var InputValidator = function(form) {
            //form validation rules
            $("#signupform").validate({
                errorClass: "invalid",
                validClass: "success",
                rules: {
                    handle: {
                        required: true,
                        minlength: 1,
                    },
                    email: {
                        required: true,
                        email: true
                    },
                    password: {
                        required: true,
                        minlength: 8
                    },
                    confirmPassword: {
                        required: true,
                        equalTo: "#pass1"
                    },
                },

                messages: {
                    handle: "Please enter your handle",
                    password: {
                        required: "Please provide a password",
                        minlength: "Your password must be at least 8 characters long"
                    },
                    confirmPassword: {
                        required: "Please provide a password",
                        minlength: "Your password must be at least 8 characters long"
                    },
                    email: "Please enter a valid email address",
                }
            });

            $('#loginForm').validate({
                errorClass: "invalid",
                validClass: "success",
                rules: {
                    handle: {
                        required: true,
                        minlength: 1,
                    },
                    password: {
                        required: true,
                        minlength: 8
                    }
                },
                messages: {
                    handle: "Please enter your handle",
                    password: {
                        required: "Please provide a password",
                        minlength: "Remember, your password is 8 or more characters"
                    },
                }
            });

            //disable signup button when credentials are not correct
            $('#signupform input').on('keyup blur', function () { // fires on every keyup & blur
                if ($('#signupform').valid()) {                   // checks form for validity
                    $('#signup').prop('disabled', false);        // enables button
                } else {
                    $('#signup').prop('disabled', 'disabled');   // disables button
                }
            });

            //disable login button when credentials are not correct
            $('#loginForm input').on('keyup blur', function () { // fires on every keyup & blur
                if ($('#loginForm').valid()) {                   // checks form for validity
                    $('#login').prop('disabled', false);        // enables button
                } else {
                    $('#login').prop('disabled', 'disabled');   // disables button
                }
            });
        }

    exports.InputValidator = InputValidator;
});