-- +goose Up
-- +goose StatementBegin
CREATE TABLE post (
	id integer primary key,
	user text,
	content text,
	created datetime default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post;
-- +goose StatementEnd

