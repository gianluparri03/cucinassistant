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
		Path:       "/info",
		GetHandler: handlers.GetInfo,
	},

	{
		Path:       "/menus",
		GetHandler: handlers.GetMenus,
	},
	{
		Path:        "/menus/new",
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
		PostHandler: handlers.PostClearEntries,
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
		PostHandler: handlers.PostSignOut,
	},
	{
		Path:        "/user/forgot_password",
		GetHandler:  handlers.GetForgotPassword,
		PostHandler: handlers.PostForgotPassword,
	},
	{
		Path:        "/user/reset_password",
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
		Path:        "/user/delete_1",
		GetHandler:  handlers.GetDeleteUser1,
		PostHandler: handlers.PostDeleteUser1,
	},
	{
		Path:        "/user/delete_2",
		GetHandler:  handlers.GetDeleteUser2,
		PostHandler: handlers.PostDeleteUser2,
	},
}
