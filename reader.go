package gomp4

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	// BoxHeaderSize is the size of a box header.
	BoxHeaderSize int64 = 8
)

// Box defines an Atom Box.
type Box struct {
	Type       string
	At         int64 // box header offset
	Size       int64 // box size
	DataOffset int64 // box data offset
	DataSize   int64 // data only size
}

func (b *Box) Repr() string {
	return fmt.Sprintf("box %s @%d: data ~ [%d (+%d), %d)", b.Type, b.At, b.DataOffset, b.DataSize, b.DataOffset+b.DataSize)
}

// MP4Reader defines an mp4 reader structure.
type MP4Reader struct {
	r    io.ReaderAt
	size int64

	Boxes []*Box
}

// ReadBoxData ...
func (m *MP4Reader) ReadBoxData(b *Box) (data []byte, err error) {
	return m.readNBytesAt(b.DataSize, b.DataOffset)
}

// readNBytesAt reads bytes from [offset, offset+n).
func (m *MP4Reader) readNBytesAt(n int64, offset int64) (data []byte, err error) {
	data = make([]byte, n)
	bytesRead, err := m.r.ReadAt(data, offset)
	if err != nil {
		if int64(bytesRead) == n && err == io.EOF {
			err = nil
		}
	}
	return
}

// readBox reads a box header from an offset.
/*
┌─────────────────────┐
|      Box Header     |
| Size (4) | Type (4) | Box Header = 8 Bytes
| --------------------|
|     Box Data (N)    | Box Data = N Bytes
└─────────────────────┘
           └─────────── Box Size = 8 + N bytes
*/
func (m *MP4Reader) readBox(offset int64) (b *Box, err error) {
	header, err := m.readNBytesAt(BoxHeaderSize, offset)
	if err != nil {
		err = fmt.Errorf("read box at offset %d: %w", offset, err)
		return
	}
	boxSize := int64(binary.BigEndian.Uint32(header[0:4]))
	b = &Box{
		Type:       string(header[4:BoxHeaderSize]),
		At:         offset,
		Size:       boxSize,
		DataOffset: offset + BoxHeaderSize,
		DataSize:   boxSize - BoxHeaderSize,
	}
	return
}

func (m *MP4Reader) parse() (err error) {
	var (
		pivot = int64(0)
		b     *Box
	)

	for {
		b, err = m.readBox(pivot)
		if err != nil {
			return
		}

		m.Boxes = append(m.Boxes, b)
		pivot += b.Size
		if pivot >= m.size {
			break
		}
	}

	return
}

// Open opens a MP4 file and the reader.
func Open(path string) (r *MP4Reader, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	r = &MP4Reader{
		r:     file,
		size:  info.Size(),
		Boxes: make([]*Box, 0),
	}
	err = r.parse()

	return
}
