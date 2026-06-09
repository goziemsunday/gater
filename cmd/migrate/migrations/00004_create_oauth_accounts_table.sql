-- +goose Up
CREATE TABLE oauth_accounts (
  id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id                   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider                  TEXT NOT NULL,  -- 'google'
  provider_account_id       TEXT NOT NULL,  -- Google's user ID
  scope                     TEXT,
  created_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at                TIMESTAMPTZ NOT NULL DEFAULT now(),
  
  UNIQUE  (provider, provider_account_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_accounts_user_id ON oauth_accounts(user_id);

DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at ON oauth_accounts;
CREATE TRIGGER update_oauth_accounts_updated_at
  BEFORE UPDATE ON oauth_accounts
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS oauth_accounts CASCADE;
DROP TRIGGER IF EXISTS update_oauth_accounts_updated_at on oauth_accounts;
