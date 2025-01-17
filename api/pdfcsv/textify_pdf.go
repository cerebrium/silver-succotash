package pdfcsv

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	convertapi "github.com/ConvertAPI/convertapi-go/pkg"
	"github.com/ConvertAPI/convertapi-go/pkg/config"
	"github.com/ConvertAPI/convertapi-go/pkg/param"
)

func convert_pdf_to_text(filename string) ([][]string, error) {
	// Get the current working directory

	filePath := filepath.Join("./uploads", filename)
	txt_file := strings.Replace(filePath, ".pdf", ".txt", -1)
	txt_file_destination := filepath.Join(txt_file)

	// Please don't take all my pdf's... Its really not
	// a big thing, but still.
	config.Default = config.NewDefault("token_QksTJT7R")

	_, err_arr := convertapi.ConvDef("pdf", "txt",
		param.NewPath("File", filePath, nil),
	).ToPath(txt_file_destination)

	if err_arr != nil {
		return nil, err_arr[0]
	}

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
	var result [][]string
	var isCapturing bool
	var foundTransporterId bool
	current_page_idx := 0

	// Add the first slice to the matrix
	result = append(result, []string{})

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
			result[current_page_idx] = append(result[current_page_idx], " "+line)
		}

		if isCapturing && strings.Contains(line, "Page") {
			// This means there are multiple pages of content
			// The content comes in the same order on each page.
			// We will need to write it in the same way to a different
			// results slice.
			result = append(result, []string{})
			current_page_idx++
			continue
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

	// Build the slice of string
	final_strings := []string{}
	first_break := false
	var number_of_drivers int
	driver_numbers := []int{}

	// There is a 'page n' in here, but if we append it to the last
	// list, then we can ignore it
	for m_idx := 0; m_idx < len(result)-1; m_idx++ {
		final_strings = append(final_strings, "")

		// The second list starts with an empty space
		for i := 0; i < len(result[m_idx]); i++ {
			if m_idx > 0 && m_idx > len(driver_numbers)-1 {
				// We should be able to walk down the array here
				// to find the next list of driver id's and create
				// our count
				q := 0
				for {
					if q == len(result[m_idx]) {

						err := errors.New("could not find next driver id list")
						return nil, err
					}
					if strings.TrimSpace(result[m_idx][q]) != "" {
						next_page_driver_count := len(strings.Fields(result[m_idx][q]))

						driver_numbers = append(driver_numbers, next_page_driver_count)
						break
					}
					q++
				}

				if strings.TrimSpace(result[m_idx][i]) == "" {
					continue
				}
			}
			if !first_break {
				if strings.TrimSpace(result[m_idx][i]) == "" {
					first_break = true
					// This should get the length of each slice in the matrix
					number_of_drivers = len(strings.Fields(result[m_idx][i+1]))
					driver_numbers = append(driver_numbers, number_of_drivers)
				}
				continue
			}

			// The pages past the first never create a driver number

			if strings.Contains(result[m_idx][i], "Page") {
				break
			}

			final_strings[m_idx] += result[m_idx][i]
			if strings.TrimSpace(result[m_idx][i]) == "" {
				break
			}

		}
	}

	driver_data_matrix := [][]string{}

	// We need to pop the last element off,
	// so using a slice, not an array.
	for i := 0; i < 10; i++ {
		driver_data_matrix = append(driver_data_matrix, []string{})
	}

	// split the strings, and write to the matrix
	for str_idx := 0; str_idx < len(final_strings); str_idx++ {
		final_string := final_strings[str_idx]
		// These are the actual table values we want
		split_values := strings.Fields(final_string)

		// remove the percent signs
		for i, str := range split_values {
			split_values[i] = strings.ReplaceAll(str, "%", "")
		}

		m_idx := 0
		for i := 0; i < len(split_values); i++ {
			if i%driver_numbers[str_idx] == 0 && i != 0 {
				m_idx++
			}

			driver_data_matrix[m_idx] = append(driver_data_matrix[m_idx], split_values[i])
		}

	}

	current_dataset := []string{}

	// Make the actual lines of data
	var final_data_matrix [][]string

	// We expect the last array to be absolute bogus
	driver_data_matrix = driver_data_matrix[:9]

	// We need to group the data by driver
	for i := 0; i < len(driver_data_matrix[0]); i++ {
		for x := 0; x < len(driver_data_matrix); x++ {
			current_dataset = append(current_dataset, driver_data_matrix[x][i])
		}

		final_data_matrix = append(final_data_matrix, current_dataset)

		current_dataset = []string{}
	}

	return final_data_matrix, nil
}
