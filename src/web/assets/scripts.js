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

    $('.item.hidden')
        .clone()
        .prependTo('#new-items')
        .removeClass("hidden")
        .children().each((ind, inp) => {
            inp.name = $(inp).attr('nametemplate').replace('ID', itemsCount);
        });
}

// Removes the first item (the last on added)
function removeItem(e) {
    if (itemsCount > 1) {
        itemsCount--;
        $('.item').first().remove();
    }

    if (e) { e.preventDefault(); }
} 



// Swaps two contents
function swapContent() {
    $('.pre-swap').remove();
    $('.post-swap').removeClass("hidden");
}


// Prints a recipe from the share page
function printRecipe() {
	function inner() {
		window.print();
		document.removeEventListener('htmx:afterSwap', inner);
	}

	document.addEventListener('htmx:afterSwap', inner);
	$("#back").click();
}

// Converts the expiration input to an actual date input
function activateExpirationInput(i) {
    i.type = 'date';
    i.showPicker();
} 

// Converts the quantity input to an actual number input
function activateQuantityInput(i) {
    i.type = 'number';
} 

// Converts the section input to an actual selection input
function activateSectionInput(i) {
    var s = i.nextElementSibling;
    i.remove();
    s.classList.remove('hidden');
    s.showPicker();
}
