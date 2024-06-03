CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS UserAccount (
    ID uuid PRIMARY KEY,
    Email varchar(256) NOT NULL UNIQUE,
    Password varchar(64) NOT NULL,
    CreationDate timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS ChatPair (
    ChatID uuid NOT NULL,
    PairID uuid NOT NULL,
    Request text NOT NULL,
    Response text NOT NULL,
    Extension varchar(16),
    Attachment bytea,
    PRIMARY KEY (ChatID, PairID)
);
CREATE TABLE IF NOT EXISTS Chat (
    ChatID uuid PRIMARY KEY,
    OwnerID uuid NOT NULL,
    CreationDate timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS Config (
    ConfigKey text PRIMARY KEY,
    ConfigValue text
);
