package pdfcsv

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	convertapi "github.com/ConvertAPI/convertapi-go/pkg"
	"github.com/ConvertAPI/convertapi-go/pkg/config"
	"github.com/ConvertAPI/convertapi-go/pkg/param"
)

func convert_pdf_to_text(filename string) ([][]string, error) {
	// Get the current working directory

	// Please don't take all my pdf's... Its really not 
	// a big thing, but still. 
	config.Default = config.NewDefault("token_QksTJT7R")

	filePath := filepath.Join("./uploads", filename)
	txt_file := strings.Replace(filePath, ".pdf", ".txt", -1)
	txt_file_destination := filepath.Join(txt_file)

	convertapi.ConvDef("pdf", "txt",
		param.NewPath("File", "filepath", nil),
	).ToPath("txt_file_destination")

	final_data_matrix, err := parse_text_file_created(txt_file_destination)
	if err != nil {
		return nil, err
	}
	return final_data_matrix, nil
}

func parse_text_file_created(filename string) ([][]string, error) {
	txt_file_to_parse, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer txt_file_to_parse.Close()

	scanner := bufio.NewScanner(txt_file_to_parse)
	var result []string
	var isCapturing bool
	var foundTransporterId bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check if we are starting to capture text (after 'Focus Area' is found)
		if strings.Contains(line, "Transporter ID") {
			foundTransporterId = true
		}

		if foundTransporterId && strings.Contains(line, "Focus Area") {
			isCapturing = true
		}

		// If we are capturing, add the line to result
		if isCapturing {
			result = append(result, " "+line)
		}

		// Stop capturing if 'Drivers With Working Hour Exceptions' is found
		if isCapturing && strings.Contains(line, "Drivers With Working Hour Exceptions") {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	/*
	*
	* We want to get all the values after the first empty line
	*
	* If there is another empty line we break
	*
	* We then want to iterate through everything and break on every white
	* space, so that every value is seperated. we then want to remove all
	* whitespace and group items.
	*
	 */

	var final_string string
	first_break := false
	var number_of_drivers int

	for i := 0; i < len(result); i++ {
		if !first_break {
			if strings.TrimSpace(result[i]) == "" {
				first_break = true
				number_of_drivers = len(strings.Fields(result[i+1]))
			}
			continue
		}

		final_string += result[i]
		if strings.TrimSpace(result[i]) == "" {
			break
		}
	}

	// These are the actual table values we want
	split_values := strings.Fields(final_string)

	// remove the percent signs
	for i, str := range split_values {
		split_values[i] = strings.ReplaceAll(str, "%", "")
	}

	var driver_data_matrix [][]string
	var current_dataset []string

	for i := 0; i < len(split_values); i++ {
		if i%number_of_drivers == 0 && i != 0 {
			driver_data_matrix = append(driver_data_matrix, current_dataset)
			// Reset the next array write

			current_dataset = []string{}
		}

		current_dataset = append(current_dataset, split_values[i])
	}
	driver_data_matrix = append(driver_data_matrix, current_dataset)

	current_dataset = []string{}

	// Make the actual lines of data
	var final_data_matrix [][]string

	// We expect the last array to be absolute bogus
	if len(driver_data_matrix[len(driver_data_matrix)-1]) < len(driver_data_matrix[0]) {
		driver_data_matrix = driver_data_matrix[:len(driver_data_matrix)-1]
	}

	for i := 0; i < number_of_drivers; i++ {
		for x := 0; x < len(driver_data_matrix); x++ {
			current_dataset = append(current_dataset, driver_data_matrix[x][i])
		}

		final_data_matrix = append(final_data_matrix, current_dataset)

		current_dataset = []string{}
	}

	return final_data_matrix, nil
}
