CREATE TABLE IF NOT EXISTS ca_users (
    uid SERIAL NOT NULL,

    username VARCHAR(250) NOT NULL,
    password VARCHAR(250) NOT NULL,
    email VARCHAR(250) NOT NULL,
    token VARCHAR(250),

    email_lang CHAR(2) NOT NULL DEFAULT '',

    PRIMARY KEY (uid),
    UNIQUE (username),
    UNIQUE (email)
);


CREATE TABLE IF NOT EXISTS menus (
    uid INT NOT NULL,
    mid SERIAL NOT NULL,

    name VARCHAR(64) NOT NULL,
    meals VARCHAR(4096),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (mid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS menus_uid_createdat ON menus (uid, created_at);


CREATE TABLE IF NOT EXISTS sections (
    uid INT NOT NULL,
    sid SERIAL NOT NULL,

    name VARCHAR(128) NOT NULL,

    PRIMARY KEY (sid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX IF NOT EXISTS sections_uid ON sections (uid);

CREATE TABLE IF NOT EXISTS articles (
    sid INT NOT NULL,
    aid SERIAL NOT NULL,

    name VARCHAR(250) NOT NULL,
    quantity FLOAT,
    expiration DATE NOT NULL,

    PRIMARY KEY (aid),
    FOREIGN KEY (sid) REFERENCES sections (sid) ON DELETE CASCADE,
    UNIQUE (sid, name, expiration)
);

CREATE INDEX IF NOT EXISTS articles_sid_expiration ON articles (sid, expiration, aid);


CREATE TABLE IF NOT EXISTS entries (
    uid INT NOT NULL,
    eid SERIAL NOT NULL,

    name VARCHAR(250) NOT NULL,
    marked BOOLEAN DEFAULT FALSE,

    PRIMARY KEY (eid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX IF NOT EXISTS entries_uid ON entries (uid);


CREATE TABLE IF NOT EXISTS recipes (
    uid INT NOT NULL,
    rid SERIAL NOT NULL,

    name VARCHAR(64) NOT NULL,
    stars INT NOT NULL DEFAULT 0,

    ingredients VARCHAR(4096) NOT NULL DEFAULT '',
    directions VARCHAR(4096) NOT NULL DEFAULT '',
    notes VARCHAR(4096) NOT NULL DEFAULT '',

    PRIMARY KEY (rid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX IF NOT EXISTS recipies_uid_name ON recipes (uid, name);
