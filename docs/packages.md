# Packages

CucinAssistant is written in Go, and it is divided in some packages: here I will describe briefly all of them.


## cucinassistant (main)

It just contains a `main.go` file that runs everything.
You can run everything just by using `make run` or with `go run .`; in this
case, make sure to also set the `CA_ENV` environment variable.

## cucinassistant/configs

This package contains some config files (for debugging, testing and for the ci),
and the Go files that parses them.
You can look at `configs.go` to see what the fields mean.
It also contains `version.go`, that contains the current version.

## cucinassistant/database

The `database` package contains the lowest layer of CucinAssistant, and is the 
only one that can use the database directly.
It exports the functions `Connect` and `Bootstrap` to set up a connection to the
database and all the structs and functions used by the other packages.  

This package has automatic tests, that can be run with `make test`.

## cucinassistant/email

Contains all the functions necessary to send emails.

## cucinassistant/langs

Contains both the translations (`lang_en.go` and `lang_it.go`) and the
functions used to get the needed translations.

When running `make test`, it will also ensure that all the strings
(`strings.go`) have a translation in every language.

## cucinassistant/web

The web server.
Contains a file with all the endpoints (`endpoints.go`), written in their own
package.

## cucinassistant/web/utils

This package contains a couple utility functions used in the frontend part.
I'll explain each file separately:

- `context.go` defines `Context` (a simple container of things used in all the
website handlers). With that, it defines the type `Handler` and `PHandler`
(Protected, meaning that only registered users can access it), that are the
functions that can serve an HTTP page.

- `endpoint.go` defines an `Endpoint`, which is a path with optional Get and
  Post handlers.

- `renderer.go` contains `RenderComponent`, `ShowMessage`, `ShowError` and
  `Redirect`.

- `sessions.go` adds the two functions `SaveUID` and `DropUID`, used to update
  the user's session, and `SetLang`.

## cucinassistant/web/handlers

This package contains all the handlers.

## cucinassistant/tools

This folder contains some tools that can be used in pair with CucinAssistant.

- `broadcast.go` can be used in production to send an email to every user, like
  for scheduled manteinance or for a newsletter.
  It runs interactively, asking (for every language) a subject and a content.
  Then, after a confirm, it sends the email to users, in their language.
  As for the `main.go` file, it needs the `CA_ENV` variable to be set.
- `migrate.go` is used to upgrade the database schema to the new version from
  the previous one.
