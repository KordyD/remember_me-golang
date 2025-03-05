CREATE TABLE IF NOT EXISTS users
(
    id        uuid primary key,
    email     text unique not null,
    pass_hash text        not null
);

CREATE TABLE IF NOT EXISTS apps
(
    id     uuid primary key,
    name   text unique not null,
    secret text unique not null
)

-- TODO indexes