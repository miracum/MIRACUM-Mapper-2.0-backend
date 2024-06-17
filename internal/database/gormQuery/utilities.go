package gormQuery

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/jackc/pgx/v5/pgconn"
)

func handlePgError(gormErr error) (*pgconn.PgError, bool) {
	err, ok := gormErr.(*pgconn.PgError)
	if !ok {
		// The error is not a *pgconn.PgError
		// Try to unwrap it and cast it again
		if unwrappedErr := errors.Unwrap(gormErr); unwrappedErr != nil {
			err, ok = unwrappedErr.(*pgconn.PgError)
			if !ok {
				// The unwrapped error is also not a *pgconn.PgError
				// Handle this case appropriately
				return nil, false
			}
		} else {
			// The error could not be unwrapped
			// Handle this case appropriately
			return nil, false
		}
	}

	// Now err is a *pgconn.PgError
	// You can use it as you wish
	// For example, return it directly
	return err, true
}

func extractIDFromErrorDetail(detail, fieldName string) (string, error) {
	pattern := fmt.Sprintf(`\((%s)\)=\((.*?)\)`, fieldName)
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(detail)
	if len(matches) > 2 {
		return matches[2], nil
	}
	return "", errors.New("unable to extract ID from error detail")
}
