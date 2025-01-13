package pdfcsv

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/pdfcsv"
	"hopdf.com/helpers"
)

// This handler expects the body of th incoming request
// to have a pdf in it. The pdf will have a table that
// needs to be parsed.
func UploadHandler(c echo.Context) error {
	pdf, err := handler(c)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	fmt.Println("the pdf: ", pdf)

	csv := []string{"a", "b", "c"}

	csv_obj := &pdfcsv.PdfCsv{
		Csv: csv,
	}

	return helpers.Success(c, csv_obj)
}

func handler(c echo.Context) (*UploadPdfBody, error) {
	var reqBody UploadPdfBody

	if err := c.Bind(&reqBody); err != nil {
		return nil, err
	}

	// Process the data
	return &reqBody, nil
}
