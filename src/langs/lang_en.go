package langs

import (
	"cucinassistant/database"
)

var english *Lang = &Lang{
	Tag:  "en",
	Name: "English",

	Strings: map[String]string{
		String(database.ERR_ARTICLE_DUPLICATED):         "An article with this name and expiration already exists",
		String(database.ERR_ARTICLE_EXPIRATION_INVALID): "Invalid expiration",
		String(database.ERR_ARTICLE_NOT_FOUND):          "Article not found",
		String(database.ERR_ARTICLE_QUANTITY_INVALID):   "Invalid quantity",
		String(database.ERR_ENTRY_DUPLICATED):           "Entry already in list",
		String(database.ERR_ENTRY_NOT_FOUND):            "Entry not found",
		String(database.ERR_MENU_NOT_FOUND):             "Menu not found",
		String(database.ERR_RECIPE_DUPLICATED):          "A recipe with this name already exists",
		String(database.ERR_RECIPE_NOT_FOUND):           "Recipe not found",
		String(database.ERR_SECTION_DUPLICATED):         "A section with this name already exists",
		String(database.ERR_SECTION_NOT_FOUND):          "Section not found",
		String(database.ERR_UNKNOWN):                    "Unknown error",
		String(database.ERR_USER_MAIL_INVALID):          "Invalid email",
		String(database.ERR_USER_MAIL_UNAVAIL):          "Email not available",
		String(database.ERR_USER_NAME_TOO_SHORT):        "Invalid username: must be at least 5 characters long",
		String(database.ERR_USER_NAME_UNAVAIL):          "Username not available",
		String(database.ERR_USER_PASS_TOO_SHORT):        "Invalid password: must be at least 8 characters long",
		String(database.ERR_USER_UNKNOWN):               "Unknown user",
		String(database.ERR_USER_WRONG_CREDENTIALS):     "Wrong credentials",
		String(database.ERR_USER_WRONG_TOKEN):           "Something went wrong",
		STR_ADD:                                         "Add",
		STR_ADD_ARTICLES:                                "Add articles",
		STR_APPEND_ENTRIES:                              "Add entries",
		STR_ARTICLES_ADDED:                              "Articles added succesfully",
		STR_CANCEL:                                      "Cancel",
		STR_CHANGE_EMAIL:                                "Email change",
		STR_CHANGE_PASSWORD:                             "Password change",
		STR_CHANGE_USERNAME:                             "Username change",
		STR_CLONE:                                       "Clone",
		STR_CONFIRM:                                     "Confirm",
		STR_CURRENT_SEARCH:                              "Current search",
		STR_DELETE:                                      "Delete",
		STR_DELETE_CONFIRM_EMAIL:                        "to permanently delete your account,",
		STR_DELETE_MENU:                                 "Delete menu",
		STR_DELETE_MENU_TEXT:                            "Are you sure you want to delete this menu?",
		STR_DELETE_RECIPE:                               "Delete recipe",
		STR_DELETE_RECIPE_TEXT:                          "Are you sure you want to delete this recipe?",
		STR_DELETE_SECTION_TEXT:                         "Are you sure? All the articles inside this section will be deleted.",
		STR_DELETE_SELECTED:                             "Delete selected",
		STR_DELETE_USER:                                 "Delete account",
		STR_DELETE_USER_TEXT1:                           "Are you sure to delete your account?",
		STR_DELETE_USER_TEXT2:                           "Are you REALLY sure to delete your account? This action is irreversible.",
		STR_DIRECTIONS:                                  "Directions",
		STR_EDIT:                                        "Edit",
		STR_EDIT_ARTICLE:                                "Edit article",
		STR_EDIT_ENTRY:                                  "Edit entry",
		STR_EDIT_MENU:                                   "Edit menu",
		STR_EDIT_RECIPE:                                 "Edit recipe",
		STR_EDIT_SECTION:                                "Edit section",
		STR_EMAIL:                                       "Email",
		STR_EMAIL_CHANGED:                               "Email changed succesfully",
		STR_EMAIL_SENT:                                  "We've sent you an email: please, check your inbox",
		STR_EXPIRATION:                                  "Expiration date",
		STR_FORGOT_PASSWORD:                             "Forgot password",
		STR_FRIDAY:                                      "Friday",
		STR_GOOD_MORNING:                                "Good morning " + placeholder + ",",
		STR_GOODBYE:                                     "Goodbye",
		STR_GOODBYE_EMAIL:                               "your account has been permanently deleted.",
		STR_INFO:                                        "Further informations",
		STR_INGREDIENTS:                                 "Ingredients",
		STR_LANG_CHANGED:                                "Language switched successfully",
		STR_LANGUAGE:                                    "Language",
		STR_LOGOUT:                                      "Logout",
		STR_MENU_DELETED:                                "Menu deleted succesfully",
		STR_MENUS:                                       "Menus",
		STR_MONDAY:                                      "Monday",
		STR_NAME:                                        "Name",
		STR_NEW_EMAIL:                                   "New email",
		STR_NEW_MENU:                                    "New menu",
		STR_NEW_PASSWORD:                                "New password",
		STR_NEW_RECIPE:                                  "New recipe",
		STR_NEW_SECTION:                                 "New section",
		STR_NEW_USERNAME:                                "New username",
		STR_NOREPLY:                                     "This email is automatically generated. Please do not reply.",
		STR_NOTES:                                       "Notes",
		STR_OK:                                          "Ok",
		STR_OLD_PASSWORD:                                "Old password",
		STR_ORDER_CHANGED:                               "The order of the articles has changed",
		STR_PAGE_NOT_FOUND:                              "Page not found",
		STR_PASSWORD:                                    "Password",
		STR_PASSWORD_CHANGED:                            "Password changed successfully",
		STR_PASSWORD_CHANGED_EMAIL:                      "recently your password has been changed.",
		STR_QUANTITY:                                    "Quantity",
		STR_RECIPE_DELETED:                              "Recipe deleted succesfully",
		STR_RECIPES:                                     "Recipes",
		STR_RECIPES_EMPTY:                               "No recipes found.",
		STR_REGARDS:                                     "Regards",
		STR_REPEAT_PASSWORD:                             "Repeat password",
		STR_RESET_PASSWORD:                              "Reset password",
		STR_RESET_PASSWORD_EMAIL:                        "to reset your password,",
		STR_SATURDAY:                                    "Saturday",
		STR_SAVE:                                        "Save",
		STR_SEARCH_ARTICLES:                             "Search articles",
		STR_SEARCH_EMPTY:                                "No articles found.",
		STR_SECTION:                                     "Section",
		STR_SECTION_DELETED:                             "Section deleted successfully",
		STR_SECTION_EMPTY:                               "This section is empty.",
		STR_SECTIONS:                                    "Storage sections",
		STR_SET_EMAIL_LANG:                              "Set emails language",
		STR_SETTINGS:                                    "Settings",
		STR_SHOPPINGLIST:                                "Shopping list",
		STR_SHOPPINGLIST_EMPTIED:                        "Entries deleted successfully",
		STR_SHOPPINGLIST_EMPTY:                          "The list is empty.",
		STR_SIGNIN:                                      "Sign in",
		STR_SIGNUP:                                      "Sign up",
		STR_SIGNUP_DONE:                                 "Succesfully signed up",
		STR_STARS:                                       "Stars",
		STR_STATS:                                       "Statistics",
		STR_STATS_ARTICLES:                              placeholder + " articles",
		STR_STATS_ENTRIES:                               placeholder + " entries",
		STR_STATS_MENUS:                                 placeholder + " menus",
		STR_STATS_RECIPES:                               placeholder + " recipes",
		STR_STATS_SECTIONS:                              placeholder + " sections",
		STR_STATS_USERS:                                 placeholder + " users",
		STR_STORAGE:                                     "Storage",
		STR_SUNDAY:                                      "Sunday",
		STR_THURSDAY:                                    "Thursday",
		STR_TUESDAY:                                     "Tuesday",
		STR_TUTORIAL:                                    "Tutorial",
		STR_UNKNOWN_LANG:                                "Unknown language",
		STR_UNKNOWN_REQUEST:                             "Unknown request",
		STR_UNMATCHING_PASSWORDS:                        "The two passwords do not match",
		STR_USER_CREATED:                                "Account created succesfully",
		STR_USER_DELETED:                                "Account deleted succesfully",
		STR_USERNAME:                                    "Username",
		STR_USERNAME_CHANGED:                            "Username changed succesfully",
		STR_WEDNESDAY:                                   "Wednesday",
		STR_WELCOME_EMAIL:                               "Welcome to CucinAssistant!",
		STR_WELCOMEBACK:                                 "Welcome back, " + placeholder + "!",
		STR_CLICK_HERE:                                  "click here",
	},
}
