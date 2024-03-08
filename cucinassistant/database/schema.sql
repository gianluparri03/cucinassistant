CREATE TABLE IF NOT EXISTS users (
    uid INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(250) NOT NULL UNIQUE,
    password VARCHAR(250) NOT NULL,
    email VARCHAR(250) NOT NULL UNIQUE,
    token VARCHAR(250)
);

CREATE TABLE IF NOT EXISTS menus (
    user INT REFERENCES users (uid) ON DELETE CASCADE,
    id INT NOT NULL,
    menu VARCHAR(4096) NOT NULL,
    prev INT,
    next INT,

    PRIMARY KEY (user, id),
    FOREIGN KEY (user, prev) REFERENCES menus (user, id),
    FOREIGN KEY (user, next) REFERENCES menus (user, id)
);

CREATE TABLE IF NOT EXISTS storage (
    user INT NOT NULL REFERENCES users (uid) ON DELETE CASCADE,
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(250) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    expiration DATE NOT NULL DEFAULT '2004-02-05',

    UNIQUE (user, name, expiration)
);

CREATE TABLE IF NOT EXISTS shopping (
    user INT NOT NULL REFERENCES users (uid) ON DELETE CASCADE,
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(250) NOT NULL,

    UNIQUE (user, name)
);

CREATE TABLE IF NOT EXISTS ideas (
    user INT NOT NULL REFERENCES users (uid) ON DELETE CASCADE,
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(250) NOT NULL,

    UNIQUE (user, name)
);
