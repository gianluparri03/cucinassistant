package database

import (
	"reflect"
	"testing"
)

var testingArticlesN int = 0

func (sa StringArticle) getExpectedArticle() Article {
	a, _ := sa.Parse()
	a.fixExpiration()

	testingArticlesN++
	a.AID = testingArticlesN
	return a
}

func TestStorageAddArticles(t *testing.T) {
	user, _ := getTestingUser(t)
	section, _ := user.Storage().NewSection("section")
	sid := section.SID
	otherSection, _ := user.Storage().NewSection("otherSection")

	name := "article"
	qty := "10.07"
	exp := "2024-10-05"

	inList := []StringArticle{
		{Name: "NoQty", Quantity: "", Expiration: exp, Section: sid},
		{Name: "Full", Quantity: qty, Expiration: exp, Section: sid},
		{Name: "Empty", Quantity: "", Expiration: "", Section: sid},
		{Name: "NoExp", Quantity: qty, Expiration: "", Section: sid},
		{Name: "otherSection", Quantity: qty, Expiration: exp, Section: otherSection.SID},
	}

	outSimple := make(map[int][]Article)
	outDoubled := make(map[int][]Article)
	for _, sa := range inList {
		simple := sa.getExpectedArticle()
		outSimple[sa.Section] = append(outSimple[sa.Section], simple)

		doubled := simple
		if simple.Quantity == nil {
			doubled.Quantity = nil
		} else {
			qty := (*simple.Quantity) * 2
			doubled.Quantity = &qty
		}
		outDoubled[sa.Section] = append(outDoubled[sa.Section], doubled)
	}

	otherUser, _ := getTestingUser(t)
	notMySection, _ := otherUser.Storage().NewSection("section")

	testingArticlesN += 5

	type data struct {
		User     User
		Articles []StringArticle

		ExpectedErr      error
		ExpectedArticles map[int][]Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.Storage().AddArticles(d.Articles...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				for sid, articles := range d.ExpectedArticles {
					section, _ := d.User.Storage().GetArticles(sid, "")
					if !reflect.DeepEqual(section.Articles, articles) {
						t.Errorf("%s: expected list <%v>, got <%v>", msg, articles, section.Articles)
					}
				}
			}
		},

		Cases: []testCase[data]{
			{
				"added articles to unknown section",
				data{User: user, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added articles to another user's section",
				data{User: user, Articles: []StringArticle{{Name: name, Section: notMySection.SID}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user added articles to section",
				data{User: otherUser, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added article with invalid quantity",
				data{User: user, Articles: []StringArticle{{Name: name, Quantity: "a lot", Expiration: exp, Section: sid}}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"added article with invalid expiration",
				data{User: user, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: "dunno", Section: sid}}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"(original)",
				data{User: user, Articles: inList, ExpectedArticles: outSimple},
			},
			{
				"(doubled)",
				data{User: user, Articles: inList, ExpectedArticles: outDoubled},
			},
		},
	}.Run(t)
}

func TestStorageDeleteArticle(t *testing.T) {
	user, _ := getTestingUser(t)

	section, _ := user.Storage().NewSection("section")
	sid := section.SID

	sArticle := StringArticle{Name: "article", Expiration: "2024-12-18", Quantity: "5", Section: sid}
	sNextArticle := StringArticle{Name: "next-article", Section: sid}
	sSimilarArticle := StringArticle{Name: "article", Expiration: "2023-02-18", Section: sid}
	user.Storage().AddArticles(sArticle, sNextArticle, sSimilarArticle)

	article := sArticle.getExpectedArticle()
	nextArticle := sNextArticle.getExpectedArticle()
	similarArticle := sSimilarArticle.getExpectedArticle()

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		AID  int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.Storage().DeleteArticle(d.AID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if d.AID != similarArticle.AID {
				if _, err := user.Storage().GetArticle(similarArticle.AID); err != nil {
					t.Errorf("%s, deleted similar article", msg)
				}
			}

			article, _ := user.Storage().GetArticle(d.AID)
			if !d.ShouldExist && article.AID != 0 {
				t.Errorf("%s, article wasn't deleted", msg)
			} else if d.ShouldExist && article.AID == 0 {
				t.Errorf("%s, article was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted article",
				data{User: otherUser, AID: article.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown article",
				data{User: user, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: false},
			},
			{
				"(has next)",
				data{User: user, AID: article.AID, ShouldExist: false},
			},
			{
				"(has not next)",
				data{User: user, AID: nextArticle.AID, ShouldExist: false},
			},
		},
	}.Run(t)
}

func TestStorageDeleteSection(t *testing.T) {
	user, _ := getTestingUser(t)
	section, _ := user.Storage().NewSection("s")
	user.Storage().AddArticles(StringArticle{Name: "article", Section: section.SID})
	testingArticlesN++

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		SID  int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.Storage().DeleteSection(d.SID); err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			section, _ = user.Storage().GetSection(d.SID)
			if !d.ShouldExist && section.SID != 0 {
				t.Errorf("%s, section wasn't deleted", msg)
			} else if d.ShouldExist && section.SID == 0 {
				t.Errorf("%s, section was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
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

func TestStorageEditArticle(t *testing.T) {
	user, _ := getTestingUser(t)

	section, _ := user.Storage().NewSection("section")
	sid := section.SID

	sArticle := StringArticle{Name: "article", Expiration: "2024-12-18", Quantity: "5", Section: sid}
	sOtherArticle := StringArticle{Name: "otherArticle", Expiration: "2024-12-31", Quantity: "9", Section: sid}
	user.Storage().AddArticles(sArticle, sOtherArticle)

	article := sArticle.getExpectedArticle()
	testingArticlesN++

	newQty := StringArticle{Name: "article", Expiration: "2024-12-18"}
	newAll := StringArticle{Name: "Article", Expiration: "2024-12-25", Quantity: "6"}
	newAllWithChange := StringArticle{Name: "article", Expiration: "2025-01-02", Quantity: "9"}

	otherUser, _ := getTestingUser(t)

	type data struct {
		User    User
		AID     int
		NewData StringArticle

		ExpectedErr     error
		ExpectedChanged bool
		CheckEdits      bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.Storage().EditArticle(d.AID, d.NewData); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if d.CheckEdits {
				gotData, _ := user.Storage().GetArticle(d.AID)
				expectedData, _ := d.NewData.Parse()

				if !reflect.DeepEqual(gotData.Name, expectedData.Name) {
					t.Errorf("%s, name not saved", msg)
				} else if !reflect.DeepEqual(gotData.Expiration, expectedData.Expiration) {
					t.Errorf("%s, expiration not saved", msg)
				} else if !reflect.DeepEqual(gotData.Quantity, expectedData.Quantity) {
					t.Errorf("%s, quantity not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited article",
				data{User: otherUser, AID: article.AID, NewData: newAll, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"edited unknown article",
				data{User: user, NewData: newAll, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"invalid quantity",
				data{User: user, AID: article.AID, NewData: StringArticle{Quantity: "a few"}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"invalid expiration",
				data{User: user, AID: article.AID, NewData: StringArticle{Expiration: "next year"}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"duplicated article",
				data{User: user, AID: article.AID, NewData: sOtherArticle, ExpectedErr: ERR_ARTICLE_DUPLICATED},
			},
			{
				"(same)",
				data{User: user, AID: article.AID, NewData: sArticle, ExpectedChanged: false, CheckEdits: true},
			},
			{
				"(only quantity)",
				data{User: user, AID: article.AID, NewData: newQty, ExpectedChanged: false, CheckEdits: true},
			},
			{
				"(all, changed=false)",
				data{User: user, AID: article.AID, NewData: newAll, ExpectedChanged: false, CheckEdits: true},
			},
			{
				"(all, changed=true)",
				data{User: user, AID: article.AID, NewData: newAllWithChange, ExpectedChanged: true, CheckEdits: true},
			},
		},
	}.Run(t)
}

func TestStorageEditSection(t *testing.T) {
	user, _ := getTestingUser(t)
	section, _ := user.Storage().NewSection("s1")
	user.Storage().NewSection("s2")

	otherUser, _ := getTestingUser(t)
	otherUser.Storage().NewSection("s3")

	type data struct {
		User    User
		SID     int
		NewName string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.User.Storage().EditSection(d.SID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				section, _ := d.User.Storage().GetSection(d.SID)
				expected := Section{SID: d.SID, Name: d.NewName}
				if !reflect.DeepEqual(section, expected) {
					t.Errorf("%v, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
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

func TestStorageGetArticle(t *testing.T) {
	user, _ := getTestingUser(t)
	otherUser, _ := getTestingUser(t)

	section, _ := user.Storage().NewSection("section")
	sid := section.SID
	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2", Section: sid}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11", Section: sid}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15", Section: sid}
	user.Storage().AddArticles(s1, s2, s3)

	otherSection, _ := user.Storage().NewSection("otherSection")
	s4 := StringArticle{Name: "middle", Expiration: "2024-10-31", Section: otherSection.SID}
	user.Storage().AddArticles(s4)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()
	a4 := s4.getExpectedArticle()

	// Makes sure that it works even after edits
	s1.Quantity = "7"
	qty := float32(7)
	a1.Quantity = &qty
	user.Storage().EditArticle(a1.AID, s1)

	type data struct {
		User User
		AID  int

		ExpectedErr     error
		ExpectedArticle Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			article, err := d.User.Storage().GetArticle(d.AID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(article, d.ExpectedArticle) {
				t.Errorf("%s: expected article <%v>, got <%v>", msg, d.ExpectedArticle, article)
			}
		},

		Cases: []testCase[data]{
			{
				"got unknown article",
				data{User: user, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"other user retrieved article",
				data{User: otherUser, AID: a1.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"(a1)",
				data{User: user, AID: a1.AID, ExpectedArticle: a1},
			},
			{
				"(a2)",
				data{User: user, AID: a2.AID, ExpectedArticle: a2},
			},
			{
				"(a3)",
				data{User: user, AID: a3.AID, ExpectedArticle: a3},
			},
			{
				"(a4)",
				data{User: user, AID: a4.AID, ExpectedArticle: a4},
			},
		},
	}.Run(t)
}

func TestStorageGetArticlesGeneral(t *testing.T) {
	user, _ := getTestingUser(t)

	section1, _ := user.Storage().NewSection("section1")
	section2, _ := user.Storage().NewSection("section2")

	s1 := StringArticle{Name: "First", Expiration: "2025-01-08", Quantity: "2", Section: section1.SID}
	s2 := StringArticle{Name: "second", Expiration: "2025-01-06", Quantity: "15.31", Section: section2.SID}
	user.Storage().AddArticles(s1, s2)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()

	otherUser, _ := getTestingUser(t)
	otherUserSection, _ := otherUser.Storage().NewSection("section")
	otherUser.Storage().AddArticles(StringArticle{Name: "other", Section: otherUserSection.SID})
	testingArticlesN++

	type data struct {
		User   User
		Filter string

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			articles, err := d.User.Storage().GetArticles(0, d.Filter)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(articles.Articles, d.ExpectedArticles) {
				t.Errorf("%s: expected list <%#v>, got <%#v>", msg, d.ExpectedArticles, articles.Articles)
			}
		},

		Cases: []testCase[data]{
			{
				"got articles of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(unfiltered)",
				data{User: user, ExpectedArticles: []Article{a2, a1}},
			},
			{
				"(filtered)",
				data{User: user, Filter: "F", ExpectedArticles: []Article{a1}},
			},
		},
	}.Run(t)
}

func TestStorageGetArticlesSection(t *testing.T) {
	user, _ := getTestingUser(t)

	section, _ := user.Storage().NewSection("section")
	sid := section.SID

	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2", Section: sid}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11", Section: sid}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15.31", Section: sid}
	user.Storage().AddArticles(s1, s2, s3)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()

	otherSection, _ := user.Storage().NewSection("otherSection")
	s4 := StringArticle{Name: "fourth", Section: otherSection.SID}
	user.Storage().AddArticles(s4)
	testingArticlesN++

	otherUser, _ := getTestingUser(t)
	emptySection, _ := otherUser.Storage().NewSection("emptySection")

	// Makes sure that it works even after edits
	s1.Quantity = "7"
	qty := float32(7)
	a1.Quantity = &qty
	user.Storage().EditArticle(a1.AID, s1)
	expected := []Article{a2, a1, a3}

	type data struct {
		User   User
		SID    int
		Filter string

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			section, err := d.User.Storage().GetArticles(d.SID, d.Filter)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(section.Articles, d.ExpectedArticles) {
				t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedArticles, section.Articles)
			}
		},

		Cases: []testCase[data]{
			{
				"got articles of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"got articles of unknown section",
				data{User: user, SID: -1, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user retrieved section",
				data{User: otherUser, SID: section.SID, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"(filled)",
				data{User: user, SID: section.SID, ExpectedArticles: expected},
			},
			{
				"(filtered)",
				data{User: user, SID: section.SID, Filter: "th", ExpectedArticles: []Article{a3}},
			},
			{
				"(empty)",
				data{User: otherUser, SID: emptySection.SID},
			},
		},
	}.Run(t)
}

func TestStorageGetNeighbours(t *testing.T) {
	user, _ := getTestingUser(t)

	section, _ := user.Storage().NewSection("section")
	sid := section.SID
	s1 := StringArticle{Expiration: "2024-11-08", Section: sid}
	s2 := StringArticle{Expiration: "2024-10-11", Section: sid}
	s3 := StringArticle{Expiration: "", Section: sid}
	user.Storage().AddArticles(s1, s2, s3)

	otherSection, _ := user.Storage().NewSection("otherSection")
	s4 := StringArticle{Expiration: "2024-10-31", Section: otherSection.SID}
	user.Storage().AddArticles(s4)

	a1 := s1.getExpectedArticle().AID
	a2 := s2.getExpectedArticle().AID
	a3 := s3.getExpectedArticle().AID
	a4 := s4.getExpectedArticle().AID

	type data struct {
		SID int
		AID int

		ExpectedPrev int
		ExpectedNext int
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			prev, next := user.Storage().GetNeighbours(d.SID, d.AID)
			if !reflect.DeepEqual(prev, d.ExpectedPrev) {
				t.Errorf("%s: expected prev <%v>, got <%v>", msg, d.ExpectedPrev, prev)
			} else if !reflect.DeepEqual(next, d.ExpectedNext) {
				t.Errorf("%s: expected next <%v>, got <%v>", msg, d.ExpectedNext, next)
			}
		},

		Cases: []testCase[data]{
			{
				"(a1, section)",
				data{SID: sid, AID: a1, ExpectedPrev: a2, ExpectedNext: a3},
			},
			{
				"(a2, section)",
				data{SID: sid, AID: a2, ExpectedNext: a1},
			},
			{
				"(a3, section)",
				data{SID: sid, AID: a3, ExpectedPrev: a1},
			},
			{
				"(a4, section)",
				data{SID: sid, AID: a4},
			},
			{
				"(a1, general)",
				data{AID: a1, ExpectedPrev: a4, ExpectedNext: a3},
			},
			{
				"(a2, general)",
				data{AID: a2, ExpectedNext: a4},
			},
			{
				"(a3, general)",
				data{AID: a3, ExpectedPrev: a1},
			},
			{
				"(a4, general)",
				data{AID: a4, ExpectedPrev: a2, ExpectedNext: a1},
			},
		},
	}.Run(t)
}

func TestStorageGetSection(t *testing.T) {
	user, _ := getTestingUser(t)
	section, _ := user.Storage().NewSection("s")

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		SID  int

		ExpectedErr     error
		ExpectedSection Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.Storage().GetSection(d.SID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(got, d.ExpectedSection) {
				t.Errorf("%s: expected section <%v>, got <%v>", msg, d.ExpectedSection, got)
			}
		},

		Cases: []testCase[data]{
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

func TestStorageGetSections(t *testing.T) {
	user, _ := getTestingUser(t)
	s1, _ := user.Storage().NewSection("s1")
	s2, _ := user.Storage().NewSection("s2")

	otherUser, _ := getTestingUser(t)

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.Storage().NewSection("s")

	type data struct {
		User User

		ExpectedErr      error
		ExpectedSections []Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			sections, err := d.User.Storage().GetSections()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(sections, d.ExpectedSections) {
				t.Errorf("%s: expected sections <%v>, got <%v>", msg, d.ExpectedSections, sections)
			}
		},

		Cases: []testCase[data]{
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

func TestStorageNewSection(t *testing.T) {
	user, _ := getTestingUser(t)

	type data struct {
		User User
		Name string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.Storage().NewSection(d.Name)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if err == nil {
				section, _ := d.User.Storage().GetSection(got.SID)
				if !reflect.DeepEqual(section, got) {
					t.Errorf("%v, returned bad section", msg)
				}
			}
		},

		Cases: []testCase[data]{
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
