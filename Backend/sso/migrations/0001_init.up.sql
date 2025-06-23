create table if not exists users (
                                     id serial primary key,
                                     email varchar not null unique,
                                     pass_hash bytea not null
);

create index if not exists idx_email on users (email);

create table if not exists apps (
                                    id serial primary key,
                                    name varchar not null unique,
                                    secret text not null unique
);
