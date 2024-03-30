// Used by: add.html, edit.html
function focusDate(target) {
    var input = $(target);
    var timestamp = input.attr('timestamp');
    input.val(timestamp);
    input.attr('type', 'date');
    input.click();
}

// Used by: add.html, edit.html
function blurDate(target) {
    var input = $(target);
    var timestamp = input.val();
    input.attr('timestamp', timestamp);
    input.attr('type', 'text');
    input.val(timestamp ? new Date(timestamp).toLocaleDateString() : '');
}

// Used by: add.html
function addInput() {
    let div = `<div class="new-element item">` +
                  `<input class="name" type="text" placeholder="Nome" required>` +
                  `<input class="expiration" type="text" placeholder="Scadenza" onblur="blurDate(this);" onfocus="focusDate(this);" timestamp="">` +
                  `<input class="quantity" type="number" min="0" placeholder="Quantit&agrave;">` +
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
    for (item of $('.item')) data += item.children[0].value + ';' + $(item.children[1]).attr('timestamp') + ';' + item.children[2].value + '|';
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

    $('.quantity').val(Math.floor(Math.max(0, eval($('.quantity').val()))));
    $('.expiration').val($('.expiration').attr('timestamp'));
}

// Used by: remove.html
function saveRemoveData() {
    var data = "";
    for (c of $('input[type=checkbox]:checked')) data += c.id.slice(4) + ';';
    $('#data').val(data.slice(0, -1));
}

// Used by: every file except add.html
window.onload = function() {
    for (input of $('.expiration')) input.value = new Date(input.value).toLocaleDateString();
};
