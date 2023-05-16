-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS GophKeeper(user_id text UNIQUE, login text UNIQUE, password text, aeskey text, time_stamp text, file text);
-- SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS GophKeeper;
-- SELECT 'down SQL query';
-- +goose StatementEnd
