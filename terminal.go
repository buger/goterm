// Inspired by
// http://en.wikipedia.org/wiki/ANSI_escape_code#Colors
// http://www.darkcoding.net/software/pretty-command-line-console-output-on-unix-in-python-and-go-lang/
package terminal

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

const RESET = "\033[0m"
const RESET_COLOR = "\033[32m"

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

// Set percent flag: num | PCT
// Check percent flag: num & PCT
// Reset percent flag: num & 0xFF
const PCT = 0x80000000

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func getColor(code int) string {
	return fmt.Sprintf("\033[3%dm", code)
}

func getBgColor(code int) string {
	return fmt.Sprintf("\033[4%dm", code)
}

func getWinsize() (*winsize, error) {
	ws := new(winsize)

	var _TIOCGWINSZ int64

	switch runtime.GOOS {
	case "linux":
		_TIOCGWINSZ = 0x5413
	case "darwin":
		_TIOCGWINSZ = 1074295912
	}

	r1, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(_TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)),
	)

	if int(r1) == -1 {
		return nil, os.NewSyscallError("GetWinsize", errno)
	}
	return ws, nil
}

var Screen *bytes.Buffer = new(bytes.Buffer)

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

func applyTransform(str string, transform sf) (out string) {
	out = ""

	for idx, line := range strings.Split(str, "\n") {
		out += transform(idx, line)
	}

	return
}

func Clear() {
	fmt.Print("\033[2J")
}

func MoveCursor(x int, y int) {
	Printf("\033[%d;%dH", x, y)
}

func MoveTo(str string, x int, y int) (out string) {
	x, y = GetXY(x, y)

	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[%d;%dH%s", y+idx, x, line)
	})
}

func Bold(str string) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("\033[1m%s\033[0m", line)
	})
}

func Color(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getColor(color), line, RESET)
	})
}

func Background(str string, color int) string {
	return applyTransform(str, func(idx int, line string) string {
		return fmt.Sprintf("%s%s%s", getBgColor(color), line, RESET)
	})
}

func Width() int {
	ws, _ := getWinsize()
	return int(ws.Col)
}

func Height() int {
	ws, _ := getWinsize()
	return int(ws.Row)
}

func CurrentHeight() int {
	return strings.Count(Screen.String(), "\n")
}

func Flush() {
	for idx, str := range strings.Split(Screen.String(), "\n") {
		if idx > Height() {
			return
		}

		fmt.Println(str)
	}

	Screen.Reset()
}

func Print(a ...interface{}) {
	fmt.Fprint(Screen, a...)
}

func Println(a ...interface{}) {
	fmt.Fprintln(Screen, a...)
}

func Printf(format string, a ...interface{}) {
	fmt.Fprintf(Screen, format, a...)
}
