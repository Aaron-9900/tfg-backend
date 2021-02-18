package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisInit(t *testing.T) {
	err := InitRedis()
	assert.NoError(t, err)
}

func TestRedisSetGetDelete(t *testing.T) {
	err := InitRedis()
	assert.NoError(t, err)

	err = RDB.SetUserAndToken(1, "test")
	assert.NoError(t, err)

	_, err = RDB.GetUserAndToken(1, "test")
	assert.NoError(t, err)

	count, err := RDB.DeleteToken(1, "test")
	assert.NoError(t, err)
	assert.Equal(t, count, int64(1))

	_, err = RDB.GetUserAndToken(1, "test")
	assert.Error(t, err)

}
