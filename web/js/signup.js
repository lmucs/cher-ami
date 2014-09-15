$(function() {

  $('#signup').click(function() {
    $.post(
      "http://localhost:8228/signup",
      $("#signupform").serialize()
    );
  });

  $('#pass2').keyup(function(){
    //Store the password field objects into variables ...
    var pass1 = $('#pass1');
    var pass2 = $('#pass2');
    //Store the Confimation Message Object ...
    var message = $('#confirmMessage');
    //Set the colors we will be using ...
    var goodColor = "#66cc66";
    var badColor = "#ff6666";
    //Compare the values in the password field 
    //and the confirmation field
    if(pass1.val() == pass2.val()){
       //The passwords match. 
       //Set the color to the good color and inform
       //the user that they have entered the correct password 
       pass2.css( "background-color", goodColor);
       message.css("color", goodColor);
       message.text("Passwords Match!");
    } else {
       //The passwords do not match.
       //Set the color to the bad color and
       //notify the user.
       pass2.css( "background-color", badColor);
       message.text("Passwords Do Not Match!");
       message.innerHTML = "Passwords Do Not Match!"
      }
  });
  
});