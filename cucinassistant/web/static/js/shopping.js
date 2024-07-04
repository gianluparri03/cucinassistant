// Used by: view.html
function clickHandler(e) {
    let check = {icon: "fa-check", name: "Spunta", handler: () => alert('unimplemented')};
    let uncheck = {icon: "fa-undo", name: "Riaggiungi", handler: () => alert('unimplemented')};

    showMessage("Cosa vuoi fare per " + e.target.textContent + "?",
        [
            {icon: "fa-times", name: "Annulla", handler: () => {}},
            Math.floor(Math.random() * 10) % 2 ? check : uncheck,
            {icon: "fa-edit", name: "Modifica", handler: () => location.href = '/spesa/modifica/' + e.target.htmlFor},
        ]
    );
}

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
