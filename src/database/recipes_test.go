package database

import (
	"reflect"
	"testing"
)

func TestRecipesNew(t *testing.T) {
	user, _ := getTestingUser(t)
	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		Name string

		ExpectedErr  error
		ExpectedName string
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			name := "testRecipe"

			if r, err := d.User.Recipes().New(name); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if err == nil {
				if r.Name != name {
					t.Errorf("%s: expected name <%s>, got <%v>", msg, name, r.Name)
				} else if _, err = d.User.Recipes().GetOne(r.RID); err != nil {
					t.Errorf("%s: recipe not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"unknown user created recipe",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"",
				data{User: user},
			},
			{
				"created duplicated recipe",
				data{User: user, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{User: otherUser},
			},
		},
	}.Run(t)
}

func TestRecipesGetAll(t *testing.T) {
	user, _ := getTestingUser(t)
	r1, _ := user.Recipes().New("r1")
	r2, _ := user.Recipes().New("r2")
	user.Recipes().Edit(r1.RID, Recipe{
		Name: "r1", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
	})

	userWithoutRecipes, _ := getTestingUser(t)

	userWithRecipies, _ := getTestingUser(t)
	userWithRecipies.Recipes().New("r")

	type data struct {
		User User

		ExpectedErr     error
		ExpectedRecipes []Recipe
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			recipes, err := d.User.Recipes().GetAll()
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(recipes, d.ExpectedRecipes) {
				t.Errorf("%s: expected recipes <%v>, got <%v>", msg, d.ExpectedRecipes, recipes)
			}
		},

		Cases: []testCase[data]{
			{
				"got recipies of unknown user",
				data{User: unknownUser, ExpectedErr: ERR_USER_UNKNOWN},
			},
			{
				"(no recipies)",
				data{User: userWithoutRecipes},
			},
			{
				"(some recipies)",
				data{User: user, ExpectedRecipes: []Recipe{r1, r2}},
			},
		},
	}.Run(t)
}

func TestRecipesGetOne(t *testing.T) {
	user, _ := getTestingUser(t)
	recipe, _ := user.Recipes().New("r")
	recipe, _ = user.Recipes().Edit(recipe.RID, Recipe{
		Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-",
	})

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		RID  int

		ExpectedErr    error
		ExpectedRecipe Recipe
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			recipe, err := d.User.Recipes().GetOne(d.RID)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			} else if !reflect.DeepEqual(recipe, d.ExpectedRecipe) {
				t.Errorf("%s: expected recipe <%v>, got <%v>", msg, d.ExpectedRecipe, recipe)
			}
		},

		Cases: []testCase[data]{
			{
				"got data of unkown recipe",
				data{User: user, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"other user retrieved recipe",
				data{User: otherUser, RID: recipe.RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"",
				data{User: user, RID: recipe.RID, ExpectedRecipe: recipe},
			},
		},
	}.Run(t)
}

func TestRecipesEdit(t *testing.T) {
	user, _ := getTestingUser(t)
	recipe, _ := user.Recipes().New("oldName")
	user.Recipes().New("takenName")

	newDataSameName := Recipe{RID: recipe.RID, Name: "oldName"}
	newDataTakenName := Recipe{RID: recipe.RID, Name: "takenName"}
	newData := Recipe{RID: recipe.RID, Name: "newName", Stars: 4, Ingredients: "flour", Directions: "Mix", Notes: "-"}

	otherUser, _ := getTestingUser(t)
	otherUser.Recipes().New("newName")

	type data struct {
		User User
		RID  int

		NewData Recipe

		ExpectedErr error
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			got, err := d.User.Recipes().Edit(d.RID, d.NewData)
			if err != d.ExpectedErr {
				t.Errorf("%s: expected <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			if d.ExpectedErr == nil {
				if !reflect.DeepEqual(got, d.NewData) {
					t.Errorf("%s: expected <%v>, got <%v>", msg, d.NewData, got)
				} else if r, _ := d.User.Recipes().GetOne(d.RID); !reflect.DeepEqual(r, got) {
					t.Errorf("%v, changes not saved", msg)
				}
			}
		},

		Cases: []testCase[data]{
			{
				"other user edited recipe",
				data{User: otherUser, RID: recipe.RID, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"replaced unknown recipe",
				data{User: user, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"(same name)",
				data{User: user, RID: recipe.RID, NewData: newDataSameName},
			},
			{
				"(taken name)",
				data{User: user, RID: recipe.RID, NewData: newDataTakenName, ExpectedErr: ERR_RECIPE_DUPLICATED},
			},
			{
				"",
				data{User: user, RID: recipe.RID, NewData: newData},
			},
		},
	}.Run(t)
}

func TestRecipesDelete(t *testing.T) {
	user, _ := getTestingUser(t)
	recipe, _ := user.Recipes().New("")

	otherUser, _ := getTestingUser(t)

	type data struct {
		User User
		RID  int

		ExpectedErr error
		ShouldExist bool
	}

	testSuite[data]{
		Target: func(t *testing.T, msg string, d data) {
			if err := d.User.Recipes().Delete(d.RID); err != d.ExpectedErr {
				t.Errorf("%s: expected err <%v>, got <%v>", msg, d.ExpectedErr, err)
			}

			recipe, _ = user.Recipes().GetOne(d.RID)
			if !d.ShouldExist && recipe.RID != 0 {
				t.Errorf("%s, recipe wasn't deleted", msg)
			} else if d.ShouldExist && recipe.RID == 0 {
				t.Errorf("%s, recipe was deleted anyway", msg)
			}
		},

		Cases: []testCase[data]{
			{
				"deleted unknown recipe",
				data{User: user, ExpectedErr: ERR_RECIPE_NOT_FOUND},
			},
			{
				"other user deleted recipe",
				data{User: otherUser, RID: recipe.RID, ExpectedErr: ERR_RECIPE_NOT_FOUND, ShouldExist: true},
			},
			{
				"",
				data{User: user, RID: recipe.RID},
			},
		},
	}.Run(t)
}
