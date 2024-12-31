package main

import (
	"errors"

	"github.com/google/uuid"
)

func parseStringToUUID(s string) (uuid.UUID, error) {
	var empty uuid.UUID

	if err := uuid.Validate(s); err != nil {
		return empty, errors.New("invalid UUID")
	}

	uuid, err := uuid.Parse(s)
	if err != nil {
		return empty, errors.New("cannot parse UUID")
	}

	return uuid, nil
}
