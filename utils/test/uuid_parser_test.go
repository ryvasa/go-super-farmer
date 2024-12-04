package test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/ryvasa/go-super-farmer/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	testUUID := uuid.New()
	parsedUUID, err := utils.ParseUUID(testUUID.String())
	assert.NoError(t, err)
	assert.Equal(t, testUUID, parsedUUID)

	_, err = utils.ParseUUID("invalid-uuid")
	assert.Error(t, err)
}
