package gomp4

import (
	"encoding/binary"
	"fmt"
)

const (
	BoxFTyp = "ftyp"
	BoxMoov = "moov"
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
