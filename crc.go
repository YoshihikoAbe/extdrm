package extdrm

import (
	"io"
)

const crcBlockSize = 262140

type crcSkipReader struct {
	rd io.Reader

	buf      [crcBlockSize]byte
	len, pos int
}

func (csr *crcSkipReader) Read(p []byte) (n int, err error) {
	for n < len(p) {
		if csr.pos >= csr.len {
			if err = csr.refill(); err != nil {
				return
			}
		}

		cnt := csr.len - csr.pos
		remain := len(p) - n
		if cnt > remain {
			cnt = remain
		}
		copy(p[n:n+cnt], csr.buf[csr.pos:csr.pos+cnt])

		csr.pos += cnt
		n += cnt
	}
	return
}

func (csr *crcSkipReader) refill() error {
	// skip crc
	if _, err := io.ReadFull(csr.rd, csr.buf[:4]); err != nil {
		return err
	}

	n, err := csr.rd.Read(csr.buf[:])
	if err != nil {
		return err
	}

	csr.pos = 0
	csr.len = n

	return nil
}
