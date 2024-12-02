-- +goose Up
-- +goose StatementBegin
CREATE TABLE post (
	id integer primary key,
	name text
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE post;
-- +goose StatementEnd

