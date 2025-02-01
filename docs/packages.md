# Packages

CucinAssistant is written in Go, and it is divided in some packages: here I will describe briefly all of them.


## cucinassistant (main)

It just contains a `main.go` file that runs everything.
You can run everything just by using `make run` (this will create and start a postgresql container with docker),
or by hand with `go run .` inside the `src/` directory; in this case, make sure to also set the `CA_ENV` environment variable.

## cucinassistant/configs

This package contains some config files (for debugging, testing and for the ci), and the Go files that parses them.
You can look at `configs.go` to see what the fields mean.
It also contains `version.go`, that contains the current version of CucinAssistant.

## cucinassistant/database

The `database` package contains the lowest layer of CucinAssistant, and is the only one that can use the database
directly.
It exports the functions `Connect` and `Bootstrap` to set up a connection to the database and all the structs and
functions used from the website.  

This package is the only one with tests. They can be run with `make test`.

## cucinassistant/email

This package exports only a function, `SendMail`, that sends an email to one (or more) recipients.
There is also a `templates/` folder, that contains all the possible email templates.

## cucinassistant/web

It registers all the routes to a mux router, that can be started with the `Start` function.

## cucinassistant/web/utils

This package contains a couple utility functions used in the frontend part. I'll explain each file separately:

- `context.go` defines `Context` (a simple container of things used in all the website handlers). With that, it defines
the type `Handler` and `PHandler` (protected), that are all the functions that can serve an HTTP page. The latter can be
accessed only by registere users.

- `endpoint.go` defines an `Endpoint`, which is a path with optional Get and Post handlers.

- `renderer.go` contains `RenderPage`, `Show`, `Redirect` and `ShowAndRedirect`

- `sessions.go` adds the two functions `SaveUID` and `DropUID`, used to update the user's session

## cucinassistant/web/handlers

This package contains all the handlers, used in `cucinassistant/web/endpoints.go`.
Thanks to the `utils` module, the handling is simplified both for the authentication part (in fact
the user who requested the page is already an input), and in the error handling part (in fact the function
has an `error` return type; if not nil, it will be shown to the user).

## cucinassistant/langs

This package exports a function, `Translate`, used to translate all the UI, using `i18n`.
Inside the folder there are the `.toml` files for all the supported languages.

The IDs of the strings to be translated are always in uppercase and snake case; if they start with
`ERR_` it means they come from the `database` package, if they start with `MSG_` they come from the
`web/utils` or `web/handlers`, and if they start with `STR_` they come from the page templates.


## cucinassistant/tools

This folder contains some tools that can be used in pair with CucinAssistant.

- `broadcast.go` can be used in production to send an email to every user, like for scheduled manteinance.
It runs interactively, sending a test email to the sender email and - after a confirm - broadcast it to
each user. As for the `main.go` file, it needs the `CA_ENV` variable to be set.
