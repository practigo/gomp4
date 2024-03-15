package gomp4

import (
	"fmt"
	"strings"
)

var dataHandler = map[string]DataParser{
	BoxFTyp: ParseFtyp,
	BoxMvhd: ParseMvhd,
	BoxTkhd: ParseTkhd,
	BoxMdhd: ParseMdhd,
}

// View prints the box-tree structure.
func View(br BoxReader) error {
	boxes, err := br.GetBoxes()
	if err != nil {
		return err
	}

	boxTreeView := make([]string, 0)

	box2str := func(b *Box) error {
		dr := b.DataRange()
		prefix := strings.Repeat("--|", int(b.Depth+1)) // depth starts from 0
		boxTreeView = append(boxTreeView, fmt.Sprintf("%s box %s @%d: data ~ [%d (+%d), %d)\n", prefix, b.Type, b.At, dr.Start, dr.Size(), dr.End))
		// data parse
		if dp, ok := dataHandler[b.Type]; ok {
			data, err := br.ReadData(b)
			if err != nil {
				return err
			}
			bd := dp(data)
			boxTreeView = append(boxTreeView, bd.Repr(prefix))
		}
		return nil
	}

	if err = boxes.Iter(box2str); err != nil {
		return err
	}

	fmt.Println(strings.Join(boxTreeView, ""))
	return nil
}
