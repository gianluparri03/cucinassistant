package database

import (
	"reflect"
	"strconv"
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
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("section")
	otherSID, _ := s.NewSection("otherSection")

	name := "article"
	qty := "10.07"
	exp := "2024-10-05"

	inList := []StringArticle{
		{Name: "NoQty", Quantity: "", Expiration: exp, Section: strconv.Itoa(SID)},
		{Name: "Full", Quantity: qty, Expiration: exp, Section: strconv.Itoa(SID)},
		{Name: "Empty", Quantity: "", Expiration: "", Section: strconv.Itoa(SID)},
		{Name: "NoExp", Quantity: qty, Expiration: "", Section: strconv.Itoa(SID)},
		{Name: "otherSection", Quantity: qty, Expiration: exp, Section: strconv.Itoa(otherSID)},
	}

	outSimple := make(map[string][]Article)
	outDoubled := make(map[string][]Article)
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

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()
	notMySID, _ := otherS.NewSection("section")

	testingArticlesN += 5

	type data struct {
		S        Storage
		Articles []StringArticle

		ExpectedErr      error
		ExpectedArticles map[string][]Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.S.AddArticles(d.Articles...)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				for SIDStr, articles := range d.ExpectedArticles {
					SID, _ := strconv.Atoi(SIDStr)
					section, _ := d.S.GetArticles(SID, "")
					if !reflect.DeepEqual(section.Articles, articles) {
						t.Errorf("%s: expected list <%v>, got <%v>", msg, articles, section.Articles)
					}
				}
			}
		},

		Cases: []testCase[data]{
			{
				"added articles to unknown section",
				data{S: s, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added articles to another user's section",
				data{S: s, Articles: []StringArticle{{Name: name, Section: strconv.Itoa(notMySID)}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user added articles to section",
				data{S: otherS, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added article with invalid quantity",
				data{S: s, Articles: []StringArticle{{Name: name, Quantity: "a lot", Expiration: exp, Section: strconv.Itoa(SID)}}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"added article with invalid expiration",
				data{S: s, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: "dunno", Section: strconv.Itoa(SID)}}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"(original)",
				data{S: s, Articles: inList, ExpectedArticles: outSimple},
			},
			{
				"(doubled)",
				data{S: s, Articles: inList, ExpectedArticles: outDoubled},
			},
		},
	}.Run(t)
}

func TestStorageDeleteArticle(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("section")

	sArticle := StringArticle{Name: "article", Expiration: "2024-12-18", Quantity: "5", Section: strconv.Itoa(SID)}
	sNextArticle := StringArticle{Name: "next-article", Section: strconv.Itoa(SID)}
	sSimilarArticle := StringArticle{Name: "article", Expiration: "2023-02-18", Section: strconv.Itoa(SID)}
	s.AddArticles(sArticle, sNextArticle, sSimilarArticle)

	article := sArticle.getExpectedArticle()
	nextArticle := sNextArticle.getExpectedArticle()
	similarArticle := sSimilarArticle.getExpectedArticle()

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()

	type data struct {
		S   Storage
		AID int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.S.DeleteArticle(d.AID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if d.AID != similarArticle.AID {
				if _, err := s.GetArticle(similarArticle.AID); err != nil {
					t.Errorf("%s, deleted similar article", msg)
				}
			}

			article, _ := s.GetArticle(d.AID)
			if !d.ShouldExist && article.AID != 0 {
				t.Errorf("%s, article wasn't deleted", msg)
			} else if d.ShouldExist && article.AID == 0 {
				t.Errorf("%s, article was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted article",
				data{S: otherS, AID: article.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown article",
				data{S: s, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: false},
			},
			{
				"(has next)",
				data{S: s, AID: article.AID, ShouldExist: false},
			},
			{
				"(has not next)",
				data{S: s, AID: nextArticle.AID, ShouldExist: false},
			},
		},
	}.Run(t)
}

func TestStorageDeleteSection(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("s")
	s.AddArticles(StringArticle{Name: "article", Section: strconv.Itoa(SID)})
	testingArticlesN++

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()

	type data struct {
		S   Storage
		SID int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.S.DeleteSection(d.SID); err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			section, _ := s.GetSection(d.SID)
			if !d.ShouldExist && section.SID != 0 {
				t.Errorf("%s, section wasn't deleted", msg)
			} else if d.ShouldExist && section.SID == 0 {
				t.Errorf("%s, section was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted section",
				data{S: otherS, SID: SID, ExpectedErr: ERR_SECTION_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown section",
				data{S: s, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"",
				data{S: s, SID: SID},
			},
		},
	}.Run(t)
}

func TestStorageEditArticle(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("section")
	sArticle := StringArticle{Name: "article", Expiration: "2024-12-18", Quantity: "5", Section: strconv.Itoa(SID)}
	s.AddArticles(sArticle)
	article := sArticle.getExpectedArticle()

	otherSID, _ := s.NewSection("otherSection")
	sOtherArticle := StringArticle{Name: "otherArticle", Expiration: "2024-12-31", Quantity: "9", Section: strconv.Itoa(otherSID)}
	s.AddArticles(sOtherArticle)
	sOtherArticle.getExpectedArticle()

	sDupArticle := sOtherArticle
	sDupArticle.Section = strconv.Itoa(SID)
	s.AddArticles(sDupArticle)
	dupArticle := sDupArticle.getExpectedArticle()

	notMyU, _ := getTestingUser(t)
	notMyS := notMyU.Storage()

	notMySID, _ := notMyS.NewSection("section")

	newQty := StringArticle{Name: "article", Expiration: "2024-12-18", Section: strconv.Itoa(SID)}
	newAll := StringArticle{Name: "Article", Expiration: "2024-12-25", Quantity: "6", Section: strconv.Itoa(SID)}
	newAllWithChange := StringArticle{Name: "article", Expiration: "2025-01-02", Quantity: "9", Section: strconv.Itoa(SID)}
	newSection := StringArticle{Section: strconv.Itoa(otherSID)}

	type data struct {
		S       Storage
		AID     int
		NewData StringArticle

		ExpectedErr error
		CheckEdits  bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.S.EditArticle(d.AID, d.NewData); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if d.CheckEdits {
				gotData, _ := s.GetArticle(d.AID)
				expectedData, _ := d.NewData.Parse()
				expectedData.fixExpiration()

				if !reflect.DeepEqual(gotData.Name, expectedData.Name) {
					t.Errorf("%s, name not saved", msg)
				} else if !reflect.DeepEqual(gotData.Expiration, expectedData.Expiration) {
					t.Errorf("%s, expiration not saved", msg)
				} else if !reflect.DeepEqual(gotData.Quantity, expectedData.Quantity) {
					t.Errorf("%s, quantity not saved", msg)
				} else if !reflect.DeepEqual(gotData.SID, expectedData.SID) {
					t.Errorf("%s, section not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited article",
				data{S: notMyS, AID: article.AID, NewData: newAll, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"edited unknown article",
				data{S: s, NewData: newAll, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"invalid section",
				data{S: s, AID: article.AID, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"invalid quantity",
				data{S: s, AID: article.AID, NewData: StringArticle{Section: strconv.Itoa(SID), Quantity: "a few"}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"invalid expiration",
				data{S: s, AID: article.AID, NewData: StringArticle{Section: strconv.Itoa(SID), Expiration: "next year"}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"duplicated article",
				data{S: s, AID: article.AID, NewData: sOtherArticle, ExpectedErr: ERR_ARTICLE_DUPLICATED},
			},
			{
				"duplicated article in target section",
				data{S: s, AID: dupArticle.AID, NewData: sOtherArticle, ExpectedErr: ERR_ARTICLE_DUPLICATED},
			},
			{
				"(same)",
				data{S: s, AID: article.AID, NewData: sArticle, CheckEdits: true},
			},
			{
				"(only quantity)",
				data{S: s, AID: article.AID, NewData: newQty, CheckEdits: true},
			},
			{
				"(expiration and quantity)",
				data{S: s, AID: article.AID, NewData: newAllWithChange, CheckEdits: true},
			},
			{
				"(section)",
				data{S: s, AID: article.AID, NewData: newSection, CheckEdits: true},
			},
			{
				"moved to other user's section",
				data{S: s, AID: article.AID, NewData: StringArticle{Section: strconv.Itoa(notMySID)}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
		},
	}.Run(t)
}

func TestStorageEditSection(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("s1")
	s.NewSection("s2")

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()
	otherS.NewSection("s3")

	type data struct {
		S       Storage
		SID     int
		NewName string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.S.EditSection(d.SID, d.NewName)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				section, _ := d.S.GetSection(d.SID)
				expected := Section{SID: d.SID, Name: d.NewName}
				if !reflect.DeepEqual(section, expected) {
					t.Errorf("%v, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited section",
				data{S: otherS, SID: SID, NewName: "s3", ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"edited unknown section",
				data{S: s, NewName: "s3", ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"duplicated section",
				data{S: s, SID: SID, NewName: "s2", ExpectedErr: ERR_SECTION_DUPLICATED},
			},
			{
				"(same)",
				data{S: s, SID: SID, NewName: "s1"},
			},
			{
				"(different)",
				data{S: s, SID: SID, NewName: "s3"},
			},
		},
	}.Run(t)
}

func TestStorageGetArticle(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()

	SID, _ := s.NewSection("section")
	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2", Section: strconv.Itoa(SID)}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11", Section: strconv.Itoa(SID)}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15", Section: strconv.Itoa(SID)}
	s.AddArticles(s1, s2, s3)

	otherSID, _ := s.NewSection("otherSection")
	s4 := StringArticle{Name: "middle", Expiration: "2024-10-31", Section: strconv.Itoa(otherSID)}
	s.AddArticles(s4)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()
	a4 := s4.getExpectedArticle()

	// Makes sure that it works even after edits
	s1.Quantity = "7"
	qty := float32(7)
	a1.Quantity = &qty
	s.EditArticle(a1.AID, s1)

	type data struct {
		S   Storage
		AID int

		ExpectedErr     error
		ExpectedArticle Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			article, err := d.S.GetArticle(d.AID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(article, d.ExpectedArticle) {
				t.Errorf("%s: expected article <%v>, got <%v>", msg, d.ExpectedArticle, article)
			}
		},

		Cases: []testCase[data]{
			{
				"got unknown article",
				data{S: s, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"other user retrieved article",
				data{S: otherS, AID: a1.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"(a1)",
				data{S: s, AID: a1.AID, ExpectedArticle: a1},
			},
			{
				"(a2)",
				data{S: s, AID: a2.AID, ExpectedArticle: a2},
			},
			{
				"(a3)",
				data{S: s, AID: a3.AID, ExpectedArticle: a3},
			},
			{
				"(a4)",
				data{S: s, AID: a4.AID, ExpectedArticle: a4},
			},
		},
	}.Run(t)
}

func TestStorageGetArticlesGeneral(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID1, _ := s.NewSection("section1")
	SID2, _ := s.NewSection("section2")

	s1 := StringArticle{Name: "First", Expiration: "2025-01-08", Quantity: "2", Section: strconv.Itoa(SID1)}
	s2 := StringArticle{Name: "second", Expiration: "2025-01-06", Quantity: "15.31", Section: strconv.Itoa(SID2)}
	s.AddArticles(s1, s2)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()
	otherSID, _ := otherS.NewSection("section")
	otherS.AddArticles(StringArticle{Name: "other", Section: strconv.Itoa(otherSID)})
	testingArticlesN++

	type data struct {
		S      Storage
		Filter string

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			articles, err := d.S.GetArticles(0, d.Filter)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(articles.Articles, d.ExpectedArticles) {
				t.Errorf("%s: expected list <%#v>, got <%#v>", msg, d.ExpectedArticles, articles.Articles)
			}
		},

		Cases: []testCase[data]{
			{
				"got articles of unknown user",
				data{S: unknownUser.Storage(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(unfiltered)",
				data{S: s, ExpectedArticles: []Article{a2, a1}},
			},
			{
				"(filtered)",
				data{S: s, Filter: "F", ExpectedArticles: []Article{a1}},
			},
		},
	}.Run(t)
}

func TestStorageGetArticlesSection(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("section")

	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2", Section: strconv.Itoa(SID)}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11", Section: strconv.Itoa(SID)}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15.31", Section: strconv.Itoa(SID)}
	s.AddArticles(s1, s2, s3)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()

	otherSID, _ := s.NewSection("otherSection")
	s4 := StringArticle{Name: "fourth", Section: strconv.Itoa(otherSID)}
	s.AddArticles(s4)
	testingArticlesN++

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()
	emptySID, _ := otherS.NewSection("emptySection")

	// Makes sure that it works even after edits
	s1.Quantity = "7"
	qty := float32(7)
	a1.Quantity = &qty
	s.EditArticle(a1.AID, s1)
	expected := []Article{a2, a1, a3}

	type data struct {
		S      Storage
		SID    int
		Filter string

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			section, err := d.S.GetArticles(d.SID, d.Filter)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(section.Articles, d.ExpectedArticles) {
				t.Errorf("%s: expected list <%v>, got <%v>", msg, d.ExpectedArticles, section.Articles)
			}
		},

		Cases: []testCase[data]{
			{
				"got articles of unknown user",
				data{S: unknownUser.Storage(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"got articles of unknown section",
				data{S: s, SID: -1, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user retrieved section",
				data{S: otherS, SID: SID, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"(filled)",
				data{S: s, SID: SID, ExpectedArticles: expected},
			},
			{
				"(filtered)",
				data{S: s, SID: SID, Filter: "th", ExpectedArticles: []Article{a3}},
			},
			{
				"(empty)",
				data{S: otherS, SID: emptySID},
			},
		},
	}.Run(t)
}

func TestStorageGetNeighbours(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("section")
	s1 := StringArticle{Expiration: "2024-11-08", Section: strconv.Itoa(SID)}
	s2 := StringArticle{Expiration: "2024-10-11", Section: strconv.Itoa(SID)}
	s3 := StringArticle{Expiration: "", Section: strconv.Itoa(SID)}
	s.AddArticles(s1, s2, s3)

	otherSID, _ := s.NewSection("otherSection")
	s4 := StringArticle{Expiration: "2024-10-31", Section: strconv.Itoa(otherSID)}
	s.AddArticles(s4)

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
			prev, next := s.GetNeighbours(d.SID, d.AID)
			if !reflect.DeepEqual(prev, d.ExpectedPrev) {
				t.Errorf("%s: expected prev <%v>, got <%v>", msg, d.ExpectedPrev, prev)
			} else if !reflect.DeepEqual(next, d.ExpectedNext) {
				t.Errorf("%s: expected next <%v>, got <%v>", msg, d.ExpectedNext, next)
			}
		},

		Cases: []testCase[data]{
			{
				"(a1, section)",
				data{SID: SID, AID: a1, ExpectedPrev: a2, ExpectedNext: a3},
			},
			{
				"(a2, section)",
				data{SID: SID, AID: a2, ExpectedNext: a1},
			},
			{
				"(a3, section)",
				data{SID: SID, AID: a3, ExpectedPrev: a1},
			},
			{
				"(a4, section)",
				data{SID: SID, AID: a4},
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
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID, _ := s.NewSection("s")
	section, _ := s.GetSection(SID)

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()

	type data struct {
		S   Storage
		SID int

		ExpectedErr     error
		ExpectedSection Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.S.GetSection(d.SID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(got, d.ExpectedSection) {
				t.Errorf("%s: expected section <%v>, got <%v>", msg, d.ExpectedSection, got)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown section",
				data{S: s, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user retrieved section",
				data{S: otherS, SID: SID, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"",
				data{S: s, SID: SID, ExpectedSection: section},
			},
		},
	}.Run(t)
}

func TestStorageGetSections(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	SID1, _ := s.NewSection("s1")
	section1, _ := s.GetSection(SID1)

	SID2, _ := s.NewSection("s2")
	section2, _ := s.GetSection(SID2)

	otherU, _ := getTestingUser(t)
	otherS := otherU.Storage()

	otherOtherUser, _ := getTestingUser(t)
	otherOtherUser.Storage().NewSection("s")

	type data struct {
		S Storage

		ExpectedErr      error
		ExpectedSections []Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			sections, err := d.S.GetSections()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(sections, d.ExpectedSections) {
				t.Errorf("%s: expected sections <%v>, got <%v>", msg, d.ExpectedSections, sections)
			}
		},

		Cases: []testCase[data]{
			{
				"got sections of unknown user",
				data{S: unknownUser.Storage(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(empty)",
				data{S: otherS},
			},
			{
				"(filled)",
				data{S: s, ExpectedSections: []Section{section1, section2}},
			},
		},
	}.Run(t)
}

func TestStorageNewSection(t *testing.T) {
	u, _ := getTestingUser(t)
	s := u.Storage()

	type data struct {
		S    Storage
		Name string

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			SID, err := d.S.NewSection(d.Name)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if err == nil {
				section, _ := d.S.GetSection(SID)
				if !reflect.DeepEqual(section, section) {
					t.Errorf("%v, returned bad section", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created section",
				data{S: unknownUser.Storage(), Name: "s", ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{S: s, Name: "s1"},
			},
			{
				"created duplicated section",
				data{S: s, Name: "s1", ExpectedErr: ERR_SECTION_DUPLICATED},
			},
			{
				"",
				data{S: s, Name: "s2"},
			},
		},
	}.Run(t)
}
