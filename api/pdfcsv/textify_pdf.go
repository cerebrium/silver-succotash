package pdfcsv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func InternalConvertPdfToText(filename string) error {
	filePath := filepath.Join("./uploads", filename)
	txt_file := strings.Replace(filePath, ".pdf", ".txt", -1)
	txt_file_destination := filepath.Join(txt_file)

	unidoc_key := os.Getenv("UNICODE_SECRET_KEY")
	fmt.Println("\n\n THE CODE: ", unidoc_key, "\n\n")
	err := license.SetMeteredKey(unidoc_key)
	if err != nil {
		fmt.Printf("error with unicode key: ", err)
		return err
	}

	err = extractTextToFile(filePath, txt_file_destination)
	if err != nil {
		fmt.Printf("Error in extraction: ", err)
		return err
	}

	return nil
}

func extractTextToFile(inputPath, outputPath string) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening PDF file: %w", err)
	}
	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return fmt.Errorf("error creating PDF reader: %w", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return fmt.Errorf("error getting number of pages: %w", err)
	}

	var allText strings.Builder

	for i := 1; i <= numPages; i++ {
		page, err := pdfReader.GetPage(i)
		if err != nil {
			fmt.Printf("Error getting page %d: %v\n", i, err)
			continue
		}

		textExtractor, err := extractor.New(page)
		if err != nil {
			fmt.Printf("Error creating text extractor for page %d: %v\n", i, err)
			continue
		}

		text, err := textExtractor.ExtractText()
		if err != nil {
			fmt.Printf("Error extracting text from page %d: %v\n", i, err)
			continue
		}

		allText.WriteString(text)
		allText.WriteString("\n\n--- Page ")
		allText.WriteString(fmt.Sprintf("%d", i))
		allText.WriteString(" ---\n\n")
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outFile.Close()

	_, err = outFile.WriteString(allText.String())
	if err != nil {
		return fmt.Errorf("error writing text to output file: %w", err)
	}

	return nil
}
