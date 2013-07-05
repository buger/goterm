package terminal

import (
	"bytes"
	"strings"
)

const DEFAULT_BORDER = "- │ ┌ ┐ └ ┘"

type Box struct {
	Buf *bytes.Buffer

	Width  int
	Height int

	PaddingX int
	PaddingY int

	Border string
}

func NewBox(width, height int) *Box {
	width, height = GetXY(width, height)

	box := new(Box)
	box.Buf = new(bytes.Buffer)
	box.Width = width
	box.Height = height
	box.Border = DEFAULT_BORDER
	box.PaddingX = 1
	box.PaddingY = 0

	return box
}

func (b *Box) Write(p []byte) (int, error) {
	return b.Buf.Write(p)
}

func (b *Box) String() (out string) {
	borders := strings.Split(b.Border, " ")
	lines := strings.Split(b.Buf.String(), "\n")

	prefix := borders[1] + strings.Repeat(" ", b.PaddingX)
	suffix := strings.Repeat(" ", b.PaddingX) + borders[1]

	offset := b.PaddingY + 1
	contentWidth := b.Width - (b.PaddingX+1)*2

	for y := 0; y < b.Height; y++ {
		var line string

		switch {
		case y == 0:
			line = borders[2] + strings.Repeat(borders[0], b.Width-2) + borders[3]
		case y == (b.Height - 1):
			line = borders[4] + strings.Repeat(borders[0], b.Width-2) + borders[5]
		case y <= b.PaddingY || y >= (b.Height-b.PaddingY):
			line = borders[1] + strings.Repeat(" ", b.Width-2) + borders[1]
		default:
			if len(lines) > y-offset {
				line = lines[y-offset]
			} else {
				line = ""
			}

			if len(line) >= contentWidth-1 {
				line = line[0:contentWidth]
			} else {
				line = line + strings.Repeat(" ", contentWidth-len(line))
			}

			line = prefix + line + suffix
		}

		if y != b.Height-1 {
			line = line + "\n"
		}

		out += line
	}

	return out
}
