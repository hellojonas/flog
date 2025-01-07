CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    inactive BOOLEAN,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE applications (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    token TEXT NOT NULL,
    inactive BOOLEAN,
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users (id)
);

CREATE TABLE user_applications (
    user_id INTEGER NOT NULL,
    application_id INTEGER NOT NULL,
    PRIMARY KEY(user_id, application_id),
    FOREIGN KEY(user_id) REFERENCES users (id),
    FOREIGN KEY(application_id) REFERENCES applications (id)
);

CREATE TABLE logs (
    name TEXT NOT NULL,
    application_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(name, application_id),
    FOREIGN KEY(application_id) REFERENCES applications (id)
);
