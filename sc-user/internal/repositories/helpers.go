package repositories

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func stringToUUID(value string) (pgtype.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID %q: %w", value, err)
	}
	var bytes [16]byte
	copy(bytes[:], parsed[:])
	return pgtype.UUID{Bytes: bytes, Valid: true}, nil
}

func stringPtrToText(value *string) pgtype.Text {
	if value == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func textToStringPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}

func uuidToString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	parsed := uuid.UUID(value.Bytes)
	return parsed.String()
}

func uuidSliceToStrings(values []pgtype.UUID) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		str := uuidToString(v)
		if str != "" {
			result = append(result, str)
		}
	}
	return result
}
