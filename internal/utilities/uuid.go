package utilities

import (
	"github.com/google/uuid"
)

// var (
// 	ErrInvalidUUID = errors.New("invalid uuid provided")
// 	// Define other errors here...
// )

func ParseUUID(id string) (uuid.UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	// if err != nil {
	// 	return parsedUUID, ErrInvalidUUID
	// }
	return parsedUUID, err
}
