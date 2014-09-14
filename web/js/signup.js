$(function() {

  $('#signupform').submit(function( event ) {
    $.post(function (data) {
      console.log(data);
    })
  });

})