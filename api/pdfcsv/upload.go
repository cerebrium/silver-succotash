package pdfcsv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/pdfcsv"
	"hopdf.com/dao/stations"
	"hopdf.com/dao/tiers"
	"hopdf.com/dao/weights"
	"hopdf.com/helpers"
	"hopdf.com/localware"
)

// This handler expects the body of th incoming request
// to have a pdf in it. The pdf will have a table that
// needs to be parsed.
func UploadHandler(c echo.Context) error {
	csv, err := handler(c)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	csv_obj := &pdfcsv.PdfCsv{
		Csv: csv,
	}

	return helpers.Success(c, csv_obj)
}

func handler(c echo.Context) ([]string, error) {
	cc, ok := c.(*localware.LocalUserClerkDbContext)
	if !ok {
		c.Logger().Error("could not resolve cc")
	}
	// Retrieve the file from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file metadata
	// fileData := UploadPfdFile{
	// 	Name: fileHeader.Filename,
	// }
	//
	filePath := filepath.Join("./uploads", fileHeader.Filename)
	txt_file := strings.Replace(filePath, ".pdf", ".txt", -1)
	txt_file_destination := filepath.Join(txt_file)

	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	err = InternalConvertPdfToText(fileHeader.Filename)
	if err != nil {
		return nil, err
	}

	local_weights := weights.Weights{ID: 1}
	updated_weights, err := local_weights.Read(cc.Db)
	if err != nil {
		return nil, err
	}

	stations, err := stations.ReadAll(cc.Db)
	if err != nil {
		return nil, err
	}

	tierList, err := tiers.ReadTiers(cc.Db)
	if err != nil {
		return nil, err
	}

	percentMap := PercentMap{}
	percentMap["dcr_val"] = updated_weights.Dcr
	percentMap["dnr_dpmo_val"] = updated_weights.DnrDpmo
	percentMap["ce_val"] = updated_weights.Ce
	percentMap["pod_val"] = updated_weights.Pod
	percentMap["cc_val"] = updated_weights.Cc
	percentMap["dex_val"] = updated_weights.Dex
	percentMap["lor_val"] = updated_weights.Lor

	csv_list := []string{}

	csv_headers := "Transporter ID,Delivered,DCR,DNR DPMO,LoR DPMO,POD,CC,CE,DEX"
	csv_list = append(csv_list, csv_headers)

	stringified_pdf, err := os.Open(txt_file_destination)
	if err != nil {
		return nil, err
	}

	defer stringified_pdf.Close()

	scanner := bufio.NewScanner(stringified_pdf)
	should_compute := false
	testing := 0
	for scanner.Scan() {
		if testing > 1 {
			return []string{}, nil
		}
		line := scanner.Text()
		if strings.Contains(line, "DSP WEEKLY SUMMARY") {
			should_compute = true
			continue
		}

		if strings.Contains(line, "Drivers With Working Hour Exceptions") {
			should_compute = false
			break
		}

		if should_compute {
			if strings.Contains(line, "Transporter ID") || strings.Contains(line, "Page") || strings.TrimSpace(line) == "" {
				continue
			}

			writeStatus(line, percentMap, stations, csv_list, tierList)
			testing += 1
		}
	}

	return []string{}, nil
}

func writeStatus(line string, percentMap PercentMap, stations []stations.Station, csv_list []string, tierList []*tiers.Tiers) {
	/*
	*
	* We want to come up with an overall number out of 100.
	*  We will take each value from the line. It is always in the
	*  same order. Compute the amount, into the tier creating a
	*  percent.
	*
	* With that percent, we can then multiply the weight by the
	* created percent.
	*
	* With the final values, we can sum them and get the overal
	* int. With the, we can fit it into the tier mapping for what
	* the overall status -> percent overall is.
	*
	* Order:ID Delivered DCR DNRDPMO LoRDPMO POD CC CE DEX
	*
	 */

	csv_line := ""
	for idx, val := range strings.Split(line, " ") {
		fmt.Println("val: ", val, "\nidx: ", idx)
		switch idx {
		case 2:
			// DCR
			dcr_per := percentMap["dcr_val"]

			continue
		case 3:
			// DNRDPMO
			dnr_per := percentMap["dnr_dpmo_val"]
			continue
		case 4:
			// LoRDPMO
			ce_per := percentMap["ce_val"]
			continue
		case 5:
			// POD
			pod_per := percentMap["pod_val"]
			continue
		case 6:
			// CC
			cc_per := percentMap["cc_val"]
			continue
		case 7:
			// CE
			ce_per := percentMap["dex_val"]
			continue
		case 8:
			// DEX
			dex_per := percentMap["lor_val"]
			continue
		default:
			csv_line += "" + val + ","
		}
	}

	csv_list = append(csv_list, csv_line)
}

type (
	WrongObjCount map[string]int
	WrongObj      map[string]WrongObjCount
	PercentMap    map[string]float64
)
