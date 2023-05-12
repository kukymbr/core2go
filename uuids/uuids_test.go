package uuids_test

import (
	"testing"

	"github.com/kukymbr/core2go/uuids"
	"github.com/stretchr/testify/assert"
)

func TestNewString(t *testing.T) {
	count := 5
	generated := make(map[string]string, count)

	for i := 0; i < count; i++ {
		id := uuids.NewString()

		_, exists := generated[id]
		assert.False(t, exists)

		generated[id] = id
	}
}

func TestIsUUID_WhenValid_ExpectTrue(t *testing.T) {
	ids := [5]string{
		uuids.NewString(),
		uuids.NewString(),
		uuids.NewString(),
		uuids.NewString(),
		uuids.NewString(),
	}

	for _, id := range ids {
		assert.True(t, uuids.IsUUID(id))
	}
}

func TestIsUUID_WhenInvalid_ExpectFalse(t *testing.T) {
	ids := [5]string{
		"",
		"  ",
		"not_an_uuid_1",
		"not_an_uuid_2",
		"not-an-uuid-3333",
	}

	for _, id := range ids {
		assert.False(t, uuids.IsUUID(id))
	}
}
