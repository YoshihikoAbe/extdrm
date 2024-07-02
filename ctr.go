package extdrm

import (
	"crypto/cipher"
	"math/big"
)

const keyStreamSize = 16

var one = big.NewInt(1)

type ctr struct {
	block     cipher.Block
	cnt       *big.Int
	iv        *big.Int
	keyStream [keyStreamSize]byte
	offset    int
	annoying  [keyStreamSize * 2]byte
	magic     *big.Int
}

func newCTR(block cipher.Block, magic, iv *big.Int) cipher.Stream {
	return &ctr{
		block:  block,
		magic:  magic,
		iv:     iv,
		cnt:    big.NewInt(0),
		offset: keyStreamSize,
	}
}

func (c *ctr) refill() {
	i := big.NewInt(0).Mul(c.magic, c.cnt)
	c.cnt.Add(c.cnt, one)
	i.Add(i, c.iv)

	b := i.FillBytes(c.annoying[:])
	// output of Int.FillBytes is in big-endian, we
	// need it to be in little-endian
	reverseByteSlice(b)

	c.block.Encrypt(c.keyStream[:], b[:16])
	c.offset = 0
}

func (c *ctr) XORKeyStream(dst, src []byte) {
	for i := 0; i < len(src); i++ {
		if c.offset >= keyStreamSize {
			c.refill()
		}
		dst[i] = src[i] ^ c.keyStream[c.offset]
		c.offset++
	}
}
