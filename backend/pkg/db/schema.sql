CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS UserAccount (
    UserID uuid PRIMARY KEY,
    Email varchar(256) NOT NULL UNIQUE,
    Password varchar(64) NOT NULL,
    CreationDate timestamptz NOT NULL
);
CREATE TABLE IF NOT EXISTS QA (
    StoryID uuid NOT NULL,
    QAID uuid NOT NULL,
    Question text NOT NULL,
    Answer text NOT NULL,
    Extension varchar(16),
    Attachment bytea,
    CreationDate timestamptz NOT NULL,
    PRIMARY KEY (StoryID, QAID)
);
CREATE TABLE IF NOT EXISTS Story (
    StoryID uuid PRIMARY KEY,
    UserID uuid NOT NULL references UserAccount(UserID),
    CreationDate timestamptz NOT NULL,
    Context int[]
);
CREATE TABLE IF NOT EXISTS Config (
    ConfigKey text PRIMARY KEY,
    ConfigValue text
);
