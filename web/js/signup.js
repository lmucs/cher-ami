$(function() {

  $('#signup').click(function() {
    $.post(
      "http://localhost:8228/signup",
      $("#signupform").serialize()
    );
  });

})