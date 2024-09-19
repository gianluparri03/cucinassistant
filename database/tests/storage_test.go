package tests

import (
	"reflect"
	"testing"

	"cucinassistant/database"
)

func TestGetSections(t *testing.T) {
	user, _ := GetTestingUser(t)
	s1, _ := user.NewSection("s1")
	s2, _ := user.NewSection("s2")
	// TODO add an article to s2

	otherUser, _ := GetTestingUser(t)

	otherOtherUser, _ := GetTestingUser(t)
	otherOtherUser.NewSection("s")

	type data struct {
		User *database.User

		ExpectedErr      error
		ExpectedSections []*database.Section
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			sections, err := d.User.GetSections()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(sections, d.ExpectedSections) {
				t.Errorf("%s: expected sections <%v>, got <%v>", msg, d.ExpectedSections, sections)
			}
		},

		Cases: []TestCase[data]{
			{
				"got sections of unknown user",
				data{User: unknownUser, ExpectedErr: database.ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{User: otherUser, ExpectedSections: []*database.Section{}},
			},
			{
				"(filled)",
				data{User: user, ExpectedSections: []*database.Section{s1, s2}},
			},
		},
	}.Run(t)
}

func TestGetSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s")
	// TODO add articles

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		SID  int

		ExpectedErr     error
		ExpectedSection *database.Section
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			var expectedSections map[bool]*database.Section
			if withArticles := d.ExpectedSection; withArticles != nil {
				withoutArticles := &database.Section{SID: withArticles.SID, Name: withArticles.Name}
				expectedSections = map[bool]*database.Section{true: withArticles, false: withoutArticles}
			} else {
				expectedSections = map[bool]*database.Section{true: nil, false: nil}
			}

			for _, withArticles := range []bool{true, false} {
				expectedSection := expectedSections[withArticles]

				got, err := d.User.GetSection(d.SID, withArticles)
				if err != d.ExpectedErr {
					t.Errorf("%s (%v): expected err <%v>, got <%v>", msg, withArticles, d.ExpectedErr, err)
				} else if !reflect.DeepEqual(got, expectedSection) {
					t.Errorf("%s (%v): expected section <%v>, got <%v>", msg, withArticles, expectedSection, got)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"got data of unknown section",
				data{User: user, ExpectedErr: database.ERR_SECTION_NOT_FOUND},
			},
			{
				"other user retrieved section",
				data{User: otherUser, SID: section.SID, ExpectedErr: database.ERR_SECTION_NOT_FOUND},
			},
			{
				"",
				data{User: user, SID: section.SID, ExpectedSection: section},
			},
		},
	}.Run(t)
}

func TestNewSection(t *testing.T) {
	user, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		Name string

		ExpectedErr error
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.NewSection(d.Name)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if err == nil {
				section, _ := d.User.GetSection(got.SID, false)
				if !reflect.DeepEqual(section, got) {
					t.Errorf("%v, returned bad section", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"unknown user created section",
				data{User: unknownUser, Name: "s", ExpectedErr: database.ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, Name: "s1"},
			},
			{
				"created duplicated section",
				data{User: user, Name: "s1", ExpectedErr: database.ERR_SECTION_DUPLICATED},
			},
			{
				"",
				data{User: user, Name: "s2"},
			},
		},
	}.Run(t)
}

func TestEditSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s1")
	user.NewSection("s2")

	otherUser, _ := GetTestingUser(t)
	otherUser.NewSection("s3")

	type data struct {
		User    *database.User
		SID     int
		NewName string

		ExpectedErr error
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.EditSection(d.SID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				section, _ := d.User.GetSection(d.SID, false)
				expected := &database.Section{SID: d.SID, Name: d.NewName}
				if !reflect.DeepEqual(section, expected) {
					t.Errorf("%v, changes not saved", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user edited section",
				data{User: otherUser, SID: section.SID, NewName: "s3", ExpectedErr: database.ERR_SECTION_NOT_FOUND},
			},
			{
				"edited unknown section",
				data{User: user, NewName: "s3", ExpectedErr: database.ERR_SECTION_NOT_FOUND},
			},
			{
				"duplicated section",
				data{User: user, SID: section.SID, NewName: "s2", ExpectedErr: database.ERR_SECTION_DUPLICATED},
			},
			{
				"(same)",
				data{User: user, SID: section.SID, NewName: "s1"},
			},
			{
				"(different)",
				data{User: user, SID: section.SID, NewName: "s3"},
			},
		},
	}.Run(t)
}

func TestDeleteSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s")
	// TODO add articles

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *database.User
		SID  int

		ExpectedErr error
		ShouldExist bool
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.DeleteSection(d.SID); err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			section, _ = user.GetSection(d.SID, true)
			if !d.ShouldExist && section != nil {
				t.Errorf("%s, section wasn't deleted", msg)
			} else if d.ShouldExist && section == nil {
				t.Errorf("%s, section was deleted anyway", msg)
			}
		},

		Cases: []TestCase[data]{
			{
				"other user deleted section",
				data{User: otherUser, SID: section.SID, ExpectedErr: database.ERR_SECTION_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown section",
				data{User: user, ExpectedErr: database.ERR_SECTION_NOT_FOUND},
			},
			{
				"",
				data{User: user, SID: section.SID},
			},
		},
	}.Run(t)
}
