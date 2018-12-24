// Provides basic bulding blocks for advanced console UI
//
// Coordinate system:
//
//  1/1---X---->
//   |
//   Y
//   |
//   v
//
// Documentation for ANSI codes: http://en.wikipedia.org/wiki/ANSI_escape_code#Colors
//
// Inspired by: http://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/
package goterm

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

// Reset all custom styles
const RESET = "\033[0m"

// Reset to default color
const RESET_COLOR = "\033[39;49m"

// Return cursor to start of line and clean it
const RESET_LINE = "\r\033[K"

// List of possible colors
const (
	BLACK = iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

var Output *bufio.Writer = bufio.NewWriter(os.Stdout)

// Get ANSI escape code for given color code for foreground
func GetColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}

// Get ANSI escape code for given color code for background
func GetBgColor(code int) string {
	return fmt.Sprintf("\033[4%dm", code)
}

// Get ANSI escape code for given RGB color for foreground
func GetRGBColor(r uint8, g uint8, b uint8) string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b)
}

// Get ANSI escape code for given RGB color for background
func GetBgRGBColor(r uint8, g uint8, b uint8) string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", r, g, b)
}

// Set percent flag: num | PCT
//
// Check percent flag: num & PCT
//
// Reset percent flag: num & 0xFF
const shift = uint(^uint(0)>>63) << 4
const PCT = 0x8000 << shift

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// Global screen buffer
// Its not recommended write to buffer dirrectly, use package Print,Printf,Println fucntions instead.
var Screen *bytes.Buffer = new(bytes.Buffer)

// Get relative or absolute coordinates
// To get relative, set PCT flag to number:
//
//      // Get 10% of total width to `x` and 20 to y
//      x, y = tm.GetXY(10|tm.PCT, 20)
//
func GetXY(x int, y int) (int, int) {
	if y == -1 {
		y = CurrentHeight() + 1
	}

	if x&PCT != 0 {
		x = int((x & 0xFF) * Width() / 100)
	}

	if y&PCT != 0 {
		y = int((y & 0xFF) * Height() / 100)
	}

	return x, y
}

type sf func(int, string) string

// Apply given transformation func for each line in string
func applyTransform(str string, transform sf) (out string) {
	out = ""

	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}

	return
}

// Clear screen
func Clear() {
	Output.WriteString("\033[2J")
}

// Move cursor to given position
func MoveCursor(x int, y int) {
	fmt.Fprintf(Screen, "\033[%d;%dH", y, x)
}

// Move cursor up relative the current position
func MoveCursorUp(bias int) {
	fmt.Fprintf(Screen, "\033[%dA", bias)
}

// Move cursor down relative the current position
func MoveCursorDown(bias int) {
	fmt.Fprintf(Screen, "\033[%dB", bias)
}

// Move cursor forward relative the current position
func MoveCursorForward(bias int) {
	fmt.Fprintf(Screen, "\033[%dC", bias)
}

// Move cursor backward relative the current position
func MoveCursorBackward(bias int) {
	fmt.Fprintf(Screen, "\033[%dD", bias)
}

// Move string to possition
func MoveTo(str string, x int, y int) (out string) {
	x, y = GetXY(x, y)

	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[%d;%dH%s", y+idx, x, line)
	})
}

// Return carrier to start of line
func ResetLine(str string) (out string) {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s", RESET_LINE, line)
	})
}

// Make bold
func Bold(str string) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[1m%s\033[0m", line)
	})
}

// Apply given color to string:
//
//     tm.Color("RED STRING", tm.RED)
//
func Color(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", GetColor(color), line, RESET)
	})
}

func ColorRGB(str string, r uint8, g uint8, b uint8) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", GetRGBColor(r, g, b), line, RESET)
	})
}

func Highlight(str, substr string, color int) string {
	hiSubstr := Color(substr, color)
	return strings.Replace(str, substr, hiSubstr, -1)
}

func HighlightRegion(str string, from, to, color int) string {
	return str[:from] + Color(str[from:to], color) + str[to:]
}

// Change background color of string:
//
//     tm.Background("string", tm.RED)
//
func Background(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", GetBgColor(color), line, RESET)
	})
}

func BackgroundRGB(str string, r uint8, g uint8, b uint8) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", GetBgRGBColor(r, g, b), line, RESET)
	})
}

// Get console width
func Width() int {
	ws, err := getWinsize()

	if err != nil {
		return -1
	}

	return int(ws.Col)
}

// Get console height
func Height() int {
	ws, err := getWinsize()
	if err != nil {
		return -1
	}
	return int(ws.Row)
}

// Get current height. Line count in Screen buffer.
func CurrentHeight() int {
	return strings.Count(Screen.String(), "\n")
}

// Flush buffer and ensure that it will not overflow screen
func Flush() {
	for idx, str := range strings.SplitAfter(Screen.String(), "\n") {
		if idx > Height() {
			return
		}

		Output.WriteString(str)
	}

	Output.Flush()
	Screen.Reset()
}

func Print(a ...interface{}) (n int, err error) {
	return fmt.Fprint(Screen, a...)
}

func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(Screen, a...)
}

func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(Screen, format, a...)
}

func Context(data string, idx, max int) string {
	var start, end int

	if len(data[:idx]) < (max / 2) {
		start = 0
	} else {
		start = idx - max/2
	}

	if len(data)-idx < (max / 2) {
		end = len(data) - 1
	} else {
		end = idx + max/2
	}

	return data[start:end]
}
