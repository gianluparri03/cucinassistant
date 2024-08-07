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
		UID:      c.UID,
		Username: c.R.FormValue("username"),
		Email:    c.R.FormValue("email"),
		Password: c.R.FormValue("password"),
		Token:    c.R.FormValue("token"),
	}
}

// GetSignUp renders /user/signup
func GetSignUp(c utils.Context) {
	utils.RenderPage(c, "user/signup", nil)
}

// PostSignUp tries to sign up the user whose data is in the request
func PostSignUp(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	if user.Password != c.R.FormValue("password2") {
		utils.ShowMessage(c, "Le due password non corrispondono", "")
		return
	}

	// Tries to sign up the user
	if err := user.SignUp(); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	// Sends the welcome email
	go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.ShowMessage(c, "Account creato con successo", "/")
}

// GetSignIn renders /user/signin
func GetSignIn(c utils.Context) {
	utils.RenderPage(c, "user/signin", nil)
}

// PostSignIn tries to sign in the user whose data is in the request
func PostSignIn(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)

	// Tries to sign up the user
	if err := user.SignIn(); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// PostSignOut drops an user's session
func PostSignOut(c utils.Context) {
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.Redirect(c, "/user/signin")
}

// GetForgotPassword renders /user/forgot_password
func GetForgotPassword(c utils.Context) {
	utils.RenderPage(c, "user/forgot_pasword", nil)
}

// PostForgotPassword sends an email to the user to recover its password
func PostForgotPassword(c utils.Context) {
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

	utils.ShowMessage(c, "Ti abbiamo inviato un email. Controlla la casella di posta", "")
}

// GetResetPassword renders /user/reset_password
func GetResetPassword(c utils.Context) {
	utils.RenderPage(c, "user/reset_pasword", map[string]any{
		"Token": c.R.URL.Query().Get("token"),
	})
}

// PostResetPassword makes an user resets his password
func PostResetPassword(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	newPassword := c.R.FormValue("password-new1")
	if newPassword != c.R.FormValue("password-new2") {
		utils.ShowMessage(c, "Le due password non corrispondono", "")
		return
	}

	// Tries to reset the user's password
	if err := user.ResetPassword(newPassword); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	// Saves the session and then redirects it to the homepage
	c.S.Values["UID"] = user.UID
	utils.SaveSession(c)
	utils.Redirect(c, "/")
}

// GetSettings renders /user/settings
func GetSettings(c utils.Context) {
	utils.RenderPage(c, "user/settings", nil)
}

// GetChangeUsername renders /user/change_username
func GetChangeUsername(c utils.Context) {
	var err error

	if user, err := database.GetUser(c.UID); err == nil {
		utils.RenderPage(c, "user/change_username", map[string]any{"Username": user.Username})
		return
	}

	utils.ShowMessage(c, err.Error(), "")
}

// PostChangeUsername lets an user change its username
func PostChangeUsername(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)
	newUsername := c.R.FormValue("username-new")

	// Tries to change it
	if err := user.ChangeUsername(newUsername); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	utils.ShowMessage(c, "Nome utente cambiato con successo", "/user/settings")
}

// GetChangeEmail renders /user/change_email
func GetChangeEmail(c utils.Context) {
	var err error

	if user, err := database.GetUser(c.UID); err == nil {
		utils.RenderPage(c, "user/change_email", map[string]any{"Email": user.Email})
		return
	}

	utils.ShowMessage(c, err.Error(), "")
}

// PostChangeEmail lets an user change its email
func PostChangeEmail(c utils.Context) {
	// Fetches data
	user := retrieveUser(c)
	newEmail := c.R.FormValue("email-new")

	// Tries to change it
	if err := user.ChangeEmail(newEmail); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	utils.ShowMessage(c, "Email cambiata con successo", "/user/settings")
}

// GetChangePassword renders /user/change_password
func GetChangePassword(c utils.Context) {
	utils.RenderPage(c, "user/change_password", nil)
}

// PostChangePassword lets an user change its password
func PostChangePassword(c utils.Context) {
	// Fetches data, ensuring the two passwords are equal
	user := retrieveUser(c)
	newPassword := c.R.FormValue("password-new1")
	if newPassword != c.R.FormValue("password-new2") {
		utils.ShowMessage(c, "Le due password non corrispondono", "")
		return
	}

	// Tries to change the user's password
	if err := user.ChangePassword(newPassword); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	// Sends the email
	user, _ = database.GetUser(user.UID)
	go email.SendMail(user.Email, "Cambio password", "password_change", map[string]any{
		"Username": user.Username,
	})

	utils.ShowMessage(c, "Password cambiata con successo", "/user/settings")
}

// GetDeleteUser1 renders /user/delete_1
func GetDeleteUser1(c utils.Context) {
	utils.RenderPage(c, "user/delete", map[string]any{"Warning": true})
}

// PostDeleteUser1 sends an email to the user to delete it
func PostDeleteUser1(c utils.Context) {
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

	utils.ShowMessage(c, "Ti abbiamo inviato un email. Controlla la casella di posta", "")
}

// GetDeleteUser2 renders /user/delete_2
func GetDeleteUser2(c utils.Context) {
	utils.RenderPage(c, "user/delete", map[string]any{"Token": c.R.URL.Query().Get("token")})
}

// PostDeleteUser2 deletes the user
func PostDeleteUser2(c utils.Context) {
	// Fetches the UID and the token
	user := retrieveUser(c)

	// Fetches the email and the username
	if user_, err := database.GetUser(c.UID); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	} else {
		user.Email = user_.Email
		user.Username = user_.Username
	}

	// Tries to delete the user
	if err := user.Delete(); err != nil {
		utils.ShowMessage(c, err.Error(), "")
		return
	}

	// Sends the goodbye email
	go email.SendMail(user.Email, "Eliminazione account", "goodbye", map[string]any{
		"Username": user.Username,
	})

	// Logs it out, and shows the goodbye message
	delete(c.S.Values, "UID")
	utils.SaveSession(c)
	utils.ShowMessage(c, "Account eliminato con successo", "/user/signin")
}
