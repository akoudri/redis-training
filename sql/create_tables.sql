CREATE TABLE IF NOT EXISTS articles (
    id          SERIAL NOT NULL PRIMARY KEY,
    title       VARCHAR(50) NOT NULL UNIQUE,
    content     TEXT NOT NULL
);