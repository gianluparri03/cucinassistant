window.onload = function enableCookieBanner() {
    if (!localStorage.getItem('cookies-accepted')) {
        document.getElementById('obfuscator').style.display = 'block';
    } else {
        document.getElementById('obfuscator').innerHTML = '';
    }
}

function toggleCookieBanner() {
    localStorage.setItem('cookies-accepted', true);
    enableCookieBanner();
}
