-- +goose Up
CREATE TABLE clicks (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  url_id      UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
  clicked_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  referrer    TEXT,
  user_agent  TEXT,
  ip_hash     TEXT -- SHA-256 of IP, for uniqueness estimation without storing PII
);

CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id);
CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);

-- +goose Down
DROP TABLE IF EXISTS clicks CASCADE;
