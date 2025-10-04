# Docs

## Chapters

- [packages.md](packages.md) contains explains how the code is organized.
- [docker.md](docker.md) explains how to run your own instance of CucinAssistant
  with docker compose.
- [history.md](history.md) contains a brief history of the project and a list
  of all the released versions.
- [icons.md](icons.md) contains more informations about the icons.

## Additional notes

## Terminology

- An `User` is a person (or a group of people) registered to CucinAssistant.
  It has a unique `UID`.
- A `Menu` is a collection of 14 meals. It has a unique `MID`.
- An `Article` is an item in storage, identified by a `AID`.
  A collection of `Article`s is called a `Section` (`SID`).
- An `Entry` is an item of an user's `ShoppingList`. It has a unique `EID`.
- A `Recipe` is identified by it's `RID`.
