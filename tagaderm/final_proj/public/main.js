
var entry = document.querySelector('#username');
var output = document.querySelector('h1');

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