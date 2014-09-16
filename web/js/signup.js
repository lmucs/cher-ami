$(function() {

  $('#signup').click(function() {
    $.post(
      "http://localhost:8228/signup",
      $("#signupform").serialize()
    );
  });
  
  $('#pass2').keyup(function() {
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

    //Compare the values in the password field 
    //and the confirmation field
    if(pass1.val() == pass2.val() && pass1.val().length >= 6) {
       //The passwords match. 
       //Set the color to the good color and inform
       //the user that they have entered the correct password 
       pass1.css( "background-color", goodColor);
       pass2.css( "background-color", goodColor);
       message.css("color", goodColor);
       message.text("Passwords Match!" + " " + getStrength(pass1));
    } else if(pass1.val().length < 6) {
       pass1.css( "background-color", badColor);
       pass2.css( "background-color", badColor);
       message.css("color", badColor);
       message.text("Password has to be more than 6 charecters");
    } else {
       //The passwords do not match.
       //Set the color to the bad color and
       //notify the user.
       pass1.css( "background-color", badColor);
       pass2.css( "background-color", badColor);
       message.text("Passwords Do Not Match!");
       message.innerHTML = "Passwords Do Not Match!"
    }

  });

});