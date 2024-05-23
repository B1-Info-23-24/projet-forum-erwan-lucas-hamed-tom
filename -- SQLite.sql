-- SQLite
CREATE TABLE USER (
    id INTEGER PRIMARY KEY,

    pseudo TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE ROOMS (
    id INTEGER PRIMARY KEY,
    created_by INTEGER NOT NULL,
    max_player INTEGER NOT NULL,
    name TEXT NOT NULL,
    id_game INTEGER,
    FOREIGN KEY (created_by) REFERENCES USER(id),
    FOREIGN KEY (id_game) REFERENCES GAMES(id)
);

CREATE TABLE ROOM_USERS (
    id_room INTEGER,
    id_user INTEGER,
    score INTEGER,
    FOREIGN KEY (id_room) REFERENCES ROOMS(id),
    FOREIGN KEY (id_user) REFERENCES USER(id),
    PRIMARY KEY (id_room, id_user)
);

CREATE TABLE GAMES (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
