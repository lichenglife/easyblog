package authn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {

	encrypt, err := Encrypt("123456")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypt)
	encrypt2, err := Encrypt("123456")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypt2)
	assert.NotEqual(t, encrypt, encrypt2)
}

func TestCompare(t *testing.T) {
	encrypt, err := Encrypt("123456")
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypt)
	err = Compare(encrypt, "123456")
	assert.NoError(t, err)
	
}
