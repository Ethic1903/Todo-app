create table if not exists url(
    id serial primary key,
    alias varchar not null unique,
    url varchar not null
);
create index if not exists idx_alias on url(alias);