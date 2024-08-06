package web

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
		UID:      c.UID,
		Username: c.R.FormValue("username"),
		Email:    c.R.FormValue("email"),
		Password: c.R.FormValue("password"),
		Token:    c.R.FormValue("token"),
	}
}

// signUp tries to sign up the user whose data is in the request
func signUp(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	if user.Password != c.R.FormValue("password2") {
		utils.ShowError(c, "Le due password non corrispondono", false)
		return
	}

	// Tries to sign up the user
	if err := user.SignUp(); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.ShowError(c, "Account creato con successo", true)
}

// signIn tries to sign in the user whose data is in the request
func signIn(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)

	// Tries to sign up the user
	if err := user.SignIn(); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// signOut drops an user's session
func signOut(c utils.Context) {
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.Redirect(c, "/user/signin")
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
				"ResetLink": config.Runtime.BaseURL + "/user/reset_password?token=" + url.QueryEscape(token),
			})
		}
	}

	utils.ShowError(c, "Ti abbiamo inviato un email. Controlla la casella di posta", true)
}

// resetPassword makes an user resets his password
func resetPassword(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	newPassword := c.R.FormValue("password-new1")
	if newPassword != c.R.FormValue("password-new2") {
		utils.ShowError(c, "Le due password non corrispondono", false)
		return
	}

	// Tries to reset the user's password
	if err := user.ResetPassword(newPassword); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// changeUsername lets an user change its username
func changeUsername(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)
	newUsername := c.R.FormValue("username-new")

	// Tries to change it
	if err := user.ChangeUsername(newUsername); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	utils.ShowError(c, "Nome utente cambiato con successo", true)
}

// changeEmail lets an user change its email
func changeEmail(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)
	newEmail := c.R.FormValue("email-new")

	// Tries to change it
	if err := user.ChangeEmail(newEmail); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	utils.ShowError(c, "Nome utente cambiato con successo", true)
}

// changePassword lets an user change its password
func changePassword(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	newPassword := c.R.FormValue("password-new1")
	if newPassword != c.R.FormValue("password-new2") {
		utils.ShowError(c, "Le due password non corrispondono", false)
		return
	}

	// Tries to change the user's password
	if err := user.ChangePassword(newPassword); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	// Sends the email
	user, _ = database.GetUser(user.UID)
	go email.SendMail(user.Email, "Cambio password", "password_change", map[string]any{
		"Username":   user.Username,
	})

	utils.ShowError(c, "Password cambiata con successo", true)
}

// deleteUser1 sends an email to the user to delete it
func deleteUser1(c utils.Context) {
	// Gets the user's data
	if user, err := database.GetUser(c.UID); err == nil {
		// Tries to generate a new token
		if token, err := user.GenerateToken(); err == nil {
			// Sends the email
			go email.SendMail(user.Email, "Eliminazione account", "delete_confirm", map[string]any{
				"Username":   user.Username,
				"DeleteLink": config.Runtime.BaseURL + "/user/delete_2?token=" + url.QueryEscape(token),
			})
		}
	}

	utils.ShowError(c, "Ti abbiamo inviato un email. Controlla la casella di posta", true)
}

// deleteUser2 deletes the user
func deleteUser2(c utils.Context) {
	// Fetches the UID and the token
	user := retrieveUser(c)

	// Fetches the email and the username
	if user_, err := database.GetUser(c.UID); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	} else {
		user.Email = user_.Email
		user.Username = user_.Username
	}

	// Tries to delete the user
	if err := user.Delete(); err != nil {
		utils.ShowError(c, err.Error(), false)
		return
	}

	// Sends the goodbye email
	go email.SendMail(user.Email, "Eliminazione account", "goodbye", map[string]any{
		"Username": user.Username,
	})

	// Logs it out, and shows the goodbye message
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.ShowError(c, "Account eliminato con successo", true)
}
