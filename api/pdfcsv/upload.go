package pdfcsv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
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

	fetched_stations, err := stations.ReadAll(cc.Db)
	if err != nil {
		return nil, err
	}

	trimmed_file_name := strings.TrimSpace(fileHeader.Filename)

	var station_val stations.Station
	found := false
	for _, stat := range fetched_stations {
		if strings.Contains(trimmed_file_name, stat.Station) {
			found = true
			station_val = stat
		}
	}

	if !found {
		err := errors.New("could not find station")
		c.Logger().Error("could not find station")
		return nil, err
	}

	tierList, err := tiers.ReadTiers(cc.Db)
	if err != nil {
		return nil, err
	}

	tierMap := TierMap{}

	for _, val := range tierList {
		tierMap[val.Name] = val
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

	csv_headers := "Transporter ID,Delivered,DCR,DNR DPMO,LoR DPMO,POD,CC,CE,DEX\n"
	csv_list = append(csv_list, csv_headers)

	stringified_pdf, err := os.Open(txt_file_destination)
	if err != nil {
		return nil, err
	}

	defer stringified_pdf.Close()

	scanner := bufio.NewScanner(stringified_pdf)
	should_compute := false
	for scanner.Scan() {
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

			writeStatus(line, percentMap, station_val, &csv_list, tierMap)
		}
	}

	return csv_list, nil
}

func roundFloat(f float64, precision int) string {
	pow := math.Pow(10, float64(precision))
	roundedValue := math.Round(f*pow) / pow

	return strconv.FormatFloat(roundedValue, 'f', 2, 64)
}

func writeStatus(line string, percentMap PercentMap, station stations.Station, csv_list *[]string, tierMap TierMap) {
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
		*
	*/

	overall_tiers := []float64{
		.98, .95, .85, .7, .6,
	}

	overall_rating := []float64{
		98, 95, 85, 7, 6,
	}

	final_total := 0.00
	csv_line := ""

	for idx, val := range strings.Split(line, " ") {
		if idx > 8 {

			if final_total > overall_rating[0] {
				csv_line += "" + roundFloat(final_total, 2) + " | " + "fantastic plus\n"
				break
			}
			if final_total > overall_rating[1] {
				csv_line += "" + roundFloat(final_total, 2) + " | " + "fantastic\n"
				break
			}
			if final_total > overall_rating[2] {
				csv_line += "" + roundFloat(final_total, 2) + " | " + "great\n"
				break
			}
			if final_total > overall_rating[3] {
				csv_line += "" + roundFloat(final_total, 2) + " | " + "fair\n"
				break
			}
			if final_total > overall_rating[4] {
				csv_line += "" + roundFloat(final_total, 2) + " | " + "poor \n"
				break
			}
			csv_line += "" + roundFloat(final_total, 2) + " | " + "terrible\n"
			break
		}

		parsed_val := strings.TrimSpace(val)

		if strings.Contains(parsed_val, "%") {
			parsed_val = strings.TrimSuffix(parsed_val, "%")
		}

		if parsed_val == "-" {
			if idx == 3 || idx == 4 || idx == 7 {
				parsed_val = "0.00"
			} else {
				parsed_val = "100.00"
			}
		}

		if idx < 2 {
			csv_line += "" + val + ","
			continue
		}

		floatValue, err := strconv.ParseFloat(parsed_val, 64)
		if err != nil {
			fmt.Println("Error parsing float:", err)
			return
		}

		var tier *tiers.Tiers
		var per float64

		switch idx {
		case 2:

			// DCR
			per = percentMap["dcr_val"] * 100
			tier = tierMap["Dcr"]

		case 3:
			// DNRDPMO
			per = percentMap["dnr_dpmo_val"] * 100

			// Special case, less than instead of greater
			if floatValue < float64(station.Fan) {
				final_total += per
				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per, 2) + ","

				continue
			}
			if floatValue < float64(station.Great) {
				final_total += per * overall_tiers[1]
				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[1], 2) + ","

				continue
			}
			if floatValue < float64(station.Fair) {
				final_total += per * overall_tiers[2]
				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[2], 2) + ","

				continue
			}

			final_total += per * overall_tiers[4]
			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[2], 2) + ","
			continue
		case 4:
			// LoRDPMO
			per = percentMap["ce_val"] * 100
			tier = tierMap["Ce"]

		case 5:
			// POD
			per = percentMap["pod_val"] * 100
			tier = tierMap["Pod"]

		case 6:
			// CC
			per = percentMap["cc_val"] * 100
			tier = tierMap["Cc"]

		case 7:
			// CE
			per = percentMap["dex_val"] * 100
			tier = tierMap["Dex"]

		case 8:
			// DEX
			per = percentMap["lor_val"] * 100
			tier = tierMap["Lor"]

		default:
			continue
		}

		// Lower is better categories
		if idx == 4 || idx == 7 {
			if floatValue < tier.FanPlus {
				final_total += per

				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per, 2) + ","
				continue
			}
			if floatValue < tier.Fan {
				final_total += per * overall_tiers[1]

				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[1], 2) + ","
				continue
			}
			if floatValue < tier.Great {
				final_total += per * overall_tiers[2]
				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[2], 2) + ","
				continue
			}
			if floatValue < tier.Fair {
				final_total += per * overall_tiers[3]
				csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[3], 2) + ","
				continue
			}

			final_total += per * overall_tiers[4]
			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[4], 2) + ","

			continue
		}

		if floatValue > tier.FanPlus {
			final_total += per

			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per, 2) + ","
			continue
		}
		if floatValue > tier.Fan {
			final_total += per * overall_tiers[1]

			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[1], 2) + ","
			continue
		}
		if floatValue > tier.Great {
			final_total += per * overall_tiers[2]
			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[2], 2) + ","
			continue
		}
		if floatValue > tier.Fair {
			final_total += per * overall_tiers[3]
			csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[3], 2) + ","
			continue
		}

		final_total += per * overall_tiers[4]
		csv_line += "" + roundFloat(floatValue, 2) + " | " + roundFloat(per*overall_tiers[4], 2) + ","

	}

	*csv_list = append(*csv_list, csv_line)
}

type (
	WrongObjCount map[string]int
	WrongObj      map[string]WrongObjCount
	PercentMap    map[string]float64
	TierMap       map[string]*tiers.Tiers
)
