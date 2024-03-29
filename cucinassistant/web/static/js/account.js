// Used by: password_change.html, signup.html
function comparePasswords() {
    if ($("#password1").val() != $("#password2").val()) {
        showMessage(false, "Le due password non corrispondono");
        return false;
    } else {
        return true;
    }
}
