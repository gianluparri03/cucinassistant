// Used by: edit.html
function deleteMenu() {
    var f = document.forms[0];
    f.action = f.action.replace('modifica', 'elimina');
    f.submit();
}

// Used by: edit.html
function saveMenu() {
    var data = "";
    for (var i = 0; i < 14; i++) data += $('#entry-' + i).val() + ';';
    $('#data').val(data.slice(0, -1));
}
