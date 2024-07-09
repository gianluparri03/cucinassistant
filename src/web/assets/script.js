function showMessage(text, actions) {
    var line1 = $('<p>').html(text);
    var line2 = $('<div>').attr('hx-push-url', 'true');

    for (const action of actions) {
        let button = $("<button class='icon-text'>")
                     .html(`<i class="fas ${action.icon}"></i> ${action.label}`);

        if (action.target) {
            button.attr('hx-' + (action.method || 'get'), action.target);
            htmx.process(button.get(0));
        } else if (action.handler) {
            button.on('click', action.handler);
        }

        line2.append(button);
    }

    $('body').append(
        $('<message>').append(
            $('<content>').append(line1, line2)
        ).on('click', e => $(e.target).parents('message').remove())
    );
}

function showAlert(text, handler) {
    showMessage(text, [{label: "Ok", icon: "fa-check", handler: handler}]);
}


function clickShoppingItem(e) {
    e.preventDefault();

    let sp = e.target.classList.contains('shopping-item') ? $(e.target) : $(e.target).parent();
    let label = $(sp).find('label');
    let id = label.attr('for').replace('entry-', '');
    let name = label.text();

    let check = {label: 'Spunta', icon: 'fa-check', target: `/spesa/${id}/spunta/`, method: 'post'}
    let uncheck = {label: 'Riaggiungi', icon: 'fa-undo', target: `/spesa/${id}/spunta/`, method: 'post'}

    showMessage(`Cosa vuoi fare per ${name}?`,
        [
            {label: 'Annulla', icon: 'fa-times'},
            Math.floor(Math.random() * 10) % 2 ? check : uncheck,
            {label: 'Modifica', icon: 'fa-edit', target: `/spesa/${id}/modifica/`}
        ]
    );
}

function updateShoppingInputs(e, first=false) {
    if (first || (e.target == $("input[type=text]").get(-1) && e.target.value)) {
        let ne = $(`<div class="new-element">
                        <input class="name" type="text" name="name-${inputNo}">
                    </div>`);

        $('#new-elements').append(ne);
        ne.on("change", updateShoppingInputs);
    } else if (!e.target.value) {
        e.target.remove();
    }
}


function updateStorageInputs(e, first=false) {
    let target = e ? $(e.target).parent() : null;

    if (first || (target.get(0) == $(".storage-item").get(-1) && target.find('input.name').val())) {
        let ne = $(`<div class="storage-item">
                        <input class="name" name="name-${inputNo}" type="text" placeholder="Nome">
                        <input class="expiration" name="expiration-${inputNo} "type="date" placeholder="Scadenza">
                        <input class="quantity" name="quantity-${inputNo} "type="number" min="0" step="1" placeholder="Quantità">
                    </div>`);

        $('#new-elements').append(ne);
        ne.on("change", updateStorageInputs);
    } else if (!target.find('input.name').val()) {
        target.remove();
    }
}
