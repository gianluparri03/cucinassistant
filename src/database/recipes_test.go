package database

import (
	"reflect"
	"testing"
)

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

	newDataSameName := Recipe{RID: RID, Name: "oldName"}
	newDataTakenName := Recipe{RID: RID, Name: "takenName"}
	newData := Recipe{RID: RID, Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-"}

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
				if !reflect.DeepEqual(got, d.NewData) {
					t.Errorf("%s: expected <%v>, got <%v>", msg, d.NewData, got)
				} else if r, _ := d.R.GetOne(d.RID); !reflect.DeepEqual(r, got) {
					t.Errorf("%v, changes not saved", msg)
				}
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
				data{R: r, RID: RID, NewData: newDataSameName},
			},
			{
				"(taken name)",
				data{R: r, RID: RID, NewData: newDataTakenName, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{R: r, RID: RID, NewData: newData},
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
		Name: "r1", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
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
			} else if !reflect.DeepEqual(recipes, d.ExpectedRecipes) {
				t.Errorf("%s: expected recipes <%v>, got <%v>", msg, d.ExpectedRecipes, recipes)
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
		Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
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
			} else if !reflect.DeepEqual(recipe, d.ExpectedRecipe) {
				t.Errorf("%s: expected recipe <%v>, got <%v>", msg, d.ExpectedRecipe, recipe)
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
		Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
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
			} else if !reflect.DeepEqual(recipe, d.ExpectedRecipe) {
				t.Errorf("%s: expected recipe <%v>, got <%v>", msg, d.ExpectedRecipe, recipe)
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
				if !reflect.DeepEqual(recipe, d.ExpectedRecipe) {
					t.Errorf("%v, saved recipe not the same", msg)
				}
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
