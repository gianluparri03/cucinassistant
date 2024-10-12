# Docs

## Chapters

- [packages.md](packages.md) contains explains how the code is organized.
- [docker.md](docker.md) explains how to run your own instance of CucinAssistant with docker compose.
- [history.md](history.md) contains a brief history of the project and a list of all the released versions.
- [contributing.md](contributing.md) contains some informations about how to contribute to the project and a list
  of the next features to be implemented.

## Additional notes

## Terminology

- An `User` is a person (or a group of people) registered to cucinassistant. It has a unique `UID`.
- A `Menu` is a collection of 14 meals. It has a unique `MID`.
- An `Article` is an item in storage, identified by a `AID`. A collection of `Article`s is called a `Section` (`SID`).
- An `Entry` is an item of an user's `ShoppingList`. It has a unique `EID`.

## HTMX and HTML status codes

This project uses [HTMX](https://htmx.org/) to load data from the server and update the client's
page.

The first response from the server will contain the complete page. The next ones, that will contain
a specific header, will contain only the new data to be displayed.

For simplicity, a status code `200` means that the response contains some data to be shown in the page;
on the other hand, a `400` status code means that the response contains a message, that must be shown to
user, therefore it can be also an informative message, not necessarily an actual error.
