function showMessage(text, actions) {
    var line1 = $('<p>').html(text);
    var line2 = $('<div>');

    function makeFunc(h) { return (e) => { $(e.target).parents('message').remove(); h(); }; }
    for (let i = 0; i < actions.length; i++) {
        line2.append($("<button class='icon-text'>")
                     .html('<i class="fas ' + actions[i].icon + '"></i> ' + actions[i].name)
                     .on('click', makeFunc(actions[i].handler)));
    }

    let message = $('<message>').append($('<content>').append(line1, line2));
    $('body').append(message);
}

function showAlert(text, handler) {
    showMessage(text, [{icon: "fa-check", name: "Ok", handler: handler}]);
}

function showConfirm(text, handler) {
    showMessage(text, [
        {icon: "fa-times", name: "Annulla", handler: () => {}},
        {icon: "fa-check", name: "Conferma", handler: handler}]);
}


$(window).on('load', function () {
    if (!localStorage.getItem('cookies-accepted')) {
        showAlert(`CucinAssistant utilizza i cookie (non di terze parti) per funzionare. <br>
                   Continuando, si presta il consenso all'utilizzo.`,
                  () => localStorage.setItem('cookies-accepted', true));
    }
});
