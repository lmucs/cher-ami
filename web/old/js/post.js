var textContainer, textareaSize, input;
var autoSize = function () {
  // also can use textContent or innerText
  textareaSize.innerHTML = input.value + '\n';
};

document.addEventListener('DOMContentLoaded', function() {
  textContainer = document.querySelector('.textarea-container');
  textareaSize = textContainer.querySelector('.textarea-size');
  input = textContainer.querySelector('textarea');

  autoSize();
  input.addEventListener('input', autoSize);
});

$("#submitButton").click(function() {
	var text = $('#postArea').val();
	var testPost = $('#testPost');
	console.log("doop:", text);
	console.log("testPost: ", testPost);
	$('#testPost').append(text);
})
