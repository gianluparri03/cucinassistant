function showMessage(canDecline, text, handler=(() => {})) {
    var line1 = $('<p>').html(text);
    var line2 = $('<div>');

    function makeFunc(h) { return (e) => { $(e.target).parents('message').remove(); h(); }; }
    var decline = $('<button class="icon">').html('<i class="fas fa-times">').on('click', makeFunc(() => {}));
    var accept = $('<button class="icon">').html('<i class="fas fa-check">').on('click', makeFunc(handler));

    if (canDecline) { line2.append(decline); }
    line2.append(accept);

    var content = $('<content>');
    content.append(line1);
    content.append(line2);

    $('body').append($('<message>').append(content));
    accept.focus();
}
