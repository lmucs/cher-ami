define(function(require, exports, module) {

   var PassCheck = function(field1, field2, messageLocation) {
      //Store the Confimation Message Object ...
      var message = messageLocation;
      //Set the colors we will be using ...
      var goodColor = "#66cc66";
      var badColor = "#ff6666";
      //measure password strength
      var getStrength = function(password) {
         var strength = 0;
         if (field1.length > 7) {strength += 1};
         if (field1.val().match(/([a-zA-Z])/) && field1.val().match(/([0-9])/))  {strength += 1};
         if (field1.val().match(/([!,%,&,@,#,$,^,*,?,_,~])/))  {strength += 1};
         if (field1.val().match(/(.*[!,%,&,@,#,$,^,*,?,_,~].*[!,%,&,@,#,$,^,*,?,_,~])/)) {strength += 1};
         if(strength < 2) {
          return 'Weak';
         } else if(strength == 2) {
          return 'Strong';
         } else {
          return 'Very Strong';
         }
      }

      if(field1.val() == field2.val() && field1.val().length >= 8) {
         field1.css( "background-color", goodColor);
         field2.css( "background-color", goodColor);
         message.css("color", goodColor);
         message.text("Passwords Match!" + " " + getStrength(field1));
         return true;
      } else if(field1.val().length < 8) {
         field1.css( "background-color", badColor);
         field2.css( "background-color", badColor);
         message.css("color", badColor);
         message.text("Password has to be more than 8 charecters");
      } else {
         //The passwords do not match.
         //Set the color to the bad color and
         //notify the user.
         field1.css( "background-color", badColor);
         field2.css( "background-color", badColor);
         message.text("Passwords Do Not Match!");
         message.innerHTML = "Passwords Do Not Match!"
      }
   }

   exports.PassCheck = PassCheck;

});