package handlers

import (
	"errors"
	"net/url"

	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/web/utils"
)

var (
	ERR_PASSWORDS_DO_NOT_MATCH = errors.New("Le due password non corrispondono")
	MSG_EMAIL_SENT             = "Ti abbiamo inviato una mail. Controlla la casella di posta."
)

// comparePasswords returns an error if the value
// of the fields field1 and field2 (in the request)
// do not match
func comparePasswords(c utils.Context, field1 string, field2 string) error {
	if c.R.FormValue(field1) == c.R.FormValue(field2) {
		return nil
	} else {
		return ERR_PASSWORDS_DO_NOT_MATCH
	}
}

// GetSignUp renders /user/signup
func GetSignUp(c utils.Context) error {
	utils.RenderPage(c, "user/signup", nil)
	return nil
}

// PostSignUp tries to sign up the user whose data is in the request
func PostSignUp(c utils.Context) (err error) {
	// Fetches data from the request
	user := database.User{
		Username: c.R.FormValue("username"),
		Email:    c.R.FormValue("email"),
		Password: c.R.FormValue("password"),
	}

	// Ensures the two passwords are equal
	if err = comparePasswords(c, "password", "password2"); err == nil {
		// Tries to sign it up
		if err = user.SignUp(); err == nil {
			// Sends the welcome email
			go email.SendMail(user.Email, "Registrazione effettuata", "welcome", map[string]any{"Username": user.Username})

			// Saves the UID and redirects to /
			utils.SaveUID(c, user.UID, "Account creato con successo")
		}
	}

	return
}

// GetSignIn renders /user/signin
func GetSignIn(c utils.Context) error {
	utils.RenderPage(c, "user/signin", nil)
	return nil
}

// PostSignIn tries to sign in the user whose data is in the request
func PostSignIn(c utils.Context) (err error) {
	// Fetches data from the request
	user := database.User{
		Username: c.R.FormValue("username"),
		Password: c.R.FormValue("password"),
	}

	// Tries to sign it in
	if err = user.SignIn(); err == nil {
		// Saves the UID and redirects to /
		utils.SaveUID(c, user.UID, "")
	}

	return
}

// PostSignOut drops an user's session
func PostSignOut(c utils.Context) error {
	utils.DropUID(c, "")
	return nil
}

// GetForgotPassword renders /user/forgot_password
func GetForgotPassword(c utils.Context) error {
	utils.RenderPage(c, "user/forgot_pasword", nil)
	return nil
}

// PostForgotPassword sends an email to the user to recover its password
func PostForgotPassword(c utils.Context) (err error) {
	// Fetches the user's email
	userEmail := c.R.FormValue("email")

	// Retrieves the user
	var user database.User
	if user, err = database.GetUserFromEmail(userEmail); err == nil {
		// Tries to generate its token
		var token string
		if token, err = user.GenerateToken(); err == nil {
			// Sends it the email
			go email.SendMail(user.Email, "Recupero password", "reset_password", map[string]any{
				"Username":  user.Username,
				"ResetLink": config.Runtime.BaseURL + "/user/reset_password?token=" + url.QueryEscape(token),
			})

			// Shows the popup
			utils.Show(c, MSG_EMAIL_SENT)
		}
	}

	return
}

// GetResetPassword renders /user/reset_password
func GetResetPassword(c utils.Context) error {
	utils.RenderPage(c, "user/reset_pasword", map[string]any{
		"Token": c.R.URL.Query().Get("token"),
	})
	return nil
}

// PostResetPassword makes an user resets his password
func PostResetPassword(c utils.Context) (err error) {
	// Fetches data
	user := database.User{
		Email: c.R.FormValue("email"),
		Token: c.R.FormValue("token"),
	}

	// Ensures the passwords are equal
	if err = comparePasswords(c, "password-new1", "password-new2"); err == nil {
		// Tries to reset the user's password
		newPassword := c.R.FormValue("password-new1")
		if err = user.ResetPassword(newPassword); err == nil {
			var user database.User
			if user, err = database.GetUser(user.UID); err == nil {
				// Sends the email
				go email.SendMail(user.Email, "Cambio password", "password_change", map[string]any{
					"Username": user.Username,
				})

				// Saves the UID and redirects to /
				utils.SaveUID(c, user.UID, "Password cambiata con successo")
			}
		}
	}

	return
}

// GetSettings renders /user/settings
func GetSettings(c utils.Context) error {
	utils.RenderPage(c, "user/settings", nil)
	return nil
}

// GetChangeUsername renders /user/change_username
func GetChangeUsername(c utils.Context) error {
	utils.RenderPage(c, "user/change_username", map[string]any{"Username": c.U.Username})
	return nil
}

// PostChangeUsername lets an user change its username
func PostChangeUsername(c utils.Context) (err error) {
	// Fetches data
	newUsername := c.R.FormValue("username-new")

	// Tries to change it
	if err = c.U.ChangeUsername(newUsername); err == nil {
		utils.ShowAndRedirect(c, "Nome utente cambiato con successo", "/user/settings")
	}

	return
}

// GetChangeEmail renders /user/change_email
func GetChangeEmail(c utils.Context) error {
	utils.RenderPage(c, "user/change_email", map[string]any{"Email": c.U.Email})
	return nil
}

// PostChangeEmail lets an user change its email
func PostChangeEmail(c utils.Context) (err error) {
	// Fetches data
	newEmail := c.R.FormValue("email-new")

	// Tries to change it
	if err = c.U.ChangeEmail(newEmail); err == nil {
		utils.ShowAndRedirect(c, "Email cambiata con successo", "/user/settings")
	}

	return
}

// GetChangePassword renders /user/change_password
func GetChangePassword(c utils.Context) error {
	utils.RenderPage(c, "user/change_password", nil)
	return nil
}

// PostChangePassword lets an user change its password
func PostChangePassword(c utils.Context) (err error) {
	// Ensures the two passwords are equal
	if err = comparePasswords(c, "password-new-1", "password-new-2"); err == nil {
		newPassword := c.R.FormValue("password-new-1")

		// Tries to change the password
		if err = c.U.ChangePassword(newPassword); err == nil {
			// Sends the email
			go email.SendMail(c.U.Email, "Cambio password", "password_change", map[string]any{
				"Username": c.U.Username,
			})

			// Shows the popup
			utils.ShowAndRedirect(c, "Password cambiata con successo", "/user/settings")
		}
	}

	return
}

// GetDeleteUser1 renders /user/delete_1
func GetDeleteUser1(c utils.Context) error {
	utils.RenderPage(c, "user/delete", map[string]any{"Warning": true})
	return nil
}

// // PostDeleteUser1 sends an email to the user to delete it
func PostDeleteUser1(c utils.Context) (err error) {
	// Tries to generate a new token
	var token string
	if token, err = c.U.GenerateToken(); err == nil {
		// Sends the email
		go email.SendMail(c.U.Email, "Eliminazione account", "delete_confirm", map[string]any{
			"Username":   c.U.Username,
			"DeleteLink": config.Runtime.BaseURL + "/user/delete_2?token=" + url.QueryEscape(token),
		})

		utils.Show(c, MSG_EMAIL_SENT)
	}

	return
}

// GetDeleteUser2 renders /user/delete_2
func GetDeleteUser2(c utils.Context) error {
	utils.RenderPage(c, "user/delete", map[string]any{"Token": c.R.URL.Query().Get("token")})
	return nil
}

// PostDeleteUser2 deletes the user
func PostDeleteUser2(c utils.Context) (err error) {
	// Fetches data from the request
	c.U.Token = c.R.FormValue("token")

	// Tries to delete the user
	if err = c.U.Delete(); err == nil {
		// Sends the goodbye email
		go email.SendMail(c.U.Email, "Eliminazione account", "goodbye", map[string]any{
			"Username": c.U.Username,
		})

		// Drops the session and redirects to /user/signin
		utils.DropUID(c, "Account eliminato con successo")
	}

	return
}
