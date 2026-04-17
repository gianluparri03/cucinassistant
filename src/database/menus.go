package database

import (
	"database/sql"
	"github.com/lib/pq"
)

// Menu is a collection of meals, divided into days
type Menu struct {
	// MID is the Menu ID
	MID int

	// Name is the name of the menu
	Name string

	// Days is the list of days of which the menu is composed
	Days []Day
}

// Day is a component of a menu, and contains a name and a list of meals
type Day struct {
	// MID is the Menu's ID
	MID int

	// Name is the day's name
	Name string

	// Position is used to identify the day
	Position int

	// Meals contains the meals of the day
	Meals []string
}

// Menus is used to manage all the menus
type Menus struct {
	uid int
}

// Menus returns the menus manager for the user
func (u User) Menus() Menus {
	return Menus{uid: u.UID}
}

// AddDay adds a new day in a menu
func (m Menus) AddDay(MID int, name string) error {
	// Gets the menu
	_, err := m.GetOne(MID)
	if err != nil {
		return err
	}

	// Adds the new day
	_, err = db.Exec(`INSERT INTO days (mid, position, name, meals) SELECT $1, max(position)+1, $2, $3 FROM days WHERE mid=$1;`, MID, name, pq.Array([]string{}))
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// Delete deletes a menu
func (m Menus) Delete(MID int) error {
	// Deletes the menu
	res, err := db.Exec(`DELETE FROM menus WHERE uid=$1 AND mid=$2;`, m.uid, MID)
	if err != nil {
		return ERR_UNKNOWN
	} else if ra, _ := res.RowsAffected(); ra < 1 {
		// If the query has failed, makes sure that the menu (and the user) exist
		_, err := m.GetOne(MID)
		return err
	}

	return nil
}

// Duplicate creates a copy of the menu. It returns the MID of the copy
func (m Menus) Duplicate(srcMID int) (int, error) {
	var dstMID int

	// Copies the menu
	err := db.QueryRow(`INSERT INTO menus (uid, name) SELECT uid, name FROM menus WHERE uid=$1 AND mid=$2 returning mid;`, m.uid, srcMID).Scan(&dstMID)
	if err != nil {
		return 0, handleNoRowsError(err, m.uid, ERR_MENU_NOT_FOUND)
	}

	// Copies the days
	_, err = db.Exec(`INSERT INTO days (mid, position, name, meals) SELECT $2, position, name, meals FROM days WHERE mid=$1;`, srcMID, dstMID)
	if err != nil {
		return dstMID, ERR_UNKNOWN
	}

	return dstMID, nil
}

// editDay is used by MoveDay, SetDayMeals and SetDayName
func (m Menus) editDay(MID int, day int, fetch bool, editMeals bool, meals []string, editName bool, name string) error {
	if fetch {
		// Gets the menu
		_, err := m.GetDay(MID, day)
		if err != nil {
			return err
		}
	}

	if editMeals {
		// Saves the new meals
		_, err := db.Exec(`UPDATE days SET meals=$3 WHERE mid=$1 AND position=$2`, MID, day, pq.Array(meals))
		if err != nil {
			return ERR_UNKNOWN
		}
	}

	if editName {
		// Saves the new name
		_, err := db.Exec(`UPDATE days SET name=$3 WHERE mid=$1 AND position=$2`, MID, day, name)
		if err != nil {
			return ERR_UNKNOWN
		}
	}

	return nil
}

// GetAll returns a list of menus (days not included) ordered by creation date
func (m Menus) GetAll() ([]Menu, error) {
	var menus []Menu

	// Queries the entries
	var rows *sql.Rows
	rows, err := db.Query(`SELECT mid, name FROM menus WHERE uid=$1 ORDER BY mid;`, m.uid)
	if err != nil {
		return menus, ERR_UNKNOWN
	}

	// Appends them to the list
	defer rows.Close()
	for rows.Next() {
		var m Menu
		rows.Scan(&m.MID, &m.Name)
		menus = append(menus, m)
	}

	// If no menus have been found, makes sure the user exists
	if len(menus) == 0 {
		_, err := GetUser("UID", m.uid)
		return menus, err
	}

	return menus, nil
}

// GetDay returns a day of a menu
func (m Menus) GetDay(MID int, dpos int) (Day, error) {
	var day Day

	// Scans the menu
	var a int
	err := db.QueryRow(`SELECT 1 FROM menus WHERE uid=$1 AND mid=$2;`, m.uid, MID).Scan(&a)
	if err != nil {
		return day, handleNoRowsError(err, m.uid, ERR_MENU_NOT_FOUND)
	}

	// Queries the day
	err = db.QueryRow(`SELECT mid, name, position, meals FROM days WHERE mid=$1 AND position=$2;`, MID, dpos).
		Scan(&day.MID, &day.Name, &day.Position, pq.Array(&day.Meals))
	if err != nil {
		return day, handleNoRowsError(err, m.uid, ERR_DAY_NOT_FOUND)
	}

	return day, nil
}

// GetOne returns a specific menu, with days and meals
func (m Menus) GetOne(MID int) (Menu, error) {
	var menu Menu

	// Scans the menu
	err := db.QueryRow(`SELECT mid, name FROM menus WHERE uid=$1 AND mid=$2;`, m.uid, MID).Scan(&menu.MID, &menu.Name)
	if err != nil {
		return menu, handleNoRowsError(err, m.uid, ERR_MENU_NOT_FOUND)
	}

	// Queries the days
	var rows *sql.Rows
	rows, err = db.Query(`SELECT name, position, meals FROM days WHERE mid=$1 ORDER BY position;`, menu.MID)
	if err != nil {
		return menu, ERR_UNKNOWN
	}
	defer rows.Close()

	// Appends the days and the meals to the menu
	for rows.Next() {
		day := Day{MID: MID}
		err := rows.Scan(&day.Name, &day.Position, pq.Array(&day.Meals))
		if err != nil {
			return menu, ERR_UNKNOWN
		}
		menu.Days = append(menu.Days, day)
	}

	return menu, nil
}

// MoveDay moves a day in the list (+1 switches it down, -1 switches it up, etc.)
func (m Menus) MoveDay(MID int, day int, delta int) error {
	// Gets the days
	dayA, err := m.GetDay(MID, day)
	if err != nil {
		return err
	}
	dayB, err := m.GetDay(MID, day+delta)
	if err != nil {
		return ERR_DAY_NOT_MOVED
	}

	// Switches the contents
	if err = m.editDay(MID, day, false, true, dayB.Meals, true, dayB.Name); err == nil {
		if err = m.editDay(MID, day+delta, false, true, dayA.Meals, true, dayA.Name); err == nil {
			return nil
		}
	}

	m.editDay(MID, day, false, true, dayA.Meals, true, dayA.Name)
	m.editDay(MID, day+delta, false, true, dayB.Meals, true, dayB.Name)

	return err
}

// New creates a new menu and return its MID
func (m Menus) New(name string, daysNames []string, mealsN int) (int, error) {
	var MID int

	// Ensures the user exists
	if _, err := GetUser("UID", m.uid); err != nil {
		return MID, err
	}

	// Ensures the number of meals is valid
	if mealsN < 0 {
		return MID, ERR_MEALS_NEGATIVE
	}

	// Prepares the statement for the days
	stmt, err := db.Prepare(`INSERT INTO days (mid, position, name, meals) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return MID, ERR_UNKNOWN
	}
	defer stmt.Close()

	// Adds the menu
	err = db.QueryRow(`INSERT INTO menus (uid, name) VALUES ($1, $2) RETURNING mid;`, m.uid, name).Scan(&MID)
	if err != nil {
		return MID, ERR_UNKNOWN
	}

	// Adds the days
	for dpos, dname := range daysNames {
		meals := make([]string, mealsN)
		_, err := stmt.Exec(MID, dpos, dname, pq.Array(meals))
		if err != nil {
			return MID, ERR_UNKNOWN
		}
	}

	return MID, nil
}

// RemoveDay removes a day from a menu
func (m Menus) RemoveDay(MID int, day int) error {
	// Gets the day
	_, err := m.GetDay(MID, day)
	if err != nil {
		return err
	}

	// Removes the desired day
	_, err = db.Exec(`DELETE FROM days WHERE mid=$1 AND position=$2;`, MID, day)
	if err != nil {
		return ERR_UNKNOWN
	}

	// Adjust the remaining positions
	_, err = db.Exec(`UPDATE days SET position=position-1 WHERE mid=$1 AND position>$2;`, MID, day)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// SetDayMeals is used to set a day's meals
func (m Menus) SetDayMeals(MID int, day int, meals []string) error {
	return m.editDay(MID, day, true, true, meals, false, "")
}

// SetDayName is used to set a day's name
func (m Menus) SetDayName(MID int, day int, name string) error {
	return m.editDay(MID, day, true, false, nil, true, name)
}

// SetName is used to set the menu's name
func (m Menus) SetName(MID int, name string) error {
	// Gets the menu
	_, err := m.GetOne(MID)
	if err != nil {
		return err
	}

	// Saves the new name
	_, err = db.Exec(`UPDATE menus SET name=$2 WHERE mid=$1;`, MID, name)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}
