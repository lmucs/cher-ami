define(function(require, exports, module) {
    var PostValidator = function(form) {

        $('#postContainer').bootstrapValidator({
            container: '#warningMessages',
            feedbackIcons: {
                valid: 'glyphicon glyphicon-ok',
                invalid: 'glyphicon glyphicon-remove',
                validating: 'glyphicon glyphicon-refresh'
            },

            fields: {
                postAreaa: {
                    validators: {
                        notEmpty: {
                            message: 'The content is required and cannot be empty'
                        },
                        stringLength: {
                            min: 1,
                            max: 126,
                            message: 'The content must be less than 500 characters long'
                        },
                    }
                }
            }
        });
    }
    
   exports.PostValidator = PostValidator;
});