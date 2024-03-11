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

// Range is the byte-address range of the memory [Start, End).
type Range struct {
	Start int64
	End   int64 // not inclusive
}

func (r Range) Size() int64 {
	return r.End - r.Start
}

// Box defines an Atom/ISOBMFF object.
// https://en.wikipedia.org/wiki/ISO_base_media_file_format
/*
┌─────────────────────┐
|      Box Header     |
| Size (4) | Type (4) | Box Header = 8 Bytes
| --------------------|
|     Box Data (N)    | Box Data = N Bytes
└─────────────────────┘
           └─────────── Box Size = 8 + N bytes
*/
type Box struct {
	// from the header
	Type string
	Size int64 // box size including the header
	// from parsing
	At    int64 // box offset w.r.t. the file
	Depth uint  // optional info

	Children Boxes
}

// DataRange returns the byte-range storing the box data.
func (b *Box) DataRange() Range {
	return Range{
		Start: b.At + BoxHeaderSize,
		End:   b.At + b.Size,
	}
}

// HandleBox handles a single Box.
type HandleBox func(*Box) error

// Boxes is an alias of []*Box.
type Boxes []*Box

// Iter loops thru the box-tree in a depth-first manner.
func (bs Boxes) Iter(f func(*Box) error) error {
	for _, b := range bs {
		if err := f(b); err != nil {
			return err
		}
		if len(b.Children) > 0 {
			if err := b.Children.Iter(f); err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadNBytesAt reads n bytes from [offset, offset + n).
func ReadNBytesAt(r io.ReaderAt, n int64, offset int64) (data []byte, err error) {
	data = make([]byte, n)
	bytesRead, err := r.ReadAt(data, offset)
	if err != nil {
		if int64(bytesRead) == n && err == io.EOF {
			err = nil
		}
	}
	return
}

// ReadNBytesAt reads from range r.
func ReadNBytesFrom(reader io.ReaderAt, r Range) (data []byte, err error) {
	return ReadNBytesAt(reader, r.Size(), r.Start)
}

// ReadBox reads and parses the 8-bytes box header at offset.
func ReadBox(r io.ReaderAt, offset int64) (b *Box, err error) {
	header, err := ReadNBytesAt(r, BoxHeaderSize, offset)
	if err != nil {
		err = fmt.Errorf("read box at offset %d: %w", offset, err)
		return
	}
	boxSize := int64(binary.BigEndian.Uint32(header[0:4]))
	b = &Box{
		Type: string(header[4:BoxHeaderSize]),
		Size: boxSize,
		At:   offset,
	}
	return
}

// DefaultNested ...
func DefaultNested() map[string]bool {
	return map[string]bool{
		BoxMoov: true,
	}
}

func getBoxes(m io.ReaderAt, r Range, depth uint, nested map[string]bool) (bs Boxes, err error) {
	var (
		b   *Box
		cur = r.Start
	)

	for {
		b, err = ReadBox(m, cur)
		if err != nil {
			return
		}
		b.Depth = depth
		if isNested, ok := nested[b.Type]; ok && isNested {
			b.Children, err = getBoxes(m, b.DataRange(), depth+1, nested)
		}
		bs = append(bs, b)
		// update
		cur += b.Size
		if cur >= r.End {
			break
		}
	}

	return
}

// BoxReader reads the boxes from an ISOBMFF file.
type BoxReader interface {
	// GetBoxes return the box-tree recursively w.r.t. box-type map nested.
	GetBoxes() (bs Boxes, err error)

	// ReadData reads the data bytes pointed by the box b.
	ReadData(b *Box) (data []byte, err error)
}

type mp4Reader struct {
	r      io.ReaderAt
	s      int64
	nested map[string]bool
}

func (m *mp4Reader) GetBoxes() (bs Boxes, err error) {
	return getBoxes(m.r, Range{0, m.s}, 0, m.nested)
}

func (m *mp4Reader) ReadData(b *Box) (data []byte, err error) {
	return ReadNBytesFrom(m.r, b.DataRange())
}

// Open opens a MP4 file and return the reader.
func Open(path string) (r BoxReader, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &mp4Reader{
		r:      file,
		s:      info.Size(),
		nested: DefaultNested(),
	}, nil
}
