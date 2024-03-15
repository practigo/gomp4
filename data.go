package gomp4

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	BoxFTyp = "ftyp"
	BoxMoov = "moov"
	BoxMvhd = "mvhd"
	BoxTrak = "trak"
	BoxTkhd = "tkhd"
	BoxMdia = "mdia"
	BoxMdhd = "mdhd"
	BoxMinf = "minf"
	BoxStbl = "stbl"
	BoxEdts = "edts"
)

type BoxData interface {
	// Repr returns a string representation of the box data.
	Repr(prefix string) string
}

func dataStr(lines []string, prefix string) string {
	new_lines := []string{prefix + " =="}
	new_lines = append(new_lines, lines...)
	new_lines = append(new_lines, "==\n")
	return strings.Join(new_lines, "\n"+prefix+" ")
}

// DataParser ...
type DataParser func(data []byte) BoxData

// FTyp is the FileType Data.
// Container: File
// Mandatory: Yes
// Quantity: Exactly one
type FTyp struct {
	MajorBrand       string   // Brand identifer.
	MinorVersion     uint32   // Informative integer for the minor version of the major brand.
	CompatibleBrands []string // A list, to the end of the box, of brands.
}

func (b *FTyp) Repr(prefix string) string {
	return dataStr([]string{
		"FileType data:",
		fmt.Sprintf("major_brand: %s, minor_version: %d", b.MajorBrand, b.MinorVersion),
		fmt.Sprintf("compatible_brands: %v", b.CompatibleBrands),
	}, prefix)
}

// ParseFTyp parses the FileType box data.
func ParseFtyp(data []byte) BoxData {
	b := &FTyp{}
	b.MajorBrand = string(data[0:4])
	b.MinorVersion = binary.BigEndian.Uint32(data[4:8])
	if len(data) > 8 {
		for i := 8; i < len(data); i += 4 {
			b.CompatibleBrands = append(b.CompatibleBrands, string(data[i:i+4]))
		}
	}
	return b
}

type vflags struct {
	Version uint8 // 0 or 1
	Flags   []byte
}

func parseVFlags(data []byte) vflags {
	return vflags{
		Version: uint8(data[0]),
		Flags:   data[1:4],
	}
}

type header struct {
	CreationTime     uint64 // 32 or 64 bits
	ModificationTime uint64
	Timescale        uint32
	Duration         uint64
}

func parseHd(data []byte, ver uint8) (header, int) {
	if ver == 0 {
		return header{
			CreationTime:     uint64(binary.BigEndian.Uint32(data[4:8])),
			ModificationTime: uint64(binary.BigEndian.Uint32(data[8:12])),
			Timescale:        binary.BigEndian.Uint32(data[12:16]),
			Duration:         uint64(binary.BigEndian.Uint32(data[16:20])),
		}, 20
	}
	return header{
		CreationTime:     uint64(binary.BigEndian.Uint32(data[4:12])),
		ModificationTime: uint64(binary.BigEndian.Uint32(data[12:20])),
		Timescale:        binary.BigEndian.Uint32(data[20:24]),
		Duration:         uint64(binary.BigEndian.Uint32(data[24:32])),
	}, 32
}

// Mvhd is the Movie Header Box data.
// Container: Movie Box (‘moov’)
// Mandatory: Yes
// Quantity: Exactly one
type Mvhd struct {
	vflags
	header
	// others: TODO:
}

func (b *Mvhd) Repr(prefix string) string {
	return dataStr([]string{
		"MovieHeader data:",
		fmt.Sprintf("version: %d, flags: %08b, create: %d, modify: %d", b.Version, b.Flags, b.CreationTime, b.ModificationTime),
		fmt.Sprintf("timeScale: %d, duration %d", b.Timescale, b.Duration),
	}, prefix)
}

// ParseMvhd ...
func ParseMvhd(data []byte) BoxData {
	vf := parseVFlags(data)
	hd, _ := parseHd(data, vf.Version)
	return &Mvhd{vf, hd}
}

// Tkhd is the Track Header Box data.
// Container: Track Box (‘trak’)
// Mandatory: Yes
// Quantity: Exactly one
type Tkhd struct {
	vflags
	CreationTime     uint64 // 32 or 64 bits
	ModificationTime uint64
	TrackID          uint32
	Reserved         uint32 // =0
	Duration         uint64
	Width, Height    uint32
}

func (b *Tkhd) Repr(prefix string) string {
	return dataStr([]string{
		"TrackHeader data:",
		fmt.Sprintf("version: %d, flags: %08b, create: %d, modify: %d", b.Version, b.Flags, b.CreationTime, b.ModificationTime),
		fmt.Sprintf("trackID: %d, duration %d, WxH: %dx%d", b.TrackID, b.Duration, b.Width, b.Height),
	}, prefix)
}

// ParseTkhd ...
func ParseTkhd(data []byte) BoxData {
	b := &Tkhd{
		vflags: parseVFlags(data),
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

// Mdhd is the Media Header Box data.
// Container: MediaBox(‘mdia’)
// Mandatory: Yes
// Quantity: Exactly one
type Mdhd struct {
	vflags
	header
	// others
	Language []rune // ISO-639-2/T language code
}

func (b *Mdhd) Repr(prefix string) string {
	return dataStr([]string{
		"MediaHeader data:",
		fmt.Sprintf("version: %d, flags: %08b, create: %d, modify: %d", b.Version, b.Flags, b.CreationTime, b.ModificationTime),
		fmt.Sprintf("timeScale: %d, duration %d", b.Timescale, b.Duration),
		fmt.Sprintf("language: %s (raw: %v)", string(b.Language), b.Language),
	}, prefix)
}

/*
decode language at a special 16bit/2bytes. The input data is offset to 0-start.

bit(1)   pad = 0;
unsigned int(5)[3]   language;

language declares the language code for this media. See ISO 639‐2/T for the set of three
character codes. Each character is packed as the difference between its ASCII value and 0x60.
Since the code is confined to being three lower‐case letters, these values are strictly positive.
*/
func decLang(data []byte) (codes []rune) {
	_ = data[1] // bounds check hint to compiler
	// [0]: 1(pad)+5(a)+2(b)
	// [1]: 3(b)+5(c)
	codes = append(codes,
		rune(uint8(data[0])>>2+0x60),
		rune(uint8(data[0]&3)<<3+uint8(data[1])>>5+0x60),
		rune(uint8(data[1]&31)+0x60))
	return codes
}

// ParseMdhd ...
func ParseMdhd(data []byte) BoxData {
	vf := parseVFlags(data)
	hd, lastIdx := parseHd(data, vf.Version)
	lang := decLang(data[lastIdx:])

	return &Mdhd{vf, hd, lang}
}
