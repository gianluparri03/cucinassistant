// Used by: add.html
function addInput() {
    let div = `<div class="new-element item">` +
                  `<input class="name" type="text" placeholder="Nome" required>` +
                  `<input class="expiration" type="date" placeholder="Scadenza">` +
                  `<input class="quantity" type="number" min="0" step="1" placeholder="Quantit&agrave;">` +
              `</div>`;

    $('#new-elements').append(div);
}

// Used by: add.html
function removeInput() {
    if ($('.new-element').length > 1) {
        $('#new-elements').find('div:last').remove();
    }
}

// Used by: add.html
function saveAddData() {
    var data = "";
    for (item of $('.item')) data += item.children[0].value + ';' + item.children[1].value + ';' + item.children[2].value + '|';
    $('#data').val(data.slice(0, -1));
}

// Used by: edit.html
function saveEditData(e) {
    var re = /^\d+([\+\-\*\/]\d+)?$/;
    if (re.test($('.quantity').replaceAll(' ', ''))) {
        showMessage(false, "Quantit&agrave; non valida");
        e.preventDefault();
        return;
    }

    $('.quantity').val(Math.floor(Math.max(0, eval($('.quantity').val()))) || 0);
}

// Used by: remove.html
function saveRemoveData() {
    var data = "";
    for (c of $('input[type=checkbox]:checked')) data += c.id.slice(4) + ';';
    $('#data').val(data.slice(0, -1));
}
