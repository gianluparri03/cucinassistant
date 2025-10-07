// Makes sure the fields password-1 and password-2 match,
// otherwise shows an error.
function comparePasswords(event) {
    // Remove the previous warnin
    $('#unmatching-passwords').addClass('hidden');

    // Gets the fields
    let p1 = $('#password-1');
    let p2 = $('#password-2');
    
    // Check if the passwords match
    if (p1.val() != p2.val()) {
        // If not, resets them
        p1.val("");
        p2.val("");

        // Shows the warning
        $('#unmatching-passwords').removeClass('hidden');

        // And stop the process
        event.preventDefault();
    }
}



let itemsCount;

function setItemsCount(n) { itemsCount = n; }

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


let locale;

function setLocale(l) { locale = l; }

// Sets the menu's name and days according to the `from` and `to` inputs
function calcDates(event) {
    const nameDateOpt = {month: 'short', day: 'numeric'};
    const dayDateOpt = {weekday: 'long', month: 'short', day: 'numeric'};

    $('#invalid-dates').addClass('hidden');

    var from = new Date($('#days-from').val());
    var to = new Date($('#days-to').val());
    if (isNaN(from.getTime()) || isNaN(to.getTime()) || to < from) {
        $('#invalid-dates').removeClass('hidden');
        event.preventDefault();
        return;
    }

    $('#menu-name').val(
        from.toLocaleDateString(locale, nameDateOpt) + ' - ' +
        to.toLocaleDateString(locale, nameDateOpt)
    );

    var days = '';
    for (var d = from; d <= to; d.setDate(d.getDate()+1)) {
        if (days) days += '\n';
        days += d.toLocaleDateString(locale, dayDateOpt);
    }
    $('#menu-days').val(days);
}
