package web

import (
	"cucinassistant/web/handlers"
	"cucinassistant/web/utils"
)

var endpoints []utils.Endpoint = []utils.Endpoint{
	{
		Path:       "/",
		GetHandler: handlers.GetIndex,
	},
	{
		Path:        "/info",
		Unprotected: true,
		GetHandler:  handlers.GetInfo,
	},
	{
		Path:        "/lang",
		Unprotected: true,
		GetHandler:  handlers.GetLang,
	},
	{
		Path:        "/stats",
		Unprotected: true,
		GetHandler:  handlers.GetStats,
	},

	{
		Path:       "/menus",
		GetHandler: handlers.GetMenus,
	},
	{
		Path:        "/menus/new",
		GetHandler:  handlers.GetMenusNew,
		PostHandler: handlers.PostMenusNew,
	},
	{
		Path:       "/menus/{MID}",
		GetHandler: handlers.GetMenu,
	},
	{
		Path:        "/menus/{MID}/edit",
		GetHandler:  handlers.GetMenuEdit,
		PostHandler: handlers.PostMenuEdit,
	},
	{
		Path:        "/menus/{MID}/delete",
		PostHandler: handlers.PostMenuDelete,
	},
	{
		Path:        "/menus/{MID}/duplicate",
		PostHandler: handlers.PostMenuDuplicate,
	},

	{
		Path:        "/public_recipes/{code}",
		Unprotected: true,
		GetHandler:  handlers.GetPublicRecipe,
	},
	{
		Path:        "/public_recipes/{code}/save",
		PostHandler: handlers.PostPublicRecipeSave,
	},
	{
		Path:       "/recipes",
		GetHandler: handlers.GetRecipes,
	},
	{
		Path:        "/recipes/new",
		GetHandler:  handlers.GetRecipesNew,
		PostHandler: handlers.PostRecipesNew,
	},
	{
		Path:       "/recipes/{RID}",
		GetHandler: handlers.GetRecipe,
	},
	{
		Path:        "/recipes/{RID}/edit",
		GetHandler:  handlers.GetRecipeEdit,
		PostHandler: handlers.PostRecipeEdit,
	},
	{
		Path:        "/recipes/{RID}/delete",
		PostHandler: handlers.PostRecipeDelete,
	},
	{
		Path:        "/recipes/{RID}/share",
		GetHandler:  handlers.GetRecipeShare,
		PostHandler: handlers.PostRecipeShare,
	},
	{
		Path:        "/recipes/{RID}/unshare",
		PostHandler: handlers.PostRecipeUnshare,
	},

	{
		Path:       "/shopping_list",
		GetHandler: handlers.GetShoppingList,
	},
	{
		Path:        "/shopping_list/append",
		GetHandler:  handlers.GetShoppingListAppend,
		PostHandler: handlers.PostShoppingListAppend,
	},
	{
		Path:        "/shopping_list/clear",
		PostHandler: handlers.PostShoppingListClear,
	},
	{
		Path:        "/shopping_list/{EID}/edit",
		GetHandler:  handlers.GetEntryEdit,
		PostHandler: handlers.PostEntryEdit,
	},
	{
		Path:        "/shopping_list/{EID}/toggle",
		PostHandler: handlers.PostEntryToggle,
	},

	{
		Path:       "/storage",
		GetHandler: handlers.GetStorage,
	},
	{
		Path:        "/storage/new",
		GetHandler:  handlers.GetStorageNew,
		PostHandler: handlers.PostStorageNew,
	},
	{
		Path:       "/storage/{SID}",
		GetHandler: handlers.GetStorageSection,
	},
	{
		Path:        "/storage/{SID}/add",
		GetHandler:  handlers.GetStorageSectionAdd,
		PostHandler: handlers.PostStorageSectionAdd,
	},
	{
		Path:        "/storage/{SID}/delete",
		PostHandler: handlers.PostStorageSectionDelete,
	},
	{
		Path:        "/storage/{SID}/edit",
		GetHandler:  handlers.GetStorageSectionEdit,
		PostHandler: handlers.PostStorageSectionEdit,
	},
	{
		Path:       "/storage/{SID}/search",
		GetHandler: handlers.GetStorageSectionSearch,
	},
	{
		Path:        "/storage/{SID}/{AID}",
		GetHandler:  handlers.GetStorageArticle,
		PostHandler: handlers.PostStorageArticle,
	},
	{
		Path:        "/storage/{SID}/{AID}/delete",
		PostHandler: handlers.PostStorageArticleDelete,
	},

	{
		Path:        "/user/change_email",
		GetHandler:  handlers.GetUserChangeEmail,
		PostHandler: handlers.PostUserChangeEmail,
	},
	{
		Path:        "/user/change_email_settings",
		GetHandler:  handlers.GetUserChangeEmailSettings,
		PostHandler: handlers.PostUserChangeEmailSettings,
	},
	{
		Path:        "/user/change_password",
		GetHandler:  handlers.GetUserChangePassword,
		PostHandler: handlers.PostUserChangePassword,
	},
	{
		Path:        "/user/change_username",
		GetHandler:  handlers.GetUserChangeUsername,
		PostHandler: handlers.PostUserChangeUsername,
	},
	{
		Path:        "/user/delete_1",
		GetHandler:  handlers.GetUserDelete1,
		PostHandler: handlers.PostUserDelete1,
	},
	{
		Path:        "/user/delete_2",
		GetHandler:  handlers.GetUserDelete2,
		PostHandler: handlers.PostUserDelete2,
	},
	{
		Path:        "/user/forgot_password",
		Unprotected: true,
		GetHandler:  handlers.GetForgotPassword,
		PostHandler: handlers.PostForgotPassword,
	},
	{
		Path:        "/user/reset_password",
		Unprotected: true,
		GetHandler:  handlers.GetResetPassword,
		PostHandler: handlers.PostResetPassword,
	},
	{
		Path:       "/user/settings",
		GetHandler: handlers.GetUserSettings,
	},
	{
		Path:        "/user/signin",
		Unprotected: true,
		GetHandler:  handlers.GetUserSignIn,
		PostHandler: handlers.PostUserSignIn,
	},
	{
		Path:        "/user/signout",
		Unprotected: true,
		PostHandler: handlers.PostUserSignOut,
	},
	{
		Path:        "/user/signup",
		Unprotected: true,
		GetHandler:  handlers.GetUserSignUp,
		PostHandler: handlers.PostUserSignUp,
	},
}
