package fs

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestMD5Gen(t *testing.T) {
	byt, err := MD5Gen("/home/rustic/Downloads/hugo_0.48_Linux-64bit.tar.gz")
	assert.Equal(t, err, nil)
	assert.Equal(t, byt, true)
}

func TestCheckMD5(t *testing.T) {
	check := CheckMD5("/home/rustic/Downloads/hugo_0.48_Linux-64bit.tar.gz")
	assert.Equal(t, check, true)
}
