package handlers

import (
	"cucinassistant/config"
	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetIndex renders /
func GetIndex(c *utils.Context) error {
	utils.RenderPage(c, "misc/home", map[string]any{"Username": c.U.Username})
	return nil
}

// GetInfo renders /info
func GetInfo(c *utils.Context) error {
	utils.RenderPage(c, "misc/info", map[string]any{
		"Config":          config.Runtime,
		"VersionCodeName": config.VersionCodeName,
		"VersionNumber":   config.VersionNumber,
		"UsersNumber":     database.GetUsersNumber(),
	})
	return nil
}
