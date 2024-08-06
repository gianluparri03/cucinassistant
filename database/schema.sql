CREATE TABLE IF NOT EXISTS users (
    uid INT NOT NULL,

    username VARCHAR(250) NOT NULL,
    password VARCHAR(250) NOT NULL,
    email VARCHAR(250) NOT NULL,
    token VARCHAR(250),

    PRIMARY KEY (uid),
    UNIQUE (username),
    UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS menus (
    uid INT NOT NULL,
    mid INT NOT NULL,

    meal_1  VARCHAR(256),
    meal_2  VARCHAR(256),
    meal_3  VARCHAR(256),
    meal_4  VARCHAR(256),
    meal_5  VARCHAR(256),
    meal_6  VARCHAR(256),
    meal_7  VARCHAR(256),
    meal_8  VARCHAR(256),
    meal_9  VARCHAR(256),
    meal_10 VARCHAR(256),
    meal_11 VARCHAR(256),
    meal_12 VARCHAR(256),
    meal_13 VARCHAR(256),
    meal_14 VARCHAR(256),

    PRIMARY KEY (uid, mid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS storage_sections (
    uid INT NOT NULL,
    sid INT NOT NULL,

    name VARCHAR(128) NOT NULL,

    PRIMARY KEY (uid, sid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS storage_articles (
    uid INT NOT NULL,
    sid INT NOT NULL,
    aid INT NOT NULL,

    name VARCHAR(250) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    expiration DATE NOT NULL DEFAULT '2004-02-05',

    PRIMARY KEY (sid, aid),
    FOREIGN KEY (uid, sid) REFERENCES storage_sections (uid, sid) ON DELETE CASCADE,
    UNIQUE (sid, name, expiration)
);

CREATE TABLE IF NOT EXISTS shopping_entries (
    uid INT NOT NULL,
    eid INT NOT NULL,

    name VARCHAR(250) NOT NULL,
    marked BOOLEAN DEFAULT FALSE,

    PRIMARY KEY (uid, eid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);
