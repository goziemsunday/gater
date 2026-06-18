package cursor

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Cursor struct {
	Timestamp time.Time
	ID        uuid.UUID
}

func Encode(c Cursor) string {
	raw := fmt.Sprintf("%s,%s", c.Timestamp.Format(time.RFC3339Nano), c.ID.String())
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func Decode(encoded string) (Cursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return Cursor{}, fmt.Errorf("cursor: decode cursor: %w", err)
	}

	parts := strings.SplitN(string(raw), ",", 2)
	if len(parts) != 2 {
		return Cursor{}, fmt.Errorf("cursor: decode cursor: invalid cursor")
	}

	ts, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return Cursor{}, fmt.Errorf("cursor: decode cursor: %w", err)
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return Cursor{}, fmt.Errorf("cursor: decode cursor: %w", err)
	}

	return Cursor{
		Timestamp: ts,
		ID:        id,
	}, nil
}
