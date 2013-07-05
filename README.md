## Description

This library provides basic building blocks for building advanced console UIs.

Initially created for [Gor](github.com/buger/gor).

Full API documentation: http://godoc.org/github.com/buger/goterm

## Basic usage

Full screen console app, priting current time:

```go
import (
    tm "github.com/buger/goterm"
    "time"
)

func main() {
    tm.Clear() // Clear current screen

    for {
        // By moving cursor to top-left position we ensure that console output will be overwritten each time, instead of adding new.
        tm.MoveCursor(1,1)

        tm.Println("Current Time:", time.Now().Format(time.RFC1123))

        tm.Flush() // Call it every time at the end of rendering

        time.Sleep(time.Second)
    }
}
```


Print red bold message on white background:

```go    
tm.Println(tm.Backgound(tm.Color(tm.Bold("Important header"), tm.RED), tm.WHITE))
```


Create box and move it to center of the screen:

```go
// Create Box with 30% width of current screen, and height of 20 lines
box := tm.NewBox(30|tm.PCT, 20)

fmt.Fprint(box, "Some box content")

// Move Box to approx center of the screen
tm.Print(tm.MoveTo(box.String(), 40|tm.PCT, 40|tm.PCT))
```


Draw table:

```go
// Based on http://golang.org/pkg/text/tabwriter
totals := tm.NewTable(0, 10, 5, ' ', 0)
fmt.Fprintf(totals, "Time\tStarted\tActive\tFinished\n")
fmt.Fprintf(totals, "%s\t%d\t%d\t%d\n", "All", started, started-finished, finished)
tm.Println(totals)
```