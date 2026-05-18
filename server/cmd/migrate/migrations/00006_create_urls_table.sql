-- +goose Up
CREATE TABLE urls (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slug        TEXT NOT NULL UNIQUE,
  long_url    TEXT NOT NULL,
  expires_at  TIMESTAMPTZ, -- nullable, for expiring links
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_urls_slug ON urls(slug);
CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);

DROP TRIGGER IF EXISTS update_urls_updated_at ON urls;
CREATE TRIGGER update_urls_updated_at
  BEFORE UPDATE ON urls
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS urls CASCADE;
DROP TRIGGER IF EXISTS update_urls_updated_at ON urls;
