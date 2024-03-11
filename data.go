package gomp4

import (
	"encoding/binary"
	"fmt"
)

const (
	BoxFTyp = "ftyp"
	BoxMoov = "moov"
	BoxMvhd = "mvhd"
	BoxTrak = "trak"
	BoxTkhd = "tkhd"
	BoxMdia = "mdia"
	BoxMinf = "minf"
)

// FTyp is the FileType Data.
// Container: File
// Mandatory: Yes
// Quantity: Exactly one
type FTyp struct {
	MajorBrand       string   // Brand identifer.
	MinorVersion     uint32   // Informative integer for the minor version of the major brand.
	CompatibleBrands []string // A list, to the end of the box, of brands.
}

func (b *FTyp) Repr() string {
	return fmt.Sprintf(`
- FileType data
| major_brand: %s
| minor_version: %d
| compatible_brands: %v
`, b.MajorBrand, b.MinorVersion, b.CompatibleBrands)
}

// ParseFTyp parses the FileType box data.
func ParseFTyp(data []byte) (b *FTyp) {
	b = &FTyp{}
	b.MajorBrand = string(data[0:4])
	b.MinorVersion = binary.BigEndian.Uint32(data[4:8])
	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(data[i:i+4]))
		}
	}
	return
}

// Mvhd is the Movie Header Box data.
// Container: Movie Box (‘moov’)
// Mandatory: Yes
// Quantity: Exactly one
type Mvhd struct {
	Version          uint8 // 0 or 1
	Flags            []byte
	CreationTime     uint64 // 32 or 64 bits
	ModificationTime uint64
	Timescale        uint32
	Duration         uint64
	// others: TODO:
}

func (b *Mvhd) Repr() string {
	return fmt.Sprintf(`
- MovieHeader data
| version: %d, flags: %08b,
| create: %d, modify: %d
| timeScale: %d, duration %d
`, b.Version, b.Flags, b.CreationTime, b.ModificationTime, b.Timescale, b.Duration)
}

// ParseMvhd ...
func ParseMvhd(data []byte) (b *Mvhd) {
	b = &Mvhd{
		Version: uint8(data[0]),
		Flags:   data[1:4],
	}
	if b.Version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(data[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(data[8:12]))
		b.Timescale = binary.BigEndian.Uint32(data[12:16])
		b.Duration = uint64(binary.BigEndian.Uint32(data[16:20]))
		// last = 20
	} else {
		b.CreationTime = binary.BigEndian.Uint64(data[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(data[12:20])
		b.Timescale = binary.BigEndian.Uint32(data[20:24])
		b.Duration = binary.BigEndian.Uint64(data[24:32])
		// last = 32
	}
	return b
}

// Tkhd is the Track Header Box data.
// Container: Track Box (‘trak’)
// Mandatory: Yes
// Quantity: Exactly one
type Tkhd struct {
	Version          byte
	Flags            []byte
	CreationTime     uint64 // 32 or 64 bits
	ModificationTime uint64
	TrackID          uint32
	Reserved         uint32 // =0
	Duration         uint64
	Width, Height    uint32
}

func (b *Tkhd) Repr() string {
	return fmt.Sprintf(`
- TrackHeader data
| version: %d, flags: %08b,
| create: %d, modify: %d
| trackID: %d, duration %d, WxH: %dx%d
`, b.Version, b.Flags, b.CreationTime, b.ModificationTime, b.TrackID, b.Duration, b.Width, b.Height)
}

// ParseTkhd ...
func ParseTkhd(data []byte) (b *Tkhd) {
	b = &Tkhd{
		Version: uint8(data[0]),
		Flags:   data[1:4],
	}
	last := 0
	if b.Version == 0 {
		b.CreationTime = uint64(binary.BigEndian.Uint32(data[4:8]))
		b.ModificationTime = uint64(binary.BigEndian.Uint32(data[8:12]))
		b.TrackID = binary.BigEndian.Uint32(data[12:16])
		// reserved [16:20]
		b.Duration = uint64(binary.BigEndian.Uint32(data[20:24]))
		last = 24
	} else {
		b.CreationTime = binary.BigEndian.Uint64(data[4:12])
		b.ModificationTime = binary.BigEndian.Uint64(data[12:20])
		b.TrackID = binary.BigEndian.Uint32(data[20:24])
		// reserved [24:28]
		b.Duration = binary.BigEndian.Uint64(data[28:36])
		last = 36
	}
	/*
		const unsigned int(32)[2] reserved = 0;
		template int(16) layer = 0;
		template int(16) alternate_group = 0;
		template int(16) volume = {if track_is_audio 0x0100 else 0}; const unsigned int(16) reserved = 0;
		template int(32)[9] matrix=
		{ 0x00010000,0,0,0,0x00010000,0,0,0,0x40000000 };
	*/
	last += 4*2 + 2 + 2 + 2 + 4*9
	b.Width = binary.BigEndian.Uint32(data[last : last+4])
	b.Height = binary.BigEndian.Uint32(data[last+4 : last+8])
	return b
}
