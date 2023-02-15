-- +migrate Up

create table public.users
(
    id           serial,
    username     text,
    password     text,
    created_at   timestamp,
);

create unique index users_username_uindex
    on public.users (username);

-- +migrate Down

DROP TABLE users;
