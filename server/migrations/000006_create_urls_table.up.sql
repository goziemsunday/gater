CREATE TABLE urls (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  slug        TEXT NOT NULL UNIQUE,
  long_url    TEXT NOT NULL,
  expires_at  TIMESTAMPTZ, -- nullable, for expiring links
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_urls_slug ON urls(slug);
CREATE INDEX idx_urls_user_id ON urls(user_id);

CREATE TRIGGER update_urls_updated_at
  BEFORE UPDATE ON urls
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();