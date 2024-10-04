package handlers

import (
	"github.com/gorilla/mux"
	"strconv"

	"cucinassistant/database"
	"cucinassistant/web/utils"
)

// GetSID returns the SID written in the url
func GetSID(c *utils.Context) (SID int, err error) {
	// Fetches the SID from the url, then converts it to an int
	SID, err = strconv.Atoi(mux.Vars(c.R)["SID"])
	if err != nil {
		err = database.ERR_SECTION_NOT_FOUND
	}

	return
}

// GetSections renders /storage
func GetSections(c *utils.Context) (err error) {
	var list []database.Section
	if list, err = c.U.GetSections(); err == nil {
		utils.RenderPage(c, "storage/dashboard", map[string]any{"List": list})
	}

	return
}

// GetNewSection renders /storage/new
func GetNewSection(c *utils.Context) (err error) {
	utils.RenderPage(c, "storage/new_section", nil)
	return
}

// PostNewSection tries to create a new section
func PostNewSection(c *utils.Context) (err error) {
	var s database.Section
	if s, err = c.U.NewSection(c.R.FormValue("name")); err == nil {
		utils.ShowAndRedirect(c, "Sezione creata con successo", "/storage/"+strconv.Itoa(s.SID))
	}

	return
}

// GetSection renders /storage/{SID}
func GetSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = GetSID(c); err == nil {
		var section database.Section
		if section, err = c.U.GetSection(SID, true); err == nil {
			utils.RenderPage(c, "storage/view", map[string]any{"Section": section, "Filter": ""})
		}
	}

	return
}

// GetEditSection renders /storage/{SID}/edit
func GetEditSection(c *utils.Context) (err error) {
	// Retrieves the SID
	var SID int
	if SID, err = GetSID(c); err == nil {
		var section database.Section
		if section, err = c.U.GetSection(SID, false); err == nil {
			utils.RenderPage(c, "storage/edit_section", map[string]any{"Section": section})
		}
	}

	return
}

// PostEditSection tries to change a section's name
func PostEditSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = GetSID(c); err == nil {
		if err = c.U.EditSection(SID, c.R.FormValue("name")); err == nil {
			utils.ShowAndRedirect(c, "Modifiche applicate con successo", "/storage/"+strconv.Itoa(SID))
		}
	}

	return
}

// PostDeleteSection tries to delete a section
func PostDeleteSection(c *utils.Context) (err error) {
	var SID int
	if SID, err = GetSID(c); err == nil {
		if err = c.U.DeleteSection(SID); err == nil {
			utils.ShowAndRedirect(c, "Sezione eliminata con successo", "/storage")
		}
	}

	return
}
