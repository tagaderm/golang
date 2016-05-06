
var entry = document.querySelector('#username');
var output = document.querySelector('h1');
// var check_submit = document.querySelector('#check_username');
// var check_div = document.querySelector('#check_div');
// var update_div = document.querySelector('#update_form');

entry.addEventListener('input', function(){
    console.log('ENTRY: ', entry.value);
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/api/username_check');
    xhr.send(entry.value);
    xhr.addEventListener('readystatechange', function(){
        if (xhr.readyState === 4 && xhr.status === 200) {
            var taken = xhr.responseText;
            console.log('TAKEN:', taken, '\n\n');
            if (taken == 'true') {
                output.textContent = 'Username Taken!';
            } else {
                output.textContent = 'Username not Taken';
            }
        }
    });
});

// check_submit.addEventListener('onclick', function(){
//     check_div.style.display = 'none';
//     update_div.style.display = 'block';
// });