package bencode

import (
	"bytes"
	"errors"
	"io"
)

type reader struct {
	*bytes.Reader
}

func newReader(data []byte) *reader {
	return &reader{bytes.NewReader(data)}
}

func (r *reader) readUntil(c byte) ([]byte, error) {
	res := []byte("")
	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return []byte{}, err
		}
		if b == c {
			break
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *reader) readNBytes(n uint64) ([]byte, error) {
	res := []byte("")
	var i uint64
	for i = 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return []byte(""), err
		}
		res = append(res, b)
	}
	return res, nil
}
