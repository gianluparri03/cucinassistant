// Used by: password_change.html, signup.html
function comparePasswords() {
    if ($("#password1").val() != $("#password2").val()) {
        showAlert("Le due password non corrispondono", () => {});
        return false;
    } else {
        return true;
    }
}
