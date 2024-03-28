CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    login TEXT NOT NULL UNIQUE,
    passHash BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS adverts (
    id INTEGER PRIMARY KEY,
    header TEXT NOT NULL,
    body TEXT,
    imageURL TEXT,
    price INTEGER NOT NULL,
    date DATETIME,
    authorLogin TEXT NOT NULL,
    FOREIGN KEY (authorLogin) REFERENCES users(login) ON DELETE CASCADE
);