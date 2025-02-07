package handlers

import (
	"net/url"

	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/langs"
	"cucinassistant/web/utils"
)

// GetSignUp renders /user/signup
func GetSignUp(c *utils.Context) error {
	utils.RenderPage(c, "user/signup", nil)
	return nil
}

// PostSignUp tries to sign up the user whose data is in the request
func PostSignUp(c *utils.Context) (err error) {
	// Fetches data from the request
	username := c.R.FormValue("username")
	email_ := c.R.FormValue("email")
	password := c.R.FormValue("password")

	// Tries to sign it up
	var user database.User
	if user, err = database.SignUp(username, email_, password); err == nil {
		// Saves the email language (based on the current one)
		user.SetEmailLang(c.L)

		// Sends the welcome email
		go email.Welcome.Write(&user, nil).Send()

		// Saves the UID and redirects to /
		utils.SaveUID(c, user.UID, "MSG_USER_CREATED")
	}

	return
}

// GetSignIn renders /user/signin
func GetSignIn(c *utils.Context) error {
	utils.RenderPage(c, "user/signin", nil)
	return nil
}

// PostSignIn tries to sign in the user whose data is in the request
func PostSignIn(c *utils.Context) (err error) {
	// Fetches data from the request
	username := c.R.FormValue("username")
	password := c.R.FormValue("password")

	// Tries to sign it in
	var user database.User
	if user, err = database.SignIn(username, password); err == nil {
		// Saves the UID and redirects to /
		utils.SaveUID(c, user.UID, "")
	}

	return
}

// PostSignOut drops an user's session
func PostSignOut(c *utils.Context) error {
	utils.DropUID(c, "")
	return nil
}

// GetForgotPassword renders /user/forgot_password
func GetForgotPassword(c *utils.Context) error {
	utils.RenderPage(c, "user/forgot_password", nil)
	return nil
}

// PostForgotPassword sends an email to the user to recover its password
func PostForgotPassword(c *utils.Context) (err error) {
	// Fetches the user's email
	userEmail := c.R.FormValue("email")

	// Retrieves the user
	var user database.User
	if user, err = database.GetUser("email", userEmail); err == nil {
		// Tries to generate its token
		var token string
		if token, err = user.GenerateToken(); err == nil {
			// Sends it the email
			go email.ResetPassword.Write(&user, map[string]any{
				"ResetLink": configs.BaseURL + "/user/reset_password?token=" + url.QueryEscape(token),
			}).Send()

			// Shows the popup
			utils.Show(c, "MSG_EMAIL_SENT")
		}
	}

	return
}

// GetResetPassword renders /user/reset_password
func GetResetPassword(c *utils.Context) error {
	utils.RenderPage(c, "user/reset_password", map[string]any{
		"Token": c.R.URL.Query().Get("token"),
	})
	return nil
}

// PostResetPassword makes an user resets his password
func PostResetPassword(c *utils.Context) (err error) {
	// Fetches data
	token := c.R.FormValue("token")
	newPassword := c.R.FormValue("password")
	user := database.User{
		Email: c.R.FormValue("email"),
	}

	// Tries to reset the user's password
	if err = user.ResetPassword(token, newPassword); err == nil {
		var user database.User
		if user, err = database.GetUser("UID", user.UID); err == nil {
			// Sends the email
			go email.PasswordChanged.Write(&user, nil).Send()

			// Saves the UID and redirects to /
			utils.SaveUID(c, user.UID, "MSG_PASSWORD_CHANGED")
		}
	}

	return
}

// GetSettings renders /user/settings
func GetSettings(c *utils.Context) error {
	utils.RenderPage(c, "user/settings", nil)
	return nil
}

// GetChangeUsername renders /user/change_username
func GetChangeUsername(c *utils.Context) error {
	utils.RenderPage(c, "user/change_username", map[string]any{"Username": c.U.Username})
	return nil
}

// PostChangeUsername lets an user change its username
func PostChangeUsername(c *utils.Context) (err error) {
	// Fetches data
	newUsername := c.R.FormValue("username-new")

	// Tries to change it
	if err = c.U.ChangeUsername(newUsername); err == nil {
		utils.ShowAndRedirect(c, "MSG_USERNAME_CHANGED", "/user/settings")
	}

	return
}

// GetChangeEmail renders /user/change_email
func GetChangeEmail(c *utils.Context) error {
	utils.RenderPage(c, "user/change_email", map[string]any{"Email": c.U.Email})
	return nil
}

// PostChangeEmail lets an user change its email
func PostChangeEmail(c *utils.Context) (err error) {
	// Fetches data
	newEmail := c.R.FormValue("email-new")

	// Tries to change it
	if err = c.U.ChangeEmail(newEmail); err == nil {
		utils.ShowAndRedirect(c, "MSG_EMAIL_CHANGED", "/user/settings")
	}

	return
}

// GetChangePassword renders /user/change_password
func GetChangePassword(c *utils.Context) error {
	utils.RenderPage(c, "user/change_password", nil)
	return nil
}

// PostChangePassword lets an user change its password
func PostChangePassword(c *utils.Context) (err error) {
	oldPassword := c.R.FormValue("old-password")
	newPassword := c.R.FormValue("new-password")

	// Tries to change the password
	if err = c.U.ChangePassword(oldPassword, newPassword); err == nil {
		// Sends the email
		go email.PasswordChanged.Write(c.U, nil).Send()

		// Shows the popup
		utils.ShowAndRedirect(c, "MSG_PASSWORD_CHANGED", "/user/settings")
	}

	return
}

// GetSetEmailLang renders /user/set_email_lang
func GetSetEmailLang(c *utils.Context) error {
	utils.RenderPage(c, "user/set_email_lang", map[string]any{
		"Langs":   langs.Available,
		"Current": c.U.EmailLang,
	})
	return nil
}

// PostSetEmailLang lets an user change its username
func PostSetEmailLang(c *utils.Context) (err error) {
	// Fetches data
	lang := c.R.FormValue("lang")

	// Tries to change it
	if err = c.U.SetEmailLang(lang); err == nil {
		utils.ShowAndRedirect(c, "MSG_LANG_CHANGED", "/user/settings")
	}

	return
}

// GetDeleteUser1 renders /user/delete_1
func GetDeleteUser1(c *utils.Context) error {
	utils.RenderPage(c, "user/delete", map[string]any{"Warning": true})
	return nil
}

// // PostDeleteUser1 sends an email to the user to delete it
func PostDeleteUser1(c *utils.Context) (err error) {
	// Tries to generate a new token
	var token string
	if token, err = c.U.GenerateToken(); err == nil {
		// Sends the email
		go email.DeleteConfirm.Write(c.U, map[string]any{
			"DeleteLink": configs.BaseURL + "/user/delete_2?token=" + url.QueryEscape(token),
		}).Send()

		utils.Show(c, "MSG_EMAIL_SENT")
	}

	return
}

// GetDeleteUser2 renders /user/delete_2
func GetDeleteUser2(c *utils.Context) error {
	utils.RenderPage(c, "user/delete", map[string]any{"Token": c.R.URL.Query().Get("token")})
	return nil
}

// PostDeleteUser2 deletes the user
func PostDeleteUser2(c *utils.Context) (err error) {
	// Fetches data from the request
	token := c.R.FormValue("token")

	// Tries to delete the user
	if err = c.U.Delete(token); err == nil {
		// Sends the goodbye email
		go email.Goodbye.Write(c.U, nil).Send()

		// Drops the session and redirects to /user/signin
		utils.DropUID(c, "MSG_USER_DELETED")
	}

	return
}
