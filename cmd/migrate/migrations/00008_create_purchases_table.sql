-- +goose Up
CREATE TABLE purchases (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tier_id     UUID NOT NULL REFERENCES ticket_tiers(id) ON DELETE CASCADE,
  quantity    INTEGER NOT NULL,
  total       INTEGER NOT NULL, -- price * quantity at time of purchase, in snallest denomination
  status      TEXT NOT NULL DEFAULT 'confirmed',  -- 'confirmed', 'cancelled', 'refunded'
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT  positive_quantity CHECK (quantity > 0),
  CONSTRAINT  positive_total CHECK (total >= 0)
);

CREATE INDEX idx_purchases_user_id ON purchases(user_id);
CREATE INDEX idx_purchases_tier_id ON purchases(tier_id);

DROP TRIGGER IF EXISTS update_purchases_updated_at ON purchases;
CREATE TRIGGER update_purchases_updated_at
  BEFORE UPDATE ON purchases
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- +goose Down
DROP TABLE IF EXISTS purchases CASCADE;
DROP TRIGGER IF EXISTS update_purchases_updated_at ON purchases;
