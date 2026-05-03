\c chat_service_db

CREATE TABLE IF NOT EXISTS permissions(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS roles(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    name TEXT UNIQUE NOT NULL
);


INSERT INTO permissions (id, name) VALUES ('1', 'send messages');
INSERT INTO permissions (id, name) VALUES ('2', 'add users to group');
insert into permissions (id, name) values ('3', 'delete messages');
INSERT INTO permissions (id, name) VALUES ('4', 'remove users frosend messagesm group');
INSERT INTO permissions (id, name) VALUES ('5', 'ban users');


INSERT INTO roles (id, name) VALUES ('1', 'admin');
INSERT INTO roles (id, name) VALUES ('2', 'member');
