package database

import (
	"reflect"
	"testing"
)

var testingArticlesN int = 0

// getExpectedArticle returns an Article with
// the expected AID
func (sa StringArticle) getExpectedArticle() Article {
	a, _ := sa.Parse()
	a.fixExpiration()

	testingArticlesN++
	a.AID = testingArticlesN
	return a
}

func TestGetSections(t *testing.T) {
	user, _ := GetTestingUser(t)
	s1, _ := user.NewSection("s1")
	s2, _ := user.NewSection("s2")

	otherUser, _ := GetTestingUser(t)

	otherOtherUser, _ := GetTestingUser(t)
	otherOtherUser.NewSection("s")

	type data struct {
		User User

		ExpectedErr      error
		ExpectedSections []Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			sections, err := d.User.GetSections()
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

func TestGetSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s")

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User User
		SID  int

		ExpectedErr     error
		ExpectedSection Section
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.GetSection(d.SID)
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

func TestNewSection(t *testing.T) {
	user, _ := GetTestingUser(t)

	type data struct {
		User User
		Name string

		ExpectedErr error
	}

	testSuite[data]{
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

func TestEditSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s1")
	user.NewSection("s2")

	otherUser, _ := GetTestingUser(t)
	otherUser.NewSection("s3")

	type data struct {
		User    User
		SID     int
		NewName string

		ExpectedErr error
	}

	testSuite[data]{
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

func TestDeleteSection(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("s")
	user.AddArticles(section.SID, StringArticle{"article", "", ""})
	testingArticlesN++

	otherUser, _ := GetTestingUser(t)

	type data struct {
		User User
		SID  int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
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

func TestAddArticles(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("section")
	SID := section.SID

	otherUser, _ := GetTestingUser(t)

	name := "article"
	qty := "10"
	exp := "2024-10-05"

	inList := []StringArticle{
		{Name: "NoQty", Quantity: "", Expiration: exp},
		{Name: "Full", Quantity: qty, Expiration: exp},
		{Name: "Empty", Quantity: "", Expiration: ""},
		{Name: "NoExp", Quantity: qty, Expiration: ""},
	}

	outList1 := []Article{}
	outList2 := []Article{}
	for _, sa := range inList {
		original := sa.getExpectedArticle()
		outList1 = append(outList1, original)

		doubled := Article{AID: original.AID, Name: original.Name, Expiration: original.Expiration}
		if original.Quantity == nil {
			doubled.Quantity = nil
		} else {
			qty := (*original.Quantity) * 2
			doubled.Quantity = &qty
		}
		outList2 = append(outList2, doubled)
	}

	testingArticlesN += 4

	type data struct {
		User     User
		SID      int
		Articles []StringArticle

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
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

		Cases: []testCase[data]{
			{
				"added articles to unknown section",
				data{User: user, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"other user added articles to section",
				data{User: otherUser, SID: SID, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: exp}}, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"added article with invalid quantity",
				data{User: user, SID: SID, Articles: []StringArticle{{Name: name, Quantity: "a lot", Expiration: exp}}, ExpectedErr: ERR_ARTICLE_QUANTITY_INVALID},
			},
			{
				"added article with invalid sExpiration",
				data{User: user, SID: SID, Articles: []StringArticle{{Name: name, Quantity: qty, Expiration: "dunno"}}, ExpectedErr: ERR_ARTICLE_EXPIRATION_INVALID},
			},
			{
				"(original)",
				data{User: user, SID: SID, Articles: inList, ExpectedArticles: outList1},
			},
			{
				"(doubled)",
				data{User: user, SID: SID, Articles: inList, ExpectedArticles: outList2},
			},
		},
	}.Run(t)
}

func TestGetArticles(t *testing.T) {
	user, _ := GetTestingUser(t)

	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2"}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11"}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15"}
	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()

	section, _ := user.NewSection("section")
	user.AddArticles(section.SID, s1, s2, s3)
	expected := []Article{a2, a1, a3}

	otherSection, _ := user.NewSection("otherSection")
	user.AddArticles(otherSection.SID, s1)
	testingArticlesN++

	otherUser, _ := GetTestingUser(t)
	emptySection, _ := otherUser.NewSection("emptySection")

	type data struct {
		User   User
		SID    int
		Filter string

		ExpectedErr      error
		ExpectedArticles []Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			section, err := d.User.GetArticles(d.SID, d.Filter)
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
				data{User: user, ExpectedErr: ERR_SECTION_NOT_FOUND},
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

func TestGetArticle(t *testing.T) {
	user, _ := GetTestingUser(t)

	stringEmpty := StringArticle{Name: "empty"}
	stringFull := StringArticle{Name: "full", Expiration: "2024-06-02", Quantity: "900"}
	empty := stringEmpty.getExpectedArticle()
	full := stringFull.getExpectedArticle()

	section, _ := user.NewSection("section")
	user.AddArticles(section.SID, stringEmpty, stringFull)

	otherSection, _ := user.NewSection("otherSection")
	otherUser, _ := GetTestingUser(t)

	type data struct {
		User User
		SID  int
		AID  int

		ExpectedErr     error
		ExpectedArticle Article
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			article, err := d.User.GetArticle(d.SID, d.AID)
			if err != d.ExpectedErr {
				t.Errorf("(simple) %s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(article, d.ExpectedArticle) {
				t.Errorf("(simple) %s: expected article <%v>, got <%v>", msg, d.ExpectedArticle, article)
			}

			orderedArticle, err := d.User.GetOrderedArticle(d.SID, d.AID)
			article = orderedArticle.GetArticle()
			if err != d.ExpectedErr {
				t.Errorf("(ordered) %s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(article, d.ExpectedArticle) {
				t.Errorf("(ordered) %s: expected article <%v>, got <%v>", msg, d.ExpectedArticle, article)
			}
		},

		Cases: []testCase[data]{
			{
				"got unknown article",
				data{User: user, SID: section.SID, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"got article from wrong section",
				data{User: user, SID: otherSection.SID, AID: full.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND},
			},
			{
				"other user retrieved article",
				data{User: otherUser, SID: section.SID, AID: full.AID, ExpectedErr: ERR_SECTION_NOT_FOUND},
			},
			{
				"(filled)",
				data{User: user, SID: section.SID, AID: full.AID, ExpectedArticle: full},
			},
			{
				"(empty)",
				data{User: user, SID: section.SID, AID: empty.AID, ExpectedArticle: empty},
			},
		},
	}.Run(t)
}

func TestGetOrderedArticles(t *testing.T) {
	user, _ := GetTestingUser(t)

	section, _ := user.NewSection("section")
	s1 := StringArticle{Name: "first", Expiration: "2024-11-08", Quantity: "2"}
	s2 := StringArticle{Name: "second", Expiration: "2024-10-11"}
	s3 := StringArticle{Name: "third", Expiration: "", Quantity: "15"}
	user.AddArticles(section.SID, s1, s2, s3)

	otherSection, _ := user.NewSection("otherSection")
	s4 := StringArticle{Name: "middle", Expiration: "2024-10-31"}
	user.AddArticles(otherSection.SID, s4)

	a1 := s1.getExpectedArticle()
	a2 := s2.getExpectedArticle()
	a3 := s3.getExpectedArticle()
	a4 := s4.getExpectedArticle()

	type data struct {
		User User
		SID  int
		AID  int

		ExpectedPrev *int
		ExpectedNext *int
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			oa, _ := user.GetOrderedArticle(d.SID, d.AID)
			if !reflect.DeepEqual(oa.Prev, d.ExpectedPrev) {
				t.Errorf("%s: expected prev <%p>, got <%p>", msg, d.ExpectedPrev, oa.Prev)
			} else if !reflect.DeepEqual(oa.Next, d.ExpectedNext) {
				t.Errorf("%s: expected next <%p>, got <%p>", msg, d.ExpectedNext, oa.Next)
			}
		},

		Cases: []testCase[data]{
			{
				"",
				data{SID: section.SID, AID: a1.AID, ExpectedPrev: &a2.AID, ExpectedNext: &a3.AID},
			},
			{
				"",
				data{SID: section.SID, AID: a2.AID, ExpectedPrev: nil, ExpectedNext: &a1.AID},
			},
			{
				"",
				data{SID: section.SID, AID: a3.AID, ExpectedPrev: &a1.AID, ExpectedNext: nil},
			},
			{
				"",
				data{SID: otherSection.SID, AID: a4.AID, ExpectedPrev: nil, ExpectedNext: nil},
			},
		},
	}.Run(t)
}

func TestDeleteArticle(t *testing.T) {
	user, _ := GetTestingUser(t)
	section, _ := user.NewSection("section")

	sArticle := StringArticle{Name: "article", Expiration: "2024-12-18", Quantity: "5"}
	article := sArticle.getExpectedArticle()
	sNextArticle := StringArticle{Name: "next-article"}
	nextArticle := sNextArticle.getExpectedArticle()
	user.AddArticles(section.SID, sArticle, sNextArticle)

	sSimilarArticle := StringArticle{Name: "article", Expiration: "2023-02-18"}
	user.AddArticles(section.SID, sSimilarArticle)
	similarArticle := sSimilarArticle.getExpectedArticle()

	otherSection, _ := user.NewSection("otherSection")
	otherUser, _ := GetTestingUser(t)

	type data struct {
		User User
		SID  int
		AID  int

		ExpectedErr  error
		ExpectedNext *int
		ShouldExist  bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err, next := d.User.DeleteArticle(d.SID, d.AID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(next, d.ExpectedNext) {
				t.Errorf("%s: expected next <%p>, got <%p>", msg, d.ExpectedNext, next)
			} else if d.AID != similarArticle.AID {
				if _, err := user.GetArticle(section.SID, similarArticle.AID); err != nil {
					t.Errorf("%s, deleted similar article", msg)
				}
			}

			article, _ := user.GetArticle(section.SID, d.AID)
			if !d.ShouldExist && article.AID != 0 {
				t.Errorf("%s, article wasn't deleted", msg)
			} else if d.ShouldExist && article.AID == 0 {
				t.Errorf("%s, article was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"other user deleted article",
				data{User: otherUser, SID: section.SID, AID: article.AID, ExpectedErr: ERR_SECTION_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted article from other section",
				data{User: user, SID: otherSection.SID, AID: article.AID, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: true},
			},
			{
				"deleted unknown article",
				data{User: user, SID: section.SID, ExpectedErr: ERR_ARTICLE_NOT_FOUND, ShouldExist: false},
			},
			{
				"(next=next)",
				data{User: user, SID: section.SID, AID: article.AID, ShouldExist: false, ExpectedNext: &nextArticle.AID},
			},
			{
				"(next=prev)",
				data{User: user, SID: section.SID, AID: nextArticle.AID, ShouldExist: false, ExpectedNext: &similarArticle.AID},
			},
			{
				"(next=nil)",
				data{User: user, SID: section.SID, AID: similarArticle.AID, ShouldExist: false, ExpectedNext: nil},
			},
		},
	}.Run(t)
}
