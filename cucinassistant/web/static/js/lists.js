// Used by: add.html
function addInput() {
    let div = `<div class="new-element"><input class="name" type="text" placeholder="Nome" required></div>`;
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
    for (input of $('.new-element > input')) data += input.value + ';';
    $('#data').val(data.slice(0, -1));
}

// Used by: remove.html
function saveRemoveData() {
    var data = "";
    for (c of $('input[type=checkbox]:checked')) data += c.id.slice(4) + ';';
    $('#data').val(data.slice(0, -1));
}
