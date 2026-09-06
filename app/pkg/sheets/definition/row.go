package definition

import (
	"fmt"
	"strconv"
	"strings"

	"codeberg.org/tealeg/xlsx/v4"
)

type TypedRowKey interface {
	~int
	String() string
}

type TypedRow[T TypedRowKey] struct {
	*xlsx.Row
}

func (r TypedRow[T]) GetCell(idx T) *xlsx.Cell {
	return r.Row.GetCell(int(idx))
}

func (r TypedRow[T]) String() string {
	// headline
	list := []string{"#" + strconv.Itoa(r.GetCoordinate()+1)}
	maxWidth := 0
	keys := []string{}
	r.ForEachCell(func(cell *xlsx.Cell) error {
		_idx, _ := cell.GetCoordinates()
		idx := T(_idx)
		s := idx.String()
		l := len(s)
		if l > maxWidth {
			maxWidth = l
		}
		keys = append(keys, s)
		return nil
	})

	for idx, v := range keys {
		for len(v) < maxWidth {
			v = " " + v
		}
		keys[idx] = v
	}

	r.ForEachCell(func(cell *xlsx.Cell) error {
		idx, _ := cell.GetCoordinates()
		v := cell.String()
		// list = append(list, strconv.Itoa(idx)+": "+v)
		list = append(list, fmt.Sprintf("%s : %s", keys[idx], v))
		return nil
	})
	return strings.Join(list, "\n")
}
