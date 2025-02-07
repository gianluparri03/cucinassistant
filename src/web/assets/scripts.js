// Hides the back and the home button
function disableNavigation() {
    $("#back").remove();
    $("#home").css("opacity", "0");
}



// Makes sure the fields password-1 and password-2 match,
// otherwise shows an error.
function comparePasswords(text) {
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
        $('button').before('<div id="unmatching_passwords">'+text+'</div>');

        // And stop the process
        event.preventDefault();
        return false;
    } else {
        return true;
    }
}



let itemsCount;

// Adds the inputs for a new item
function addItem(e) {
    if (e) { e.preventDefault(); } else { itemsCount = 0; }
    itemsCount++;

    $('.item').first().clone().appendTo('#new-items');
    $('.item').last().children().each(function (a, i) {
        i.name = i.attributes['nametemplate'].textContent.replace("ID", itemsCount);
    });
    $('.item').last().removeAttr("hidden");
}

// Removes the last item
function removeItem(e) {
    if (itemsCount > 1) {
        itemsCount--;
        $('.item').last().remove();
    }

    if (e) { e.preventDefault(); }
} 



// Changes the shown message
function swapButtons() {
    $('#pre-delete').remove();
    $('#post-delete').show();
}
