package peanut

import (
	"log"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

type excelBuilder struct {
	xlsx     *excelize.File
	sw       *excelize.StreamWriter
	row      int // TODO Expose this? we can report number of rows written (to be wary of Excel's row-limit)
	filename string
}

func newExcelBuilder(filename string) (*excelBuilder, error) {
	xlsx := excelize.NewFile()
	xlsx.SetPanes("Sheet1", `{"freeze":true,"split":false,"x_split":0,"y_split":1,"top_left_cell":"A2","active_pane":"bottomLeft","panes":[{"sqref":"A2","active_cell":"A2","pane":"bottomLeft"}]}`)
	sw, err := xlsx.NewStreamWriter("Sheet1")
	if err != nil {
		return nil, err
	}
	e := excelBuilder{
		xlsx:     xlsx,
		sw:       sw,
		row:      1,
		filename: filename,
	}
	return &e, nil
}

func (e *excelBuilder) AddRow(data ...interface{}) error {
	c, err := excelize.CoordinatesToCellName(1, e.row)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	err = e.sw.SetRow(c, data)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	e.row++
	return nil
}

func (e *excelBuilder) Save() error {
	err := e.sw.Flush()
	if err != nil {
		return err
	}
	return e.xlsx.SaveAs(e.filename)
}
