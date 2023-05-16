-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS GophKeeperLocks(user_id text UNIQUE, sessionID text, time_lock text);
ALTER TABLE GophKeeper ADD COLUMN user_data bytea;
ALTER TABLE GophKeeper DROP COLUMN file; 
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
