define(function(require, exports, module) {
	var InputValidator = function(form) {

		$('#signupform').bootstrapValidator({
	        message: 'This value is not valid',
	        feedbackIcons: {
	            valid: 'glyphicon glyphicon-ok',
	            invalid: 'glyphicon glyphicon-remove',
	            validating: 'glyphicon glyphicon-refresh'
	        },

	        fields: {
	        	handle: {
	        		message: 'The username is not valid',
	        		validators: {
	        			notEmpty: {
                            message: 'The username is required and cannot be empty'
	                    },
	                    stringLength: {
	                        min: 6,
	                        max: 30,
	                        message: 'The username must be more than 6 and less than 30 characters long'
	                    },
	                    regexp: {
	                        regexp: /^[a-zA-Z0-9]+$/,
	                        message: 'The username can only consist of alphabetical and number'
	                    },
	                    different: {
	                        field: 'password',
	                        message: 'The username and password cannot be the same as each other'
	                    }
	        		}
	        	},

	            email: {
	                validators: {
	                    notEmpty: {
	                        message: 'The email is required and cannot be empty'
	                    },
	                    emailAddress: {
	                        message: 'The input is not a valid email address'
	                    }
	                }
	            },

	            password: {
	            	validators: {
	            		notEmpty: {
                            message: 'The password is required and cannot be empty'
                        },
                    	different: {
                            field: 'handle',
                            message: 'The password cannot be the same as username'
                        },
                        stringLength: {
                            min: 8,
                             message: 'The password must have at least 8 characters'
                        },
                        identical: {
                            field: 'confirmPassword',
                            message: 'The password and its confirm are not the same'
                        }
	            	}
	            },

	            confirmPassword: {
	                validators: {
	                	identical: {
                        	field: 'password',
                            message: 'The password and its confirm are not the same'
                        }
	                }
	            }

	        }
        });
	}
   exports.InputValidator = InputValidator;

});