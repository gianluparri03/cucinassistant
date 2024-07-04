// Used by: add.html
function updateInputs(e, first=false) {
    if (first || (e.target == $("input[type=text]").get(-1) && e.target.value)) {
        let ne = $(`<div class="new-element">
                        <input class="name" type="text">
                    </div>`);

        $('#new-elements').append(ne);
        ne.on("change", updateInputs);
    } else if (!e.target.value) {
        e.target.remove();
    }
}

// Used by: add.html
function saveAddData() {
    var data = "";
    for (input of $('.new-element > input')) if (input.value) data += input.value + ';';
    $('#data').val(data.slice(0, -1));
}

// Used by: view.html
function clickHandler(e) {
    var line1 = $('<p>').html("Cosa vuoi fare per " + e.target.textContent + "?");
    var line2 = $('<div>');

    function makeFunc(h) { return (e) => { $(e.target).parents('message').remove(); h(); }; }
    var b1 = $('<button class="icon-text">').html('<i class="fas fa-times"></i> Annulla').on('click', makeFunc(() => {}));
    var b2 = $('<button class="icon-text">').html('<i class="fas fa-check"></i> Spunta').on('click', makeFunc(() => {}));
    var b3 = $('<button class="icon-text">').html('<i class="fas fa-edit"></i> Modifica').on('click', makeFunc(() => {}));

    line2.append(b1, b2, b3);

    var content = $('<content>');
    content.append(line1, line2);

    $('body').append($('<message>').append(content));
}
