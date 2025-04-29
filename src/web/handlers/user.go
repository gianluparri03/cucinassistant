package handlers

import (
	"cucinassistant/configs"
	"cucinassistant/database"
	"cucinassistant/email"
	"cucinassistant/langs"
	"cucinassistant/web/components"
	"cucinassistant/web/utils"
)

func GetUserChangeEmail(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserChangeEmail(c.U.Email))
	return
}

func PostUserChangeEmail(c *utils.Context) (err error) {
	newEmail := c.R.FormValue("email-new")

	if err = c.U.ChangeEmail(newEmail); err == nil {
		utils.ShowMessage(c, langs.STR_EMAIL_CHANGED, "/user/settings")
	}

	return
}

func GetUserChangeEmailSettings(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserChangeEmailSettings(c.U, langs.Available))
	return
}

func PostUserChangeEmailSettings(c *utils.Context) (err error) {
	lang := c.R.FormValue("lang")
	newsletter := c.R.FormValue("newsletter") == "on"

	if err = c.U.ChangeEmailSettings(lang, newsletter); err == nil {
		utils.ShowMessage(c, langs.STR_SETTINGS_SAVED, "/user/settings")
	}

	return
}

func GetUserChangePassword(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserChangePassword())
	return
}

func PostUserChangePassword(c *utils.Context) (err error) {
	oldPassword := c.R.FormValue("old-password")
	newPassword := c.R.FormValue("new-password")

	if err = c.U.ChangePassword(oldPassword, newPassword); err == nil {
		go email.PasswordChanged.Write(c.U, "").Send()
		utils.ShowMessage(c, langs.STR_PASSWORD_CHANGED, "/user/settings")
	}

	return
}

func GetUserChangeUsername(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserChangeUsername(c.U.Username))
	return
}

func PostUserChangeUsername(c *utils.Context) (err error) {
	newUsername := c.R.FormValue("username-new")

	if err = c.U.ChangeUsername(newUsername); err == nil {
		utils.ShowMessage(c, langs.STR_USERNAME_CHANGED, "/user/settings")
	}

	return
}

func GetUserDelete1(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserDelete(true, ""))
	return
}

func PostUserDelete1(c *utils.Context) (err error) {
	var token string

	if token, err = c.U.GenerateToken(); err == nil {
		go email.DeleteConfirm.Write(
			c.U,
			configs.BaseURL+"/user/delete_2?token="+token,
		).Send()
		utils.ShowMessage(c, langs.STR_EMAIL_SENT, "")
	}

	return
}

func GetUserDelete2(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserDelete(false, c.R.URL.Query().Get("token")))
	return
}

func PostUserDelete2(c *utils.Context) (err error) {
	token := c.R.FormValue("token")

	if err = c.U.Delete(token); err == nil {
		go email.Goodbye.Write(c.U, "").Send()
		utils.DropUID(c, langs.STR_USER_DELETED)
	}

	return
}

func GetForgotPassword(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserForgotPassword())
	return
}

func PostForgotPassword(c *utils.Context) (err error) {
	var user database.User
	var token string

	userEmail := c.R.FormValue("email")
	if user, err = database.GetUser("email", userEmail); err == nil {
		if token, err = user.GenerateToken(); err == nil {
			go email.ResetPassword.Write(
				&user,
				configs.BaseURL+"/user/reset_password?token="+token,
			).Send()
			utils.ShowMessage(c, langs.STR_EMAIL_SENT, "")
		}
	}

	return
}

func GetResetPassword(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserResetPassword(c.R.URL.Query().Get("token")))
	return
}

func PostResetPassword(c *utils.Context) (err error) {
	token := c.R.FormValue("token")
	newPassword := c.R.FormValue("password")
	user := database.User{
		Email: c.R.FormValue("email"),
	}

	if err = user.ResetPassword(token, newPassword); err == nil {
		if user, err = database.GetUser("UID", user.UID); err == nil {
			go email.PasswordChanged.Write(&user, "").Send()
			utils.SaveUID(c, user.UID, langs.STR_PASSWORD_CHANGED)
		}
	}

	return
}

func GetUserSettings(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserSettings(configs.SupportEmail))
	return
}

func GetUserSignIn(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserSignIn())
	return
}

func PostUserSignIn(c *utils.Context) (err error) {
	var user database.User

	username := c.R.FormValue("username")
	password := c.R.FormValue("password")
	if user, err = database.SignIn(username, password); err == nil {
		utils.SaveUID(c, user.UID, langs.STR_NONE)
	}

	return
}

func PostUserSignOut(c *utils.Context) (err error) {
	utils.DropUID(c, langs.STR_NONE)
	return
}

func GetUserSignUp(c *utils.Context) (err error) {
	utils.RenderComponent(c, components.UserSignUp())
	return
}

func PostUserSignUp(c *utils.Context) (err error) {
	var user database.User

	username := c.R.FormValue("username")
	email_ := c.R.FormValue("email")
	password := c.R.FormValue("password")

	if user, err = database.SignUp(username, email_, password); err == nil {
		user.ChangeEmailSettings(c.L, true)
		go email.Welcome.Write(&user, "").Send()
		utils.SaveUID(c, user.UID, langs.STR_USER_CREATED)
	}

	return
}
