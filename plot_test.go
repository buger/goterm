package goterm

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestCreateDataTable(t *testing.T) {
	data := new(DataTable)

	data.addColumn("Gender")
	data.addColumn("Age")

	if len(data.columns) != 2 {
		t.Error("Should be 2 columns")
	}

	if data.columns[1] != "Age" {
		t.Error("Should have proper column name")
	}

	data.addRow(1, 5)
	data.addRow(0, 4)

	if len(data.rows) != 2 {
		t.Error("Should have 2 rows")
	}

	if data.rows[1][0] != 0 && data.rows[1][1] != 4 {
		t.Error("Row should be properly inserted")
	}
}

func TestLineChart(t *testing.T) {
	chart := NewLineChart(100, 20)
	chart.Flags = DRAW_INDEPENDENT //| DRAW_RELATIVE

	data := new(DataTable)
	data.addColumn("Time")
	data.addColumn("Lat")
	data.addColumn("Count")

	//data.addColumn("x*x")

	for i := 0; i < 60; i++ {
		data.addRow(float64(i+60), float64(20+rand.Intn(10)), float64(i*2+rand.Intn(i+1))) // ,*/, x*x)
	}

	fmt.Println(chart.Draw(data, 0))
}
