define(function(require, exports, module) {
    var PostValidator = function(form) {

        $('#postContainer').bootstrapValidator({
            message: 'This value is not valid',
            feedbackIcons: {
                valid: 'glyphicon glyphicon-ok',
                invalid: 'glyphicon glyphicon-remove',
                validating: 'glyphicon glyphicon-refresh'
            },

            fields: {
                postArea: {
                    message: 'The message has to have content',
                    validators: {
                        notEmpty: {
                            message: 'The message is required and cannot be empty'
                        },
                        stringLength: {
                            min: 1,
                            max: 126,
                            message: 'The message must be more than 1 and less than 126 characters long'
                        },
                    }
                },
            }
        });
    }
    
   exports.PostValidator = PostValidator;
});