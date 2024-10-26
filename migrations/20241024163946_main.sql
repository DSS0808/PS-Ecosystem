-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table if not exists users
(
    id         uuid default uuid_generate_v4() primary key,
    username varchar(100),
    coins numeric(10, 5)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
drop extension "uuid-ossp";
-- +goose StatementEnd
