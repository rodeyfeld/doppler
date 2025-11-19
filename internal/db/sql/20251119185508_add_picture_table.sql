-- +goose Up
-- +goose StatementBegin
CREATE TABLE picture (
	id integer primary key,
	post_id integer,
	filename text,
	modified datetime default current_timestamp,
	created datetime default current_timestamp,
	foreign key(post_id) references post(id) on delete cascade
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE picture;
-- +goose StatementEnd

