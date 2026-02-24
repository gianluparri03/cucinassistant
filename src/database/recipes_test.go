package database

import (
	"reflect"
	"slices"
	"testing"
)

func compareRecipes(t *testing.T, msg string, expected, got Recipe) {
	if got.Name != expected.Name {
		t.Errorf("%s: expected name <%s>, got <%s>", msg, expected.Name, got.Name)
	}
	if got.Stars != expected.Stars {
		t.Errorf("%s: expected stars <%d>, got <%d>", msg, expected.Stars, got.Stars)
	}
	if got.Ingredients != expected.Ingredients {
		t.Errorf("%s: expected ingredients <%s>, got <%s>", msg, expected.Ingredients, got.Ingredients)
	}
	if got.Directions != expected.Directions {
		t.Errorf("%s: expected directions <%s>, got <%s>", msg, expected.Directions, got.Directions)
	}
	if !reflect.DeepEqual(got.Code, expected.Code) {
		t.Errorf("%s: expected code <%v>, got <%v>", msg, expected.Code, got.Code)
	}
	if got.Notes != expected.Notes {
		t.Errorf("%s: expected notes <%s>, got <%s>", msg, expected.Notes, got.Notes)
	}

	for _, tag := range expected.Tags {
		if !slices.Contains(got.Tags, tag) {
			t.Errorf("%s: missing tag <%s>", msg, tag)
		}
	}

	for _, tag := range got.Tags {
		if !slices.Contains(expected.Tags, tag) {
			t.Errorf("%s: unwanted tag <%s>", msg, tag)
		}
	}
}

func compareRecipesList(t *testing.T, msg string, expected, got []Recipe) {
	if len(expected) != len(got) {
		t.Errorf("%s: wrong number of recipes: expected <%d>, got <%d>", msg, len(expected), len(got))
		return
	}

	for _, recipe1 := range expected {
		found2 := false
		for _, recipe2 := range got {
			if recipe1.RID == recipe2.RID {
				found2 = true
				compareRecipes(t, msg+", recipe <"+recipe1.Name+">", recipe1, recipe2)
				break
			}
		}

		if !found2 {
			t.Errorf("%s: recipe not found: name <%s>", msg, recipe1.Name)
			return
		}
	}
}

func TestRecipesDelete(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("")
	recipe, _ := r.GetOne(RID)

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()

	type data struct {
		R   Recipes
		RID int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.R.Delete(d.RID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			recipe, _ = r.GetOne(d.RID)
			if !d.ShouldExist && recipe.RID != 0 {
				t.Errorf("%s, recipe wasn't deleted", msg)
			} else if d.ShouldExist && recipe.RID == 0 {
				t.Errorf("%s, recipe was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"deleted unknown recipe",
				data{R: r, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"other user deleted recipe",
				data{R: otherR, RID: RID, ExpectedErr: ERR_RECIPE_NOT_FOUND, ShouldExist: true},
			},
			{
				"",
				data{R: r, RID: RID},
			},
		},
	}.Run(t)
}

func TestRecipesEdit(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("oldName")
	r.New("takenName")

	newData := Recipe{Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-"}
	newDataWithTags := Recipe{Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"vegan", "gluten free"}}
	newDataWithOneTag := Recipe{Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"vegan"}}

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()
	otherR.New("newName")

	type data struct {
		R       Recipes
		RID     int
		NewData Recipe

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.R.Edit(d.RID, d.NewData)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			got, _ := d.R.GetOne(d.RID)

			if d.ExpectedErr == nil {
				compareRecipes(t, msg, d.NewData, got)
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited recipe",
				data{R: otherR, RID: RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"replaced unknown recipe",
				data{R: r, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"(same name)",
				data{R: r, RID: RID, NewData: Recipe{Name: "oldName"}},
			},
			{
				"taken name",
				data{R: r, RID: RID, NewData: Recipe{Name: "takenName"}, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{R: r, RID: RID, NewData: newData},
			},
			{
				"(added tags)",
				data{R: r, RID: RID, NewData: newDataWithTags},
			},
			{
				"(removed tag)",
				data{R: r, RID: RID, NewData: newDataWithOneTag},
			},
		},
	}.Run(t)
}

func TestRecipesGetAll(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID1, _ := r.New("r1")
	recipe1, _ := r.GetOne(RID1)
	r.Edit(RID1, Recipe{
		Name: "r1", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"vegan"},
	})

	RID2, _ := r.New("r2")
	recipe2, _ := r.GetOne(RID2)

	uWithout, _ := getTestingUser(t)
	rWithout := uWithout.Recipes()

	uWith, _ := getTestingUser(t)
	rWith := uWith.Recipes()
	rWith.New("r")

	type data struct {
		R Recipes

		ExpectedErr     error
		ExpectedRecipes []Recipe
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			recipes, err := d.R.GetAll()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				compareRecipesList(t, msg, d.ExpectedRecipes, recipes)
			}
		},

		Cases: []testCase[data]{
			{
				"got recipes of unknown user",
				data{R: unknownUser.Recipes(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(no recipes)",
				data{R: rWithout},
			},
			{
				"(some recipes)",
				data{R: r, ExpectedRecipes: []Recipe{recipe1, recipe2}},
			},
		},
	}.Run(t)
}

func TestRecipesGetOne(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("r")
	r.Edit(RID, Recipe{
		Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"gluten free"},
	})
	recipe, _ := r.GetOne(RID)

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()

	type data struct {
		R   Recipes
		RID int

		ExpectedErr    error
		ExpectedRecipe Recipe
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			recipe, err := d.R.GetOne(d.RID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				compareRecipes(t, msg, d.ExpectedRecipe, recipe)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unkown recipe",
				data{R: r, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"other user retrieved recipe",
				data{R: otherR, RID: RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"",
				data{R: r, RID: RID, ExpectedRecipe: recipe},
			},
		},
	}.Run(t)
}

func TestRecipesGetPublic(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("r")
	r.Edit(RID, Recipe{
		Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"dairy free"},
	})
	code, _ := r.Share(RID)
	recipe, _ := r.GetOne(RID)

	type data struct {
		Code string

		ExpectedErr    error
		ExpectedRecipe Recipe
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			d.ExpectedRecipe.RID = 0
			recipe, err := GetPublicRecipe(d.Code)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else {
				compareRecipes(t, msg, d.ExpectedRecipe, recipe)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unknown recipe",
				data{ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"",
				data{Code: code, ExpectedRecipe: recipe},
			},
		},
	}.Run(t)
}

func TestRecipesGetTags(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID2, _ := r.New("r2")
	recipe2, _ := r.GetOne(RID2)
	r.Edit(RID2, Recipe{
		Name: "r2", Stars: 2, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"gluten free"},
	})

	RID3, _ := r.New("r3")
	recipe3, _ := r.GetOne(RID3)
	r.Edit(RID3, Recipe{
		Name: "r3", Stars: 3, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"vegan", "gluten free"},
	})

	RID4, _ := r.New("r4")
	r.Edit(RID4, Recipe{
		Name: "r4", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
	})

	RID1, _ := r.New("r1")
	recipe1, _ := r.GetOne(RID1)
	r.Edit(RID1, Recipe{
		Name: "r1", Stars: 1, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"vegan"},
	})

	tags := []Tag{
		{Name: "gluten free", Recipes: []Recipe{recipe2, recipe3}},
		{Name: "vegan", Recipes: []Recipe{recipe1, recipe3}},
	}

	uWithout, _ := getTestingUser(t)
	rWithout := uWithout.Recipes()

	type data struct {
		R Recipes

		ExpectedErr  error
		ExpectedTags []Tag
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			tags, err := d.R.GetTags()

			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
				return
			} else if len(d.ExpectedTags) != len(tags) {
				t.Errorf("%s: wrong number of tags: expected <%d>, got <%d>", msg, len(d.ExpectedTags), len(tags))
				return
			}

			for _, tag1 := range d.ExpectedTags {
				found := false
				for _, tag2 := range tags {
					if tag1.Name == tag2.Name {
						found = true
						compareRecipesList(t, msg+", tag <"+tag1.Name+">", tag1.Recipes, tag2.Recipes)
						break
					}
				}

				if !found {
					t.Errorf("%s: tag <%s> not found", msg, tag1.Name)
					return
				}
			}
		},

		Cases: []testCase[data]{
			{
				"got recipes of unknown user",
				data{R: unknownUser.Recipes(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(no recipes)",
				data{R: rWithout},
			},
			{
				"(some recipes)",
				data{R: r, ExpectedTags: tags},
			},
		},
	}.Run(t)
}

func TestRecipesNew(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()

	type data struct {
		R    Recipes
		Name string

		ExpectedErr  error
		ExpectedName string
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			name := "testRecipe"

			if RID, err := d.R.New(name); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				recipe, _ := d.R.GetOne(RID)

				if recipe.Name != name {
					t.Errorf("%s: expected name <%s>, got <%v>", msg, name, recipe.Name)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created recipe",
				data{R: unknownUser.Recipes(), ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{R: r},
			},
			{
				"created duplicated recipe",
				data{R: r, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{R: otherR},
			},
		},
	}.Run(t)
}

func TestRecipesSave(t *testing.T) {
	ownerU, _ := getTestingUser(t)
	ownerR := ownerU.Recipes()

	RID, _ := ownerR.New("recipe")
	ownerR.Edit(RID, Recipe{
		Name: "recipe", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-", Tags: []string{"dairy free"},
	})
	code, _ := ownerR.Share(RID)
	recipe, _ := ownerR.GetOne(RID)

	u, _ := getTestingUser(t)
	r := u.Recipes()

	u2, _ := getTestingUser(t)
	r2 := u2.Recipes()
	r2.New("recipe")

	type data struct {
		R    Recipes
		Code string

		ExpectedRecipe Recipe
		ExpectedErr    error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			RID, err := d.R.Save(d.Code)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				recipe, _ := d.R.GetOne(RID)
				d.ExpectedRecipe.RID = RID
				d.ExpectedRecipe.Code = nil
				compareRecipes(t, msg, d.ExpectedRecipe, recipe)
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user saved recipe",
				data{Code: code, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"saved unknown recipe",
				data{R: r, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"duplicated recipe (owner)",
				data{R: ownerR, Code: code, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"duplicated recipe (other user)",
				data{R: r2, Code: code, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{R: r, Code: code, ExpectedRecipe: recipe},
			},
		},
	}.Run(t)
}

func TestRecipesShare(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("recipe")

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()

	type data struct {
		R   Recipes
		RID int

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			code, err := d.R.Share(d.RID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				if r, _ := d.R.GetOne(d.RID); *r.Code != code {
					t.Errorf("%v, code not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user shared recipe",
				data{R: otherR, RID: RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"(new)",
				data{R: r, RID: RID},
			},
			{
				"(replace)",
				data{R: r, RID: RID},
			},
		},
	}.Run(t)
}

func TestRecipesUnshare(t *testing.T) {
	u, _ := getTestingUser(t)
	r := u.Recipes()

	RID, _ := r.New("recipe")
	recipe, _ := r.GetOne(RID)
	r.Share(recipe.RID)

	otherU, _ := getTestingUser(t)
	otherR := otherU.Recipes()
	otherRID, _ := otherR.New("otherRecipe")

	type data struct {
		R   Recipes
		RID int

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			err := d.R.Unshare(d.RID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				if r, _ := d.R.GetOne(d.RID); r.Code != nil {
					t.Errorf("%v, code not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user unshared recipe",
				data{R: otherR, RID: RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"(with code)",
				data{R: r, RID: RID},
			},
			{
				"(without code)",
				data{R: otherR, RID: otherRID},
			},
		},
	}.Run(t)
}
