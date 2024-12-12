-- +goose Up
-- +goose StatementBegin

CREATE TABLE user (
	id integer primary key,
	username text,
	password text,
	email text,
	modified datetime default current_timestamp,
	created datetime default current_timestamp
);

CREATE TABLE post (
	id integer primary key,
	user_id integer,
	title text,
	content text,
	modified datetime default current_timestamp,
	created datetime default current_timestamp,
	foreign key(user_id) references user(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user;
DROP TABLE post;

-- +goose StatementEnd

