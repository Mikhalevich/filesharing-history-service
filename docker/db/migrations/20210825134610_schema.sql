-- +goose Up
-- +goose StatementBegin
CREATE TABLE FileHistory(
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL,
    user_name varchar NOT NULL,
    file_name varchar NOT NULL,
    time integer NOT NULL,
    size integer NOT NULL,
    action integer NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP Table FileHistory;
-- +goose StatementEnd
