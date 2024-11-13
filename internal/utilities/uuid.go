package utilities

import (
	"github.com/google/uuid"
)

func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	return parsedUUID, err
}
