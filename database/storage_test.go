package database

import (
	"reflect"
	"testing"
	"time"
)

func TestGetSections(t *testing.T) {
	user, _ := GetTestingUser(t)
	s1, _ := user.NewSection("s1")
	s2, _ := user.NewSection("s2")

	otherUser, _ := GetTestingUser(t)

	otherOtherUser, _ := GetTestingUser(t)
	otherOtherUser.NewSection("s")

	type data struct {
		User *User

		ExpectedErr      error
		ExpectedSections []Section
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
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{User: otherUser},
			},
			{
				"(filled)",
				data{User: user, ExpectedSections: []Section{s1, s2}},
			},
		},
	}.Run(t)
}

func TestGetSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s")

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *User
		SID  int

		ExpectedErr     error
		ExpectedSection Section
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.GetSection(d.SID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(got, d.ExpectedSection) {
				t.Errorf("%s: expected section <%v>, got <%v>", msg, d.ExpectedSection, got)
			}
		},

		Cases: []TestCase[data]{
			{
				"got data of unknown section",
				data{User: user, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user retrieved section",
				data{User: otherUser, SID: section.SID, ExpectedErr: ERR_SECTION_NOT_FOUND},
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
		User *User
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
				section, _ := d.User.GetSection(got.SID)
				if !reflect.DeepEqual(section, got) {
					t.Errorf("%v, returned bad section", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"unknown user created section",
				data{User: unknownUser, Name: "s", ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user, Name: "s1"},
			},
			{
				"created duplicated section",
				data{User: user, Name: "s1", ExpectedErr: ERR_SECTION_DUPLICATED},
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
		User    *User
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
				section, _ := d.User.GetSection(d.SID)
				expected := Section{SID: d.SID, Name: d.NewName}
				if !reflect.DeepEqual(section, expected) {
					t.Errorf("%v, changes not saved", msg)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"other user edited section",
				data{User: otherUser, SID: section.SID, NewName: "s3", ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"edited unknown section",
				data{User: user, NewName: "s3", ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"duplicated section",
				data{User: user, SID: section.SID, NewName: "s2", ExpectedErr: ERR_SECTION_DUPLICATED},
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
	user.AddArticles(section.SID, StringArticle{"article", "", ""})

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User *User
		SID  int

		ExpectedErr error
		ShouldExist bool
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.DeleteSection(d.SID); err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			section, _ = user.GetSection(d.SID)
			if !d.ShouldExist && section.SID != 0 {
				t.Errorf("%s, section wasn't deleted", msg)
			} else if d.ShouldExist && section.SID == 0 {
				t.Errorf("%s, section was deleted anyway", msg)
			}
		},

		Cases: []TestCase[data]{
			{
				"other user deleted section",
				data{User: otherUser, SID: section.SID, ExpectedErr: ERR_SECTION_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown section",
				data{User: user, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"",
				data{User: user, SID: section.SID},
			},
		},
	}.Run(t)
}

func TestAddArticles(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("section")
	SID := section.SID

	otherUser, _ := GetTestingUser(t)

	name := "article"
	sQty := "10"
	sExp := "2024-10-05"

	qty := 10
	exp := time.Date(2024, time.October, 5, 0, 0, 0, 0, time.FixedZone("", 0))

	inList := []StringArticle{
		{Name: "Full", Quantity: sQty, Expiration: sExp},
		{Name: "NoQty", Quantity: "", Expiration: sExp},
		{Name: "NoExp", Quantity: sQty, Expiration: ""},
		{Name: "Empty", Quantity: "", Expiration: ""},
	}

	outList := []Article{
		{AID: 2, Name: "Full", Quantity: &qty, Expiration: &exp},
		{AID: 3, Name: "NoQty", Quantity: nil, Expiration: &exp},
		{AID: 4, Name: "NoExp", Quantity: &qty, Expiration: nil},
		{AID: 5, Name: "Empty", Quantity: nil, Expiration: nil},
	}

	type data struct {
		User     *User
		SID      int
		Articles []StringArticle

		ExpectedErr      error
		ExpectedArticles []Article
	}

	TestSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.AddArticles(d.SID, d.Articles...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				section, _ := d.User.GetArticles(d.SID, "")
				if !reflect.DeepEqual(section.Articles, d.ExpectedArticles) {
					t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedArticles, section.Articles)
				}
			}
		},

		Cases: []TestCase[data]{
			{
				"added articles to unknown section",
				data{User: user, Articles: []StringArticle{{Name: name, Quantity: sQty, Expiration: sExp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user added articles to section",
				data{User: otherUser, SID: SID, Articles: []StringArticle{{Name: name, Quantity: sQty, Expiration: sExp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added article with invalid quantity",
				data{User: user, SID: SID, Articles: []StringArticle{{Name: name, Quantity: "a lot", Expiration: sExp}}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"added article with invalid sExpiration",
				data{User: user, SID: SID, Articles: []StringArticle{{Name: name, Quantity: sQty, Expiration: "dunno"}}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"",
				data{User: user, SID: SID, Articles: inList, ExpectedArticles: outList},
			},
		},
	}.Run(t)
}
