define(function(require, exports, module) {

        var InputValidator = function(form) {
            //form validation rules
            $("#signupform").validate({
                rules: {
                    handle: {
                        required: true,
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
                },
            });
        }

    exports.InputValidator = InputValidator;
});