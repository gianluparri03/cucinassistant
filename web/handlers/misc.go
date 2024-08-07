package handlers

import (
	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetIndex renders /
func GetIndex(c utils.Context) {
	var err error

	if user, err := database.GetUser(c.UID); err == nil {
		utils.RenderPage(c, "misc/home", map[string]any{"Username": user.Username})
		return
	}

	utils.ShowMessage(c, err.Error(), "")
}

// GetInfo renders /info
func GetInfo(c utils.Context) {
	utils.RenderPage(c, "misc/info", map[string]any{
		"Config":      config.Runtime,
		"Version":     config.Version,
		"UsersNumber": database.GetUsersNumber(),
	})
}
