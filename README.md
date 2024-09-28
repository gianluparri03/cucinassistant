# CucinAssistant

![Website](https://img.shields.io/website?url=https%3A%2F%2Fca.gianlucaparri.me)
![Codacy Grade](https://img.shields.io/codacy/grade/54e56adbe15f43568a1819224319b423)
![Codacy Coverage](https://img.shields.io/codacy/coverage/54e56adbe15f43568a1819224319b423)
![GitHub Actions](https://img.shields.io/github/actions/workflow/status/gianluparri03/cucinassistant/push.yml)
![Code Size](https://img.shields.io/github/languages/code-size/gianluparri03/cucinassistant)

> [!WARNING]
> As of today, CucinAssistant's UI is in Italian, while the code (and the docs) are in English. I plan to add the support for multiple languages in the next updates.


## What is CucinAssistant

CucinAssistant is an utility website, with which you - and your roommates - can keep track of things
related to the kitchen. In particular, CucinAssistant features, a section for the menus management, a
shopping list (with a two-phase checking system) and a storage section, where you can save the items
stored in your fridge or cupboard (even in two separate folders), both with their quantiy and their expiration
date.  
As said before, the next steps will be a multi-language UI and a recipes section.


## Docs

You can look at the [packages.md file](docs/packages.md) to see how the code is organized.
Instead, the [docker.md file](docs/docker.md) explains how to run CucinAssistant with docker compose.


## Public instances

You can try it online, on [https://ca.gianlucaparri.me](https://ca.gianlucaparri.me), or you can also download it
from Google Play ([here](https://play.google.com/store/apps/details?id=me.gianlucaparri.ca.twa)); the app is a
Trusted Web Application, so it is the same of the website.


## A bit of history

The first version of CucinAssistant was written in Python (with Flask) and used MariaDB as its database; then,
after some experiments, I've decided to rewrite it completely, and now it is written in Go and uses PostgreSQL and HTMX.

You can look at the [releases tab](https://github.com/gianluparri03/cucinassistant/releases/) for the complete version list.

> [!NOTE]
> The new versioning system consists of a progressive number and a codename; the old releases may still contain
the old version number, with a major and a minor number.


## Credits and license

The author of the project is Gianluca Parri: if you have any suggestion, found a bug or want to contribute, you can
send me an email at [gianlucaparri03@gmail.com](mailto:gianlucaparri03@gmail.com).

CucinAssistant is released with the MIT license.

In the website, CucinAssistant uses icons from [FontAwesome](https://fontawesome.com/),
the css from [sakura](https://github.com/oxalorg/sakura) and two fonts,
[Inclusive Sans](https://fonts.google.com/specimen/Inclusive+Sans?query=inclusive+sans) and
[Satisfy](https://fonts.google.com/specimen/Satisfy?query=satisfy).
