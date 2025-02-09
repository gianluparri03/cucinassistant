package web

import (
	"cucinassistant/web/handlers"
	"cucinassistant/web/utils"
)

var endpoints []utils.Endpoint = []utils.Endpoint{
	{
		Path:        "/lang",
		Unprotected: true,
		GetHandler:  handlers.GetLang,
	},
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
		GetHandler:  handlers.GetNewMenu,
		PostHandler: handlers.PostNewMenu,
	},
	{
		Path:       "/menus/{MID}",
		GetHandler: handlers.GetMenu,
	},
	{
		Path:        "/menus/{MID}/edit",
		GetHandler:  handlers.GetEditMenu,
		PostHandler: handlers.PostEditMenu,
	},
	{
		Path:        "/menus/{MID}/duplicate",
		PostHandler: handlers.PostDuplicateMenu,
	},
	{
		Path:        "/menus/{MID}/delete",
		PostHandler: handlers.PostDeleteMenu,
	},

	{
		Path:       "/shopping_list",
		GetHandler: handlers.GetShoppingList,
	},
	{
		Path:        "/shopping_list/append",
		GetHandler:  handlers.GetAppendEntries,
		PostHandler: handlers.PostAppendEntries,
	},
	{
		Path:        "/shopping_list/{EID}/toggle",
		PostHandler: handlers.PostToggleEntry,
	},
	{
		Path:        "/shopping_list/{EID}/edit",
		GetHandler:  handlers.GetEditEntry,
		PostHandler: handlers.PostEditEntry,
	},
	{
		Path:        "/shopping_list/clear",
		PostHandler: handlers.PostClearShoppingList,
	},
	{
		Path:        "/user/signup",
		Unprotected: true,
		GetHandler:  handlers.GetSignUp,
		PostHandler: handlers.PostSignUp,
	},
	{
		Path:        "/user/signin",
		Unprotected: true,
		GetHandler:  handlers.GetSignIn,
		PostHandler: handlers.PostSignIn,
	},
	{
		Path:        "/user/signout",
		Unprotected: true,
		PostHandler: handlers.PostSignOut,
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
		GetHandler: handlers.GetSettings,
	},
	{
		Path:        "/user/change_username",
		GetHandler:  handlers.GetChangeUsername,
		PostHandler: handlers.PostChangeUsername,
	},
	{
		Path:        "/user/change_email",
		GetHandler:  handlers.GetChangeEmail,
		PostHandler: handlers.PostChangeEmail,
	},
	{
		Path:        "/user/change_password",
		GetHandler:  handlers.GetChangePassword,
		PostHandler: handlers.PostChangePassword,
	},
	{
		Path:        "/user/set_email_lang",
		GetHandler:  handlers.GetSetEmailLang,
		PostHandler: handlers.PostSetEmailLang,
	},
	{
		Path:        "/user/delete_1",
		GetHandler:  handlers.GetDeleteUser1,
		PostHandler: handlers.PostDeleteUser1,
	},
	{
		Path:        "/user/delete_2",
		GetHandler:  handlers.GetDeleteUser2,
		PostHandler: handlers.PostDeleteUser2,
	},
	{
		Path:       "/storage",
		GetHandler: handlers.GetSections,
	},
	{
		Path:        "/storage/new",
		GetHandler:  handlers.GetNewSection,
		PostHandler: handlers.PostNewSection,
	},
	{
		Path:        "/storage/add",
		GetHandler:  handlers.GetAddArticlesCommon,
		PostHandler: handlers.PostAddArticlesCommon,
	},
	{
		Path:       "/storage/{SID}",
		GetHandler: handlers.GetArticles,
	},
	{
		Path:        "/storage/{SID}/edit",
		GetHandler:  handlers.GetEditSection,
		PostHandler: handlers.PostEditSection,
	},
	{
		Path:        "/storage/{SID}/delete",
		PostHandler: handlers.PostDeleteSection,
	},
	{
		Path:        "/storage/{SID}/add",
		GetHandler:  handlers.GetAddArticlesSection,
		PostHandler: handlers.PostAddArticlesSection,
	},
	{
		Path:       "/storage/{SID}/search",
		GetHandler: handlers.GetSearchArticles,
	},
	{
		Path:        "/storage/{SID}/{AID}",
		GetHandler:  handlers.GetEditArticle,
		PostHandler: handlers.PostEditArticle,
	},
	{
		Path:        "/storage/{SID}/{AID}/delete",
		PostHandler: handlers.PostDeleteArticle,
	},
	{
		Path:       "/recipes",
		GetHandler: handlers.GetRecipes,
	},
	{
		Path:       "/recipes/new",
		GetHandler: handlers.GetNewRecipe,
	},
	{
		Path:       "/recipes/{RID}",
		GetHandler: handlers.GetRecipe,
	},
	{
		Path:       "/recipes/{RID}/edit",
		GetHandler: handlers.GetEditRecipe,
	},
}
