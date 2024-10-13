let inputN; // ( used by updateArticles and updateEntries )

// Makes sure the fields password-1 and password-2 match,
// otherwise shows an error.
// ( used by /user/sign_up, /user/change_password and /user/reset_password )
function comparePasswords() {
    // Remove the previous warning
    $('#unmatching_passwords').remove();

    // Gets the fields
    let p1 = $('#password-1');
    let p2 = $('#password-2');
    
    // Check if the passwords match
    if (p1.val() != p2.val()) {
        // If not, resets them
        p1.val("");
        p2.val("");

        // Shows the warning
        $('button').before('<div id="unmatching_passwords">Le due password non corrispondono</div>');

        // And stop the process
        event.preventDefault();
        return false;
    } else {
        return true;
    }
}

// If the first article's name has been filled, it adds a new one
// on top of it. If an article's name is cleared, the article will be dropped.
// It always keep an input visible.
// ( used by: /storage/{SID}/add )
function updateArticles(e, bootstrap=false) {
    if (bootstrap || (e.target == $("input[type=text]").get(0) && e.target.value)) {
        let na = $(`<div class="article">
                        <input class="name" name="article-${inputN}-name" type="text" placeholder="Nome">
                        <input class="expiration" name="article-${inputN}-expiration" type="date" placeholder="Scadenza">
                        <input class="quantity" name="article-${inputN}-quantity" type="number" min="0" step="1" placeholder="Quantità">
                    </div>`);

        $('#new-articles').prepend(na);
        na.on("change", updateArticles);
        inputN++;
    } else if (e.target.type == 'text' && !e.target.value) {
        e.target.parentNode.remove();
    }
}

// If the first input has been filled, it adds a new one
// on top of it. If one input is cleared, it will delete it.
// It always keep an input visible.
// ( used by: /shopping_list/append )
function updateEntries(e, bootstrap=false) {
    if (bootstrap || (e.target == $("input[type=text]").get(0) && e.target.value)) {
        let ne = $(`<div class="new-element">
                        <input class="name" type="text" name="entry-${inputN}-name">
                    </div>`);

        $('#new-elements').prepend(ne);
        ne.on("change", updateEntries);
        inputN++;
    } else if (!e.target.value) {
        e.target.remove();
    }
}

// Changes the shown message
// ( used by: /storage/{SID}/edit )
function swapButtons() {
    $('#pre-delete').remove();
    $('#post-delete').show();
}
