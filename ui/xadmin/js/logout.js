function deleteCookie() {
    document.cookie = 'nbsadmin' + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

function getUrlParameter(name) {
    name = name.replace(/[\[]/, '\\[').replace(/[\]]/, '\\]');
    var regex = new RegExp('[\\?&]' + name + '=([^&#]*)');
    var results = regex.exec(location.search);
    return results === null ? '' : decodeURIComponent(results[1].replace(/\+/g, ' '));
}

function cleanCookie() {
    var urlParams = new URLSearchParams(window.location.search);

    if (urlParams.has('logout') == true){
        deleteCookie()
    }
}



cleanCookie()