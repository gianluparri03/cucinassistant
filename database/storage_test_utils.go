package database

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
