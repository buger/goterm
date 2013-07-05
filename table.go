package terminal

import (
	"bytes"
	"text/tabwriter"
)

type Table struct {
	tabwriter.Writer

	Buf *bytes.Buffer
}

// Check http://golang.org/pkg/text/tabwriter/#Writer.Init
func NewTable(minwidth, tabwidth, padding int, padchar byte, flags uint) *Table {
	tbl := new(Table)
	tbl.Buf = new(bytes.Buffer)
	tbl.Init(tbl.Buf, minwidth, tabwidth, padding, padchar, flags)

	return tbl
}

func (t *Table) String() string {
	t.Flush()
	return t.Buf.String()
}
