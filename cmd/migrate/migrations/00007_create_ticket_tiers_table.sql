-- +goose Up
CREATE TABLE ticket_tiers (
  id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id  UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  name      TEXT NOT NULL,
  price     INTEGER NOT NULL,  -- smallest currency denomination
  quantity  INTEGER NOT NULL,  -- total inventory
  remaining INTEGER NOT NULL,  -- decremented on purchase
  status    TEXT NOT NULL DEFAULT 'available', -- 'available', 'sold_out'
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT positive_quantity CHECK (quantity > 0),
  CONSTRAINT valid_remaining CHECK (remaining >= 0 AND remaining <= quantity),
  CONSTRAINT positive_price CHECK (price >= 0)
);

CREATE INDEX idx_ticket_tiers_event_id ON ticket_tiers(event_id);

DROP TRIGGER IF EXISTS update_ticket_tiers_updated_at ON ticket_tiers;
CREATE TRIGGER update_ticket_tiers_updated_at
  BEFORE UPDATE ON ticket_tiers
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS ticket_tiers CASCADE;
DROP TRIGGER IF EXISTS update_ticket_tiers_updated_at ON ticket_tiers;
