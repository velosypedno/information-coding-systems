package internal

import (
	"bufio"
	"io"
	"os"
)

type InputStream struct {
	r       *bufio.Reader
	done    bool
	lastBit uint32
	buf     byte
	bufSize int
}

func NewInputStream(f *os.File) *InputStream {
	return &InputStream{
		r:       bufio.NewReader(f),
		buf:     0,
		bufSize: 8,
	}
}

func (s *InputStream) ReadBit() (uint32, error) {
	if s.done {
		return s.lastBit, io.EOF
	}

	if s.bufSize == 8 {
		var err error
		s.buf, err = s.r.ReadByte()
		if err != nil {
			if err == io.EOF {
				s.done = true
				return s.lastBit, io.EOF
			}
			return 0, err
		}
		s.bufSize = 0
	}

	var b uint32 = uint32(s.buf&0x80) >> 7
	s.buf <<= 1
	s.bufSize++
	s.lastBit = b
	return b, nil
}
