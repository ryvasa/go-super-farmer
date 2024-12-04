package utils

import "github.com/google/uuid"

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid, err
	}
	return uuid, nil
}
