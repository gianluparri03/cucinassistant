function showMessage(canDecline, text, handler=(() => {})) {
    var line1 = $('<p>').html(text);
    var line2 = $('<div>');

    function makeFunc(h) { return (e) => { $(e.target).parents('message').remove(); h(); }; }
    var decline = $('<button title="Annulla" class="icon">').html('<i class="fas fa-times">').on('click', makeFunc(() => {}));
    var accept = $('<button title="Conferma" class="icon">').html('<i class="fas fa-check">').on('click', makeFunc(handler));

    if (canDecline) { line2.append(decline); }
    line2.append(accept);

    var content = $('<content>');
    content.append(line1);
    content.append(line2);

    $('body').append($('<message>').append(content));
    accept.focus();
}

$(window).on('load', function () {
    if (!localStorage.getItem('cookies-accepted')) {
        var cookiesMessage = `CucinAssistant utilizza i cookie (non di terze parti) per funzionare. <br>
                              Continuando, si presta il consenso all'utilizzo.`;

        showMessage(false, cookiesMessage, function() {
            localStorage.setItem('cookies-accepted', true);

            var welcomeMessage = `Benvenuto su <b>CucinAssistant</b>! <br>
                                  Vuoi aprire la guida all'utilizzo? <br> <br>
                                  Potrai comunque aprirla in seguito cliccando il link sulla barra in basso.`;
            showMessage(true, welcomeMessage, () => window.open("/guida"));
        });
    }

    if ($('form[autocomplete=off]').length && $("body").height() > $(window).height()) {
        $('.buttons-line').clone().appendTo('form');
    }
})
