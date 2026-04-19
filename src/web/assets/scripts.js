$(document).ready(function() {
    htmx.on("htmx:sendError", () => $('#neterror-container').show());
});



// Closes the side bar
function closeSide() {
    $("#side-container").children().remove();
}



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
    
    i.removeClass("hidden").find('input, select').each((ind, inp) => {
        if ((tmpl = $(inp).attr('nametemplate'))) {
            inp.name = tmpl.replace('ID', itemsCount);
        }
    });

    htmx.process(i[0]);
}

// Removes the first item (the last one added)
function removeItem(e) {
    if (itemsCount > 1) {
        itemsCount--;
        $('.item').first().remove();
    }

    if (e) { e.preventDefault(); }
} 



// Swaps two contents inside an area
function swapContent(target) {
    let area = $(target).parents('.swap-area');
    area.find('.pre-swap').remove();
    area.find('.post-swap').removeClass('post-swap');
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



// Lets the user write an arithmetic expression evaluated at the second call.
function calculateQuantity(button, event) {
    event.preventDefault();

    let icon = $(button).children(0);
    let input = $(button).siblings('input');

    if (icon.hasClass('ph-calculator')) {
        icon.removeClass('ph-calculator');
        icon.addClass('ph-equals');

        input.attr('type', 'text');
        input.attr('prevValue', input.val());
        input.prop('onclick', null);
    } else {
        icon.addClass('ph-calculator');
        icon.removeClass('ph-equals');

        let value = input.attr('prevValue');
        try {
            value = eval(input.val());
        } catch {}

        input.attr('type', 'number');
        input.val(value);
    }
}
