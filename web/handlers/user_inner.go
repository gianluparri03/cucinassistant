package handlers

import (
	"net/url"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/web/utils"
)


// retrieveUser returns an User struct with all the data
// found in the request
func retrieveUser(c utils.Context) database.User {
    return database.User{
		Username: c.R.FormValue("username"),
		Email:    c.R.FormValue("email"),
		Password: c.R.FormValue("password"),
		Token: c.R.FormValue("token"),
    }
}

// signUpUser tries to sign up the user whose data is in the request
func signUpUser(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
    user := retrieveUser(c)
	if user.Password != c.R.FormValue("password2") {
		utils.ShowError(c, "Le due password non corrispondono")
		return
	}

	// Tries to sign up the user
	if err := user.SignUp(); err != nil {
		utils.ShowError(c, err.Error())
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// signInUser tries to sign in the user whose data is in the request
func signInUser(c utils.Context) {
	// Fetches data
    user := retrieveUser(c)

	// Tries to sign up the user
	if err := user.SignIn(); err != nil {
		utils.ShowError(c, err.Error())
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// signOutUser drops an user's session
func signOutUser(c utils.Context) {
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.Redirect(c, "/user?action=signin")
}

// forgotPassword sends an email to the user to recover its password
func forgotPassword(c utils.Context) {
	// Tries to get the user's email
    userEmail := c.R.FormValue("email")
	if user, err := database.GetUserFromEmail(userEmail); err == nil {
		// Tries to generate its token
		if token, err := user.GenerateToken(); err == nil {
			go email.SendMail(user.Email, "Recupero password", "reset_password", map[string]any{
				"Username":  user.Username,
				"ResetLink": config.Runtime.BaseURL + "/user?action=reset_password&token=" + url.QueryEscape(token),
			})
		}
	}

	utils.ShowError(c, "Ti abbiamo inviato un email. Controlla la casella di posta")
}

// resetPassword makes an user resets his password
func resetPassword(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
    user := retrieveUser(c)
	if user.Password != c.R.FormValue("password2") {
		utils.ShowError(c, "Le due password non corrispondono")
		return
	}

	// Tries to reset the user's password
	if err := user.ResetPassword(); err != nil {
		utils.ShowError(c, err.Error())
	}

    // Saves the session and then redirects it to the homepage
    c.S.Values["UID"] = user.UID
    utils.SaveSession(c)
    utils.Redirect(c, "/")
}
