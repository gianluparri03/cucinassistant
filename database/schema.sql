CREATE TABLE IF NOT EXISTS users (
    uid INT NOT NULL AUTO_INCREMENT,

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
    mid INT NOT NULL AUTO_INCREMENT,

    name VARCHAR(64) NOT NULL,
    meals VARCHAR(4096),

    PRIMARY KEY (mid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS user ON menus (uid);


CREATE TABLE IF NOT EXISTS sections (
    uid INT NOT NULL,
    sid INT NOT NULL AUTO_INCREMENT,

    name VARCHAR(128) NOT NULL,

    PRIMARY KEY (sid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX IF NOT EXISTS user_name ON sections (uid, name);

CREATE TABLE IF NOT EXISTS articles (
    uid INT NOT NULL,
    sid INT NOT NULL,
    aid INT NOT NULL AUTO_INCREMENT,

    name VARCHAR(250) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    expiration DATE NOT NULL DEFAULT '2004-02-05',

    PRIMARY KEY (aid),
    FOREIGN KEY (sid) REFERENCES sections (sid) ON DELETE CASCADE,
    UNIQUE (sid, name, expiration)
);


CREATE TABLE IF NOT EXISTS entries (
    uid INT NOT NULL,
    eid INT NOT NULL AUTO_INCREMENT,

    name VARCHAR(250) NOT NULL,
    marked BOOLEAN DEFAULT FALSE,

    PRIMARY KEY (eid),
    FOREIGN KEY (uid) REFERENCES users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX IF NOT EXISTS user ON entries (uid);
