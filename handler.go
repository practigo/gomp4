package gomp4

import (
	"fmt"
	"strings"
)

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
		// data example
		if b.Type == BoxFTyp {
			ftypData, err := br.ReadData(b)
			if err != nil {
				return err
			}
			boxTreeView = append(boxTreeView, ParseFTyp(ftypData).Repr())
		}
		return nil
	}

	if err = boxes.Iter(box2str); err != nil {
		return err
	}

	fmt.Println(strings.Join(boxTreeView, ""))
	return nil
}
