package extdrm

import (
	"bytes"
)

type SaltGenerator int

const (
	SV6CloudGenerator SaltGenerator = iota + 1
)

func (sg SaltGenerator) KeySalt(salt []byte, path string) []byte {
	if sg != SV6CloudGenerator || len(salt) < 27 {
		return salt
	}

	salt = bytes.Clone(salt)
	size := len(path)
	salt[24] = path[size/5] ^ 0x2F
	salt[25] = path[size/3] ^ 0x5E
	salt[26] = path[size/2] ^ 0x1A
	return salt
}

func (sg SaltGenerator) PathSalt(salt []byte, path string) []byte {
	if sg != SV6CloudGenerator || len(salt) < 4 {
		return salt
	}

	salt = bytes.Clone(salt)
	size := len(path)
	salt[2] = uint8(size)
	salt[3] += uint8(size) / 2
	return salt
}
