package internal

import (
	"bufio"
	"os"
)

type OutputStream struct {
	w       *bufio.Writer
	buf     byte
	bufSize int
}

func NewOutputStream(f *os.File) *OutputStream {
	return &OutputStream{w: bufio.NewWriter(f), buf: 0, bufSize: 0}
}

func (s *OutputStream) WriteBit(b uint32) {
	var bit byte = byte(b & 1)
	s.bufSize++
	s.buf = (s.buf << 1) | bit
	if s.bufSize == 8 {
		s.w.WriteByte(s.buf)
		s.w.Flush()
		s.buf = 0
		s.bufSize = 0
	}

}

func (s *OutputStream) FlushToByte(b uint32) {
	for s.bufSize != 0 {
		s.WriteBit(b)
	}
}
