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

// Sets the new value for itemsCount
function setItemsCount(n) {
    itemsCount = n;
}

// Adds the inputs for a new item
function addItem(e, last=false) {
    if (e) { e.preventDefault(); }
    itemsCount++;

    var i = $('.item.hidden').clone();

    if (last) {
        i.appendTo('#new-items');
    } else {
        i.prependTo('#new-items');
    }
    
    i.removeClass("hidden").children().each((ind, inp) => {
        if ((tmpl = $(inp).attr('nametemplate'))) {
            inp.name = tmpl.replace('ID', itemsCount);
        }
    });
}

// Removes the first item (the last one added)
function removeItem(e) {
    if (itemsCount > 1) {
        itemsCount--;
        $('.item').first().remove();
    }

    if (e) { e.preventDefault(); }
} 

// Removes the item whose button called the function.
// It does not decrement itemsCount, so that the absolute order is preserved.
function removeThisItem(e, t) {
    t.parentElement.remove();
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
