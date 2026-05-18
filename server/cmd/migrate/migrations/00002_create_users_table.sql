-- +goose Up
CREATE TABLE users (
  id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name            TEXT NOT NULL,
  email           TEXT NOT NULL UNIQUE,
  password_hash   TEXT, -- NULL for OAuth-only accounts
  email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
  image           TEXT, -- profile picture URL from OAuth
  created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS users CASCADE;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
