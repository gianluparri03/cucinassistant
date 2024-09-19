# CucinAssistant

Warning: as of today, CucinAssistant's UI is in Italian, while the code (and the docs) are in English.
I plan to add the support for multiple languages in the next updates.

You can look at the [docs file](docs.md) to see how the code is organized.
Instead, the [running file](running.md) explains how to run CucinAssistant with docker compose.


## What is CucinAssistant

CucinAssistant is an utility website, with which you - and your roommates - can keep track of things
related to the kitchen. In particular, CucinAssistant features, a section for the menus management, a
shopping list (with a two-phase checking system) and a storage section, where you can save the items
stored in your fridge or cupboard (even in two separate folders), both with their quantiy and their expiration
date.  
As said before, the next steps will be a multi-language UI and a recipes section.


## Public instances

You can try it online, on [https://ca.gianlucaparri.me](https://ca.gianlucaparri.me), or you can also download it
from Google Play ([here](https://play.google.com/store/apps/details?id=me.gianlucaparri.ca.twa)); the app is a
Trusted Web Application, so it is the same of the website.


## A bit of history

The first version of CucinAssistant was written in Python (with Flask) and used MariaDB as its database; then,
after some experiments, I've decided to rewrite it completely, and now it is written in Go and uses PostgreSQL.

Note: the new versioning system consists of a progressive number and a codename; the old releases may still contain
the old version number, with a major and a minor number.

You can look at the [releases tab](https://github.com/gianluparri03/cucinassistant/releases/) for the complete list.


## Credits and license

The author of the project is Gianluca Parri: if you have any suggestion, found a bug or want to contribute, you can
send me an email at [gianlucaparri03@gmail.com](mailto:gianlucaparri03@gmail.com).

CucinAssistant is released with the MIT license.

In the website, CucinAssistant uses icons from [FontAwesome](https://fontawesome.com/),
the css from [sakura](https://github.com/oxalorg/sakura) and two fonts,
[Inclusive Sans](https://fonts.google.com/specimen/Inclusive+Sans?query=inclusive+sans) and
[Satisfy](https://fonts.google.com/specimen/Satisfy?query=satisfy).
