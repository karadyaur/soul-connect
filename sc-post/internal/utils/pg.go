package utils

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func UUIDFromString(value string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid uuid: %w", err)
	}
	var bytes [16]byte
	copy(bytes[:], parsed[:])
	return pgtype.UUID{Bytes: bytes, Valid: true}, nil
}

func UUIDToString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	return uuid.UUID(value.Bytes).String()
}

func TextFromString(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: value, Valid: true}
}

func StringFromText(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

func TimestampToTime(value pgtype.Timestamp) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}

func NullableTextFromPointer(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}
	if *value == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}
