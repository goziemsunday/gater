CREATE TABLE verifications (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  identifier  TEXT NOT NULL,
  -- e.g. "email-verification:user@example.com", "password-reset:user@example.com"
  value       TEXT NOT NULL, -- hashed token
  expires_at  TIMESTAMPTZ NOT NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_verifications_identifier ON verifications(identifier);
CREATE INDEX idx_verifications_value ON verifications(value);

CREATE TRIGGER update_verifications_updated_at
  BEFORE UPDATE ON verifications
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();