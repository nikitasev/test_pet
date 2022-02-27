create table "user"
(
    id         serial
        constraint clients_pkey
            primary key,
    name       varchar(50) not null,
    is_deleted smallint default 0
);