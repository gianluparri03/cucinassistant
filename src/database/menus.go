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

// New creates a new menu
func (m Menus) New(name string, daysNames []string, mealsN int) (Menu, error) {
	// Ensures the user exists
	if _, err := GetUser("UID", m.uid); err != nil {
		return Menu{}, err
	}

	// Ensures the number of meals is valid
	if mealsN < 0 {
		return Menu{}, ERR_MEALS_NEGATIVE
	}

	// Prepares the statement for the days
	stmt, err := db.Prepare(`INSERT INTO days (mid, position, name, meals) VALUES ($1, $2, $3, $4);`)
	if err != nil {
		return Menu{}, ERR_UNKNOWN
	}
	defer stmt.Close()

	// Adds the menu
	menu := Menu{Name: name}
	err = db.QueryRow(`INSERT INTO menus (uid, name) VALUES ($1, $2) RETURNING mid;`, m.uid, name).Scan(&menu.MID)
	if err != nil {
		return Menu{}, ERR_UNKNOWN
	}

	// Adds the days
	for p, n := range daysNames {
		meals := make([]string, mealsN)
		_, err := stmt.Exec(menu.MID, p, n, pq.Array(meals))
		if err != nil {
			return menu, ERR_UNKNOWN
		}

		day := Day{MID: menu.MID, Name: n, Position: p, Meals: meals}
		menu.Days = append(menu.Days, day)
	}

	return menu, nil
}

// SetDayMeals is used to set a day's meals
func (m Menus) SetDayMeals(MID int, day int, meals []string) error {
	// Gets the menu
	_, err := m.GetDay(MID, day)
	if err != nil {
		return err
	}

	// Saves the new meals
	_, err = db.Exec(`UPDATE days SET meals=$3 WHERE mid=$1 AND position=$2`, MID, day, pq.Array(meals))
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
}

// SetDayName is used to set a day's name
func (m Menus) SetDayName(MID int, day int, name string) error {
	// Gets the menu
	_, err := m.GetDay(MID, day)
	if err != nil {
		return err
	}

	// Saves the new meals
	_, err = db.Exec(`UPDATE days SET name=$3 WHERE mid=$1 AND position=$2`, MID, day, name)
	if err != nil {
		return ERR_UNKNOWN
	}

	return nil
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
