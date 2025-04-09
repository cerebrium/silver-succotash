package pdfcsv

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"hopdf.com/dao/pdfcsv"
	"hopdf.com/dao/stations"
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
		cc.Logger().Error("could not resolve cc")
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

	dst, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	// We need to get the pdf into a format that the below method can
	// actually work with (text)
	// final_data_set, err := convert_pdf_to_text(fileHeader.Filename)
	// if err != nil {
	// 	return nil, err
	// }

	// err = InternalConvertPdfToText(fileHeader.Filename)
	// if err != nil {
	// 	return nil, err
	// }

	local_weights := weights.Weights{ID: 1}
	updated_weights, err := local_weights.Read(cc.Db)
	if err != nil {
		return nil, err
	}

	stations, err := stations.ReadAll(cc.Db)
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

	csv_to_write, err := process_data(final_data_set, fileData.Name, percentMap, stations)
	if err != nil {
		return nil, err
	}

	return csv_to_write, nil
}

type (
	WrongObjCount map[string]int
	WrongObj      map[string]WrongObjCount
	PercentMap    map[string]float64
)

func process_data(final_data_set [][]string, file_name string, PercentMap PercentMap, stations []stations.Station) ([]string, error) {
	trimmed_file_name := strings.TrimSpace(file_name)

	stations_list := []string{"DRG2", "DSN1", "DBS3", "DBS2", "DEX2", "DCF1", "DSA1", "DPO1", "DOX2"}

	var station string

	for i := 0; i < len(stations_list); i++ {
		if strings.Contains(trimmed_file_name, stations_list[i]) {
			station = stations_list[i]
		}
	}

	csv_to_write, err := calculateStatuses(station, final_data_set, PercentMap, stations)
	if err != nil {
		return nil, err
	}

	return csv_to_write, nil
}

func calculateStatuses(station string, final_data_set [][]string, percentMap PercentMap, stations []stations.Station) ([]string, error) {
	final_csv := []string{"Transporter ID, Status, Delivered, DCR, DNR DPMO, POD, CC, CE, DEX\n"}

	dnrFan, dnrGreat, dnrFair := 1100, 1100, 1100

	// Find the station, and set the vals
	for _, val := range stations {
		if val.Station == station {
			dnrFair = val.Fair
			dnrGreat = val.Great
			dnrFan = val.Fan
		}
	}

	var totalCount int

	fantastic, great, fair, poor := 22.0, 20.5, 18.0, 13.0

	for i := 0; i < len(final_data_set); i++ {
		line := final_data_set[i]

		totalCount++
		currentRating := 0.0
		dontInclude := []string{}

		// Extract and parse fields
		dcr := line[2]
		dnrDpmo := line[3]
		pod := line[4]
		cc := line[5]
		ce := line[6]
		dex := line[7]

		// Process DCR
		var dcrVal float64
		if dcr == "" {
			continue
		}
		dcrVal, _ = strconv.ParseFloat(strings.TrimSuffix(dcr, "%"), 64)
		switch {
		case dcrVal >= 99:
			dcrVal = fantastic
		case dcrVal >= 98.75:
			dcrVal = great
		case dcrVal >= 98:
			dcrVal = fair
		default:
			dcrVal = poor
		}

		// Process DNR DPMO
		var dnrDpmoVal float64
		if dnrDpmo == "-" {
			dontInclude = append(dontInclude, "dnr_dpmo_val")
		} else {
			dnrDpmoVal, _ = strconv.ParseFloat(dnrDpmo, 64)
			switch {
			case dnrDpmoVal < float64(dnrFan):
				dnrDpmoVal = fantastic
			case dnrDpmoVal < float64(dnrGreat):
				dnrDpmoVal = great
			case dnrDpmoVal < float64(dnrFair):
				dnrDpmoVal = fair
			default:
				dnrDpmoVal = poor
			}
		}

		// Process POD
		var podVal float64
		if pod == "-" {
			dontInclude = append(dontInclude, "pod_val")
		} else {
			podVal, _ = strconv.ParseFloat(strings.TrimSuffix(pod, "%"), 64)
			switch {
			case podVal >= 98.9:
				podVal = fantastic
			case podVal > 98:
				podVal = great
			case podVal > 97:
				podVal = fair
			default:
				podVal = poor
			}
		}

		// Process CC
		var ccVal float64
		if cc == "-" {
			dontInclude = append(dontInclude, "cc_val")
		} else {
			ccVal, _ = strconv.ParseFloat(strings.TrimSuffix(cc, "%"), 64)
			switch {
			case ccVal > 98:
				ccVal = fantastic
			case ccVal > 95:
				ccVal = great
			case ccVal > 90:
				ccVal = fair
			default:
				ccVal = poor
			}
		}

		// Process DEX
		var dexVal float64
		if dex == "-" {
			dontInclude = append(dontInclude, "dex_val")
		} else {
			dexVal, _ = strconv.ParseFloat(strings.TrimSuffix(dex, "%"), 64)
			switch {
			case dexVal > 87:
				dexVal = fantastic
			case dexVal > 83:
				dexVal = great
			case dexVal > 80:
				dexVal = fair
			default:
				dexVal = poor
			}
		}

		// Process CE
		var ceVal float64
		ceVal, _ = strconv.ParseFloat(ce, 64)
		if ceVal == 0 {
			ceVal = fantastic
		} else {
			ceVal = poor
		}

		// Calculate rating
		options := []string{"dcr_val", "dnr_dpmo_val", "ce_val", "pod_val", "cc_val", "dex_val"}
		missingPercent := 0.0
		for _, option := range options {
			if contains(dontInclude, option) {
				missingPercent += percentMap[option]
			}
		}
		multiplicative := 1 / (1 - missingPercent)

		if !contains(dontInclude, "dcr_val") {
			currentRating += dcrVal * (0.35 * multiplicative)
		}
		if !contains(dontInclude, "dnr_dpmo_val") {
			currentRating += dnrDpmoVal * (0.35 * multiplicative)
		}
		if !contains(dontInclude, "ce_val") {
			currentRating += ceVal * (0.075 * multiplicative)
		}
		if !contains(dontInclude, "pod_val") {
			currentRating += podVal * (0.075 * multiplicative)
		}
		if !contains(dontInclude, "cc_val") {
			currentRating += ccVal * (0.075 * multiplicative)
		}
		if !contains(dontInclude, "dex_val") {
			currentRating += dexVal * (0.075 * multiplicative)
		}

		// Determine status
		status := determineStatus(currentRating)
		content := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s\n", line[0], status, line[1], dcr, dnrDpmo, pod, cc, ce, dex)
		final_csv = append(final_csv, content)
	}

	return final_csv, nil
}

func determineStatus(rating float64) string {
	switch {
	case rating > 21.25:
		return "FANTASTIC_PLUS"
	case rating > 20.2:
		return "FANTASTIC"
	case rating > 18.85:
		return "GREAT"
	case rating > 17.951:
		return "FAIR"
	default:
		return "POOR"
	}
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
