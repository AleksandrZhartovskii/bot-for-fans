-- create schema telegram
create schema telegram;
alter schema telegram owner to postgres;

-- create table users
create table telegram.users
(
    id         bigserial not null
        constraint users_pk
            primary key,
    chat_id    bigint    not null,
    first_name varchar default ''::character varying,
    last_name  varchar default ''::character varying,
    user_name  varchar default ''::character varying
);
alter table telegram.users
    owner to postgres;
create unique index users_chat_id_uindex
    on telegram.users (chat_id);
create unique index users_id_uindex
    on telegram.users (id);

-- create table subscriptions
create table telegram.subscriptions
(
    id                bigserial not null
        constraint subscriptions_pk
            primary key,
    author_name       varchar   not null,
    last_update_time  timestamp default now(),
    exist             boolean   default false,
    author_spotify_id varchar   not null
);
alter table telegram.subscriptions
    owner to postgres;
create unique index subscriptions_id_uindex
    on telegram.subscriptions (id);
create unique index subscriptions_author_spotify_id_uindex
    on telegram.subscriptions (author_spotify_id);

-- create table user_subscriptions
create table telegram.user_subscriptions
(
    id              bigserial not null,
    user_id         bigserial not null,
    subscription_id bigserial not null
);
alter table telegram.user_subscriptions
    owner to postgres;
create unique index user_subscriptions_id_uindex
    on telegram.user_subscriptions (id);

commit;