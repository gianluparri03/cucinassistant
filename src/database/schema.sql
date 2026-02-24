CREATE TABLE ca_version (id INT NOT NULL);

CREATE TABLE ca_users (
    uid SERIAL NOT NULL,

    username VARCHAR(250) NOT NULL,
    password VARCHAR(250) NOT NULL,
    email VARCHAR(250) NOT NULL,
    token VARCHAR(250),

    email_lang CHAR(2),
	newsletter CHAR(16),

    PRIMARY KEY (uid),
    UNIQUE (username),
    UNIQUE (email),
    UNIQUE (newsletter)
);


CREATE TABLE menus (
    uid INT NOT NULL,
    mid SERIAL NOT NULL,

    name VARCHAR(64) NOT NULL,

    PRIMARY KEY (mid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE
);

CREATE INDEX menus_uid_mid ON menus (uid, mid);

CREATE TABLE days (
    mid INT NOT NULL,
    position INT NOT NULL,

    name VARCHAR(64) NOT NULL,
    meals VARCHAR(512)[],

    PRIMARY KEY (mid, position),
    FOREIGN KEY (mid) REFERENCES menus (mid) ON DELETE CASCADE
);


CREATE TABLE sections (
    uid INT NOT NULL,
    sid SERIAL NOT NULL,

    name VARCHAR(128) NOT NULL,

    PRIMARY KEY (sid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX sections_uid ON sections (uid);

CREATE TABLE articles (
    sid INT NOT NULL,
    aid SERIAL NOT NULL,

    name VARCHAR(250) NOT NULL,
    quantity FLOAT,
    expiration DATE NOT NULL,

    PRIMARY KEY (aid),
    FOREIGN KEY (sid) REFERENCES sections (sid) ON DELETE CASCADE,
    UNIQUE (sid, name, expiration)
);

CREATE INDEX articles_sid_expiration ON articles (sid, expiration, aid);


CREATE TABLE entries (
    uid INT NOT NULL,
    eid SERIAL NOT NULL,

    name VARCHAR(250) NOT NULL,
    marked BOOLEAN DEFAULT FALSE,

    PRIMARY KEY (eid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name)
);

CREATE INDEX entries_uid_name ON entries (uid, name);


CREATE TABLE recipes (
    uid INT NOT NULL,
    rid SERIAL NOT NULL,

    name VARCHAR(64) NOT NULL,
    stars INT NOT NULL DEFAULT 0,

    ingredients VARCHAR(4096) NOT NULL DEFAULT '',
    directions VARCHAR(4096) NOT NULL DEFAULT '',
    notes VARCHAR(4096) NOT NULL DEFAULT '',

	code CHAR(8),

    PRIMARY KEY (rid),
    FOREIGN KEY (uid) REFERENCES ca_users (uid) ON DELETE CASCADE,
    UNIQUE (uid, name),
	UNIQUE (code)
);

CREATE INDEX recipes_uid_name ON recipes (uid, name);
CREATE INDEX recipes_code ON recipes (code);

CREATE TABLE tags (
    name VARCHAR NOT NULL, 
    rid INT NOT NULL,

    PRIMARY KEY (name, rid), 
    FOREIGN KEY (rid) REFERENCES recipes (rid) ON DELETE CASCADE
);

CREATE INDEX tags_name ON tags (name);
