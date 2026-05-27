-- +goose Up
CREATE TABLE tickets (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  purchase_id   UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
  tier_id       UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE CASCADE,
  qr_token      TEXT NOT NULL UNIQUE, -- HMAC-SHA256 signed token
  status        TEXT NOT NULL DEFAULT 'unused', -- 'unused', 'used', 'cancelled'
  checked_in_at TIMESTAMPTZ,  -- nullable, set on check-in
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tickets_purchase_id ON tickets(purchase_id);
CREATE INDEX idx_tickets_qr_token ON tickets(qr_token);

DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
CREATE TRIGGER update_tickets_updated_at
  BEFORE UPDATE ON tickets
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS tickets CASCADE;
DROP TRIGGER IF EXISTS update_tickets_updated_at ON tickets;
