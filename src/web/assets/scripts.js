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



// ( used by: /shopping_list/append )
let entriesN = 0;

// Adds the inputs for a new entry
// ( used by: /shopping_list/append )
function addEntry(e) {
    entriesN++;

    $('#new-entries').prepend(`<div class="entry">
        <input class="name" type="text" name="entry-${entriesN}-name" placeholder="Nome">
    </div>`);

    if (e) { e.preventDefault(); }
}

// Removes the first entry's inputs
// ( used by: /shopping_list/append )
function removeEntry(e) {
    if (entriesN > 1) {
        entriesN--;
        $('.entry')[0].remove();
    }

    if (e) { e.preventDefault(); }
} 



// ( used by: /storage/{sid}/add )
let articlesN = 0;

// Adds the inputs for a new article
// ( used by: /storage/{sid}/add )
function addArticle(e) {
    articlesN++;

    $('#new-articles').prepend(`<div class="article">
        <input class="name" name="article-${articlesN}-name" type="text" placeholder="Nome">
        <input class="expiration" name="article-${articlesN}-expiration" type="date" placeholder="Scadenza">
        <input class="quantity" name="article-${articlesN}-quantity" type="number" min="0" step="1" placeholder="Quantità">
    </div>`);

    if (e) { e.preventDefault(); }
}

// Removes the first article's inputs
// ( used by: /storage/{sid}/add )
function removeArticle(e) {
    if (articlesN > 1) {
        articlesN--;
        $('.article')[0].remove();
    }

    if (e) { e.preventDefault(); }
} 


// Changes the shown message
// ( used by: /storage/{SID}/edit )
function swapButtons() {
    $('#pre-delete').remove();
    $('#post-delete').show();
}
